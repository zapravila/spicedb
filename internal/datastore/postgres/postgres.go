package postgres

import (
	"context"
	dbsql "database/sql"
	"errors"
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"github.com/IBM/pgxpoolprometheus"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/mattn/go-isatty"
	"github.com/ngrok/sqlmw"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/schollz/progressbar/v3"
	"go.opentelemetry.io/otel"
	"golang.org/x/sync/errgroup"

	datastoreinternal "github.com/zapravila/spicedb/internal/datastore"
	"github.com/zapravila/spicedb/internal/datastore/common"
	pgxcommon "github.com/zapravila/spicedb/internal/datastore/postgres/common"
	"github.com/zapravila/spicedb/internal/datastore/postgres/migrations"
	"github.com/zapravila/spicedb/internal/datastore/revisions"
	log "github.com/zapravila/spicedb/internal/logging"
	"github.com/zapravila/spicedb/pkg/datastore"
	"github.com/zapravila/spicedb/pkg/datastore/options"
)

func init() {
	datastore.Engines = append(datastore.Engines, Engine)
}

const (
	Engine                   = "postgres"
	tableNamespace           = "namespace_config"
	tableTransaction         = "relation_tuple_transaction"
	tableTuple               = "relation_tuple"
	tableCaveat              = "caveat"
	tableRelationshipCounter = "relationship_counter"

	colXID               = "xid"
	colTimestamp         = "timestamp"
	colNamespace         = "namespace"
	colConfig            = "serialized_config"
	colCreatedXid        = "created_xid"
	colDeletedXid        = "deleted_xid"
	colSnapshot          = "snapshot"
	colObjectID          = "object_id"
	colRelation          = "relation"
	colUsersetNamespace  = "userset_namespace"
	colUsersetObjectID   = "userset_object_id"
	colUsersetRelation   = "userset_relation"
	colCaveatName        = "name"
	colCaveatDefinition  = "definition"
	colCaveatContextName = "caveat_name"
	colCaveatContext     = "caveat_context"
	colDescription       = "description"
	colComment           = "comment"

	colCounterName         = "name"
	colCounterFilter       = "serialized_filter"
	colCounterCurrentCount = "current_count"
	colCounterSnapshot     = "updated_revision_snapshot"

	errUnableToInstantiate = "unable to instantiate datastore"

	// The parameters to this format string are:
	// 1: the created_xid or deleted_xid column name
	//
	// The placeholders are the snapshot and the expected boolean value respectively.
	snapshotAlive = "pg_visible_in_snapshot(%[1]s, ?) = ?"

	// This is the largest positive integer possible in postgresql
	liveDeletedTxnID = uint64(9223372036854775807)

	tracingDriverName = "postgres-tracing"

	gcBatchDeleteSize = 1000
)

var livingTupleConstraints = []string{"uq_relation_tuple_living_xid", "pk_relation_tuple"}

func init() {
	dbsql.Register(tracingDriverName, sqlmw.Driver(stdlib.GetDefaultDriver(), new(traceInterceptor)))
}

var (
	psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	getRevision = psql.
			Select(colXID, colSnapshot).
			From(tableTransaction).
			OrderByClause(fmt.Sprintf("%s DESC", colXID)).
			Limit(1)

	createTxn = fmt.Sprintf(
		"INSERT INTO %s DEFAULT VALUES RETURNING %s, %s",
		tableTransaction,
		colXID,
		colSnapshot,
	)

	getNow = psql.Select("NOW()")

	tracer = otel.Tracer("spicedb/internal/datastore/postgres")
)

type sqlFilter interface {
	ToSql() (string, []interface{}, error)
}

// NewPostgresDatastore initializes a SpiceDB datastore that uses a PostgreSQL
// database by leveraging manual book-keeping to implement revisioning.
//
// This datastore is also tested to be compatible with CockroachDB.
func NewPostgresDatastore(
	ctx context.Context,
	url string,
	options ...Option,
) (datastore.Datastore, error) {
	ds, err := newPostgresDatastore(ctx, url, options...)
	if err != nil {
		return nil, err
	}

	return datastoreinternal.NewSeparatingContextDatastoreProxy(ds), nil
}

func newPostgresDatastore(
	ctx context.Context,
	pgURL string,
	options ...Option,
) (datastore.Datastore, error) {
	config, err := generateConfig(options)
	if err != nil {
		return nil, common.RedactAndLogSensitiveConnString(ctx, errUnableToInstantiate, err, pgURL)
	}

	// Parse the DB URI into configuration.
	parsedConfig, err := pgxpool.ParseConfig(pgURL)
	if err != nil {
		return nil, common.RedactAndLogSensitiveConnString(ctx, errUnableToInstantiate, err, pgURL)
	}

	// Setup the default custom plan setting, if applicable.
	pgConfig, err := defaultCustomPlan(parsedConfig)
	if err != nil {
		return nil, common.RedactAndLogSensitiveConnString(ctx, errUnableToInstantiate, err, pgURL)
	}

	// Setup the credentials provider
	var credentialsProvider datastore.CredentialsProvider
	if config.credentialsProviderName != "" {
		credentialsProvider, err = datastore.NewCredentialsProvider(ctx, config.credentialsProviderName)
		if err != nil {
			return nil, err
		}
	}

	// Setup the config for each of the read and write pools.
	readPoolConfig := pgConfig.Copy()
	config.readPoolOpts.ConfigurePgx(readPoolConfig)

	readPoolConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		RegisterTypes(conn.TypeMap())
		return nil
	}

	writePoolConfig := pgConfig.Copy()
	config.writePoolOpts.ConfigurePgx(writePoolConfig)

	writePoolConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		RegisterTypes(conn.TypeMap())
		return nil
	}

	if credentialsProvider != nil {
		// add before connect callbacks to trigger the token
		getToken := func(ctx context.Context, config *pgx.ConnConfig) error {
			config.User, config.Password, err = credentialsProvider.Get(ctx, fmt.Sprintf("%s:%d", config.Host, config.Port), config.User)
			return err
		}
		readPoolConfig.BeforeConnect = getToken
		writePoolConfig.BeforeConnect = getToken
	}

	if config.migrationPhase != "" {
		log.Info().
			Str("phase", config.migrationPhase).
			Msg("postgres configured to use intermediate migration phase")
	}

	initializationContext, cancelInit := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelInit()

	readPool, err := pgxpool.NewWithConfig(initializationContext, readPoolConfig)
	if err != nil {
		return nil, common.RedactAndLogSensitiveConnString(ctx, errUnableToInstantiate, err, pgURL)
	}

	writePool, err := pgxpool.NewWithConfig(initializationContext, writePoolConfig)
	if err != nil {
		return nil, common.RedactAndLogSensitiveConnString(ctx, errUnableToInstantiate, err, pgURL)
	}

	// Verify that the server supports commit timestamps
	var trackTSOn string
	if err := readPool.
		QueryRow(initializationContext, "SHOW track_commit_timestamp;").
		Scan(&trackTSOn); err != nil {
		return nil, err
	}

	watchEnabled := trackTSOn == "on"
	if !watchEnabled {
		log.Warn().Msg("watch API disabled, postgres must be run with track_commit_timestamp=on")
	}

	if config.enablePrometheusStats {
		if err := prometheus.Register(pgxpoolprometheus.NewCollector(readPool, map[string]string{
			"db_name":    "spicedb",
			"pool_usage": "read",
		})); err != nil {
			return nil, err
		}
		if err := prometheus.Register(pgxpoolprometheus.NewCollector(writePool, map[string]string{
			"db_name":    "spicedb",
			"pool_usage": "write",
		})); err != nil {
			return nil, err
		}
		if err := common.RegisterGCMetrics(); err != nil {
			return nil, err
		}
	}

	gcCtx, cancelGc := context.WithCancel(context.Background())

	quantizationPeriodNanos := config.revisionQuantization.Nanoseconds()
	if quantizationPeriodNanos < 1 {
		quantizationPeriodNanos = 1
	}
	revisionQuery := fmt.Sprintf(
		querySelectRevision,
		colXID,
		tableTransaction,
		colTimestamp,
		quantizationPeriodNanos,
		colSnapshot,
	)

	validTransactionQuery := fmt.Sprintf(
		queryValidTransaction,
		colXID,
		tableTransaction,
		colTimestamp,
		config.gcWindow.Seconds(),
		colSnapshot,
	)

	maxRevisionStaleness := time.Duration(float64(config.revisionQuantization.Nanoseconds())*
		config.maxRevisionStalenessPercent) * time.Nanosecond

	datastore := &pgDatastore{
		CachedOptimizedRevisions: revisions.NewCachedOptimizedRevisions(
			maxRevisionStaleness,
		),
		dburl:                   pgURL,
		readPool:                pgxcommon.MustNewInterceptorPooler(readPool, config.queryInterceptor),
		writePool:               pgxcommon.MustNewInterceptorPooler(writePool, config.queryInterceptor),
		watchBufferLength:       config.watchBufferLength,
		watchBufferWriteTimeout: config.watchBufferWriteTimeout,
		optimizedRevisionQuery:  revisionQuery,
		validTransactionQuery:   validTransactionQuery,
		gcWindow:                config.gcWindow,
		gcInterval:              config.gcInterval,
		gcTimeout:               config.gcMaxOperationTime,
		analyzeBeforeStatistics: config.analyzeBeforeStatistics,
		watchEnabled:            watchEnabled,
		gcCtx:                   gcCtx,
		cancelGc:                cancelGc,
		readTxOptions:           pgx.TxOptions{IsoLevel: pgx.RepeatableRead, AccessMode: pgx.ReadOnly},
		maxRetries:              config.maxRetries,
		credentialsProvider:     credentialsProvider,
	}

	datastore.SetOptimizedRevisionFunc(datastore.optimizedRevisionFunc)

	// Start a goroutine for garbage collection.
	if datastore.gcInterval > 0*time.Minute && config.gcEnabled {
		datastore.gcGroup, datastore.gcCtx = errgroup.WithContext(datastore.gcCtx)
		datastore.gcGroup.Go(func() error {
			return common.StartGarbageCollector(
				datastore.gcCtx,
				datastore,
				datastore.gcInterval,
				datastore.gcWindow,
				datastore.gcTimeout,
			)
		})
	} else {
		log.Warn().Msg("datastore background garbage collection disabled")
	}

	return datastore, nil
}

type pgDatastore struct {
	*revisions.CachedOptimizedRevisions

	dburl                   string
	readPool, writePool     pgxcommon.ConnPooler
	watchBufferLength       uint16
	watchBufferWriteTimeout time.Duration
	optimizedRevisionQuery  string
	validTransactionQuery   string
	gcWindow                time.Duration
	gcInterval              time.Duration
	gcTimeout               time.Duration
	analyzeBeforeStatistics bool
	readTxOptions           pgx.TxOptions
	maxRetries              uint8
	watchEnabled            bool

	credentialsProvider datastore.CredentialsProvider

	gcGroup  *errgroup.Group
	gcCtx    context.Context
	cancelGc context.CancelFunc
	gcHasRun atomic.Bool
}

func (pgd *pgDatastore) SnapshotReader(revRaw datastore.Revision) datastore.Reader {
	rev := revRaw.(postgresRevision)

	queryFuncs := pgxcommon.QuerierFuncsFor(pgd.readPool)
	executor := common.QueryExecutor{
		Executor: pgxcommon.NewPGXExecutor(queryFuncs),
	}

	return &pgReader{
		queryFuncs,
		executor,
		buildLivingObjectFilterForRevision(rev),
	}
}

// ReadWriteTx starts a read/write transaction, which will be committed if no error is
// returned and rolled back if an error is returned.
func (pgd *pgDatastore) ReadWriteTx(
	ctx context.Context,
	fn datastore.TxUserFunc,
	opts ...options.RWTOptionsOption,
) (datastore.Revision, error) {
	config := options.NewRWTOptionsWithOptions(opts...)

	var err error
	for i := uint8(0); i <= pgd.maxRetries; i++ {
		var newXID xid8
		var newSnapshot pgSnapshot
		err = wrapError(pgx.BeginTxFunc(ctx, pgd.writePool, pgx.TxOptions{IsoLevel: pgx.Serializable}, func(tx pgx.Tx) error {
			var err error
			newXID, newSnapshot, err = createNewTransaction(ctx, tx)
			if err != nil {
				return err
			}

			queryFuncs := pgxcommon.QuerierFuncsFor(pgd.readPool)
			executor := common.QueryExecutor{
				Executor: pgxcommon.NewPGXExecutor(queryFuncs),
			}

			rwt := &pgReadWriteTXN{
				&pgReader{
					queryFuncs,
					executor,
					currentlyLivingObjects,
				},
				tx,
				newXID,
			}

			return fn(ctx, rwt)
		}))
		if err != nil {
			if !config.DisableRetries && errorRetryable(err) {
				pgxcommon.SleepOnErr(ctx, err, i)
				continue
			}

			return datastore.NoRevision, err
		}

		if i > 0 {
			log.Debug().Uint8("retries", i).Msg("transaction succeeded after retry")
		}

		return postgresRevision{newSnapshot.markComplete(newXID.Uint64)}, nil
	}

	if !config.DisableRetries {
		err = fmt.Errorf("max retries exceeded: %w", err)
	}

	return datastore.NoRevision, err
}

const repairTransactionIDsOperation = "transaction-ids"

func (pgd *pgDatastore) Repair(ctx context.Context, operationName string, outputProgress bool) error {
	switch operationName {
	case repairTransactionIDsOperation:
		return pgd.repairTransactionIDs(ctx, outputProgress)

	default:
		return fmt.Errorf("unknown operation")
	}
}

const batchSize = 10000

func (pgd *pgDatastore) repairTransactionIDs(ctx context.Context, outputProgress bool) error {
	conn, err := pgx.Connect(ctx, pgd.dburl)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	// Get the current transaction ID.
	currentMaximumID := 0
	if err := conn.QueryRow(ctx, queryCurrentTransactionID).Scan(&currentMaximumID); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return err
		}
	}

	// Find the maximum transaction ID referenced in the transactions table.
	referencedMaximumID := 0
	if err := conn.QueryRow(ctx, queryLatestXID).Scan(&referencedMaximumID); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return err
		}
	}

	// The delta is what this needs to fill in.
	log.Ctx(ctx).Info().Int64("current-maximum", int64(currentMaximumID)).Int64("referenced-maximum", int64(referencedMaximumID)).Msg("found transactions")
	counterDelta := referencedMaximumID - currentMaximumID
	if counterDelta < 0 {
		return nil
	}

	var bar *progressbar.ProgressBar
	if isatty.IsTerminal(os.Stderr.Fd()) && outputProgress {
		bar = progressbar.Default(int64(counterDelta), "updating transactions counter")
	}

	for i := 0; i < counterDelta; i++ {
		var batch pgx.Batch

		batchCount := min(batchSize, counterDelta-i)
		for j := 0; j < batchCount; j++ {
			batch.Queue("begin;")
			batch.Queue("select pg_current_xact_id();")
			batch.Queue("rollback;")
		}

		br := conn.SendBatch(ctx, &batch)
		if err := br.Close(); err != nil {
			return err
		}

		i += batchCount - 1
		if bar != nil {
			if err := bar.Add(batchCount); err != nil {
				return err
			}
		}
	}

	if bar != nil {
		if err := bar.Close(); err != nil {
			return err
		}
	}

	log.Ctx(ctx).Info().Msg("completed revisions repair")
	return nil
}

// RepairOperations returns the available repair operations for the datastore.
func (pgd *pgDatastore) RepairOperations() []datastore.RepairOperation {
	return []datastore.RepairOperation{
		{
			Name:        repairTransactionIDsOperation,
			Description: "Brings the Postgres database up to the expected transaction ID (Postgres v15+ only)",
		},
	}
}

func wrapError(err error) error {
	// If a unique constraint violation is returned, then its likely that the cause
	// was an existing relationship given as a CREATE.
	if cerr := pgxcommon.ConvertToWriteConstraintError(livingTupleConstraints, err); cerr != nil {
		return cerr
	}

	if pgxcommon.IsSerializationError(err) {
		return common.NewSerializationError(err)
	}

	// hack: pgx asyncClose usually happens after cancellation,
	// but the reason for it being closed is not propagated
	// and all we get is attempting to perform an operation
	// on cancelled connection. This keeps the same error,
	// but wrapped along a cancellation so that:
	// - pgx logger does not log it
	// - response is sent as canceled back to the client
	if err != nil && err.Error() == "conn closed" {
		return errors.Join(err, context.Canceled)
	}

	return err
}

func (pgd *pgDatastore) Close() error {
	pgd.cancelGc()

	if pgd.gcGroup != nil {
		err := pgd.gcGroup.Wait()
		log.Warn().Err(err).Msg("completed shutdown of postgres datastore")
	}

	pgd.readPool.Close()
	pgd.writePool.Close()
	return nil
}

func errorRetryable(err error) bool {
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return false
	}

	if pgconn.SafeToRetry(err) {
		return true
	}

	if pgxcommon.IsSerializationError(err) {
		return true
	}

	log.Warn().Err(err).Msg("unable to determine if pgx error is retryable")
	return false
}

func (pgd *pgDatastore) ReadyState(ctx context.Context) (datastore.ReadyState, error) {
	headMigration, err := migrations.DatabaseMigrations.HeadRevision()
	if err != nil {
		return datastore.ReadyState{}, fmt.Errorf("invalid head migration found for postgres: %w", err)
	}

	pgDriver, err := migrations.NewAlembicPostgresDriver(ctx, pgd.dburl, pgd.credentialsProvider)
	if err != nil {
		return datastore.ReadyState{}, err
	}
	defer pgDriver.Close(ctx)

	version, err := pgDriver.Version(ctx)
	if err != nil {
		return datastore.ReadyState{}, err
	}

	if version == headMigration {
		// Ensure a datastore ID is present. This ensures the tables have not been truncated.
		uniqueID, err := pgd.datastoreUniqueID(ctx)
		if err != nil {
			return datastore.ReadyState{}, fmt.Errorf("database validation failed: %w; if you have previously run `TRUNCATE`, this database is no longer valid and must be remigrated. See: https://spicedb.dev/d/truncate-unsupported", err)
		}

		log.Trace().Str("unique_id", uniqueID).Msg("postgres datastore unique ID")
		return datastore.ReadyState{IsReady: true}, nil
	}

	return datastore.ReadyState{
		Message: fmt.Sprintf(
			"datastore is not migrated: currently at revision `%s`, but requires `%s`. Please run `spicedb migrate`. If you have previously run `TRUNCATE`, this database is no longer valid and must be remigrated. See: https://spicedb.dev/d/truncate-unsupported",
			version,
			headMigration,
		),
		IsReady: false,
	}, nil
}

func (pgd *pgDatastore) Features(_ context.Context) (*datastore.Features, error) {
	return &datastore.Features{Watch: datastore.Feature{Enabled: pgd.watchEnabled}}, nil
}

func buildLivingObjectFilterForRevision(revision postgresRevision) queryFilterer {
	createdBeforeTXN := sq.Expr(fmt.Sprintf(
		snapshotAlive,
		colCreatedXid,
	), revision.snapshot, true)

	deletedAfterTXN := sq.Expr(fmt.Sprintf(
		snapshotAlive,
		colDeletedXid,
	), revision.snapshot, false)

	return func(original sq.SelectBuilder) sq.SelectBuilder {
		return original.Where(createdBeforeTXN).Where(deletedAfterTXN)
	}
}

func currentlyLivingObjects(original sq.SelectBuilder) sq.SelectBuilder {
	return original.Where(sq.Eq{colDeletedXid: liveDeletedTxnID})
}

// defaultCustomPlan parses a Postgres URI and determines if a plan_cache_mode
// has been specified. If not, it defaults to "force_custom_plan".
// This works around a bug impacting performance documented here:
// https://spicedb.dev/d/force-custom-plan.
func defaultCustomPlan(poolConfig *pgxpool.Config) (*pgxpool.Config, error) {
	if existing, ok := poolConfig.ConnConfig.Config.RuntimeParams["plan_cache_mode"]; ok {
		log.Info().
			Str("plan_cache_mode", existing).
			Msg("found plan_cache_mode in DB URI; leaving as-is")
		return poolConfig, nil
	}

	poolConfig.ConnConfig.Config.RuntimeParams["plan_cache_mode"] = "force_custom_plan"
	log.Warn().
		Str("details-url", "https://spicedb.dev/d/force-custom-plan").
		Str("plan_cache_mode", "force_custom_plan").
		Msg("defaulting value in Postgres DB URI")

	return poolConfig, nil
}

var _ datastore.Datastore = &pgDatastore{}
