package proxy

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"hash"
	"sync"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	log "github.com/zapravila/spicedb/internal/logging"
	"github.com/zapravila/spicedb/pkg/datastore"
	"github.com/zapravila/spicedb/pkg/datastore/options"
	corev1 "github.com/zapravila/spicedb/pkg/proto/core/v1"
	"github.com/zapravila/spicedb/pkg/spiceerrors"
	"github.com/zapravila/spicedb/pkg/tuple"
)

// KeyConfig is a configuration for a key used to sign relationships.
type KeyConfig struct {
	// ID is the unique identifier for the key.
	ID string

	// ExpiredAt is the time at which the key is no longer valid, if any.
	ExpiredAt *time.Time

	// Bytes is the raw key material.
	Bytes []byte
}

type hmacConfig struct {
	keyID     string
	expiredAt *time.Time
	pool      sync.Pool
}

var (
	versionByte = byte(0x01)
	hashLength  = 16
)

// NewRelationshipIntegrityProxy creates a new datastore proxy that ensures the integrity of
// relationships by using HMACs to sign the data. The current key is used to sign new data,
// and the expired keys are used to verify old data, if any.
func NewRelationshipIntegrityProxy(ds datastore.Datastore, currentKey KeyConfig, expiredKeys []KeyConfig) (datastore.Datastore, error) {
	// Ensure the datastore supports integrity.
	features, err := ds.OfflineFeatures()
	if err != nil {
		return nil, err
	}

	if features.IntegrityData.Status != datastore.FeatureSupported {
		return nil, spiceerrors.MustBugf("datastore does not support relationship integrity")
	}

	if len(currentKey.Bytes) == 0 {
		return nil, fmt.Errorf("contents of the current key file cannot be empty")
	}

	if len(currentKey.ID) == 0 {
		return nil, fmt.Errorf("current key ID cannot be empty")
	}

	currentKeyHMAC := &hmacConfig{
		keyID:     currentKey.ID,
		expiredAt: currentKey.ExpiredAt,
		pool:      poolForKey(currentKey.Bytes),
	}

	if currentKey.ExpiredAt != nil {
		return nil, spiceerrors.MustBugf("current key cannot have an expiration")
	}

	keysByID := make(map[string]*hmacConfig, len(expiredKeys)+1)
	keysByID[currentKey.ID] = currentKeyHMAC

	expiredKeyIDs := make([]string, 0, len(expiredKeys))
	for _, key := range expiredKeys {
		if len(key.Bytes) == 0 {
			return nil, fmt.Errorf("expired key cannot be empty")
		}

		if len(key.ID) == 0 {
			return nil, fmt.Errorf("expired key ID cannot be empty")
		}

		if key.ExpiredAt == nil {
			return nil, fmt.Errorf("expired key missing expiration time")
		}

		if _, ok := keysByID[key.ID]; ok {
			return nil, fmt.Errorf("found duplicate key ID: %s", key.ID)
		}

		keysByID[key.ID] = &hmacConfig{
			keyID:     key.ID,
			expiredAt: key.ExpiredAt,
			pool:      poolForKey(key.Bytes),
		}

		expiredKeyIDs = append(expiredKeyIDs, key.ID)
	}

	log.Debug().
		Str("current_key_id", currentKey.ID).
		Strs("expired_key_ids", expiredKeyIDs).
		Msgf("created relationship integrity proxy")

	return &relationshipIntegrityProxy{
		ds:         ds,
		primaryKey: currentKeyHMAC,
		keysByID:   keysByID,
	}, nil
}

func poolForKey(key []byte) sync.Pool {
	return sync.Pool{
		New: func() any {
			return hmac.New(sha256.New, key)
		},
	}
}

type relationshipIntegrityProxy struct {
	ds         datastore.Datastore
	primaryKey *hmacConfig
	keysByID   map[string]*hmacConfig
}

func (r *relationshipIntegrityProxy) lookupKey(keyID string) (*hmacConfig, error) {
	key, ok := r.keysByID[keyID]
	if !ok {
		return nil, fmt.Errorf("key not found: %s", keyID)
	}

	return key, nil
}

// computeRelationshipHash computes the HMAC hash of a relationship tuple.
func computeRelationshipHash(tpl *corev1.RelationTuple, key *hmacConfig) ([]byte, error) {
	bytes, err := tuple.CanonicalBytes(tpl)
	if err != nil {
		return nil, err
	}

	hasher := key.pool.Get().(hash.Hash)
	defer key.pool.Put(hasher)

	hasher.Reset()
	if _, err := hasher.Write(bytes); err != nil {
		return nil, err
	}

	return hasher.Sum(nil)[:hashLength], nil
}

func (r *relationshipIntegrityProxy) SnapshotReader(rev datastore.Revision) datastore.Reader {
	return relationshipIntegrityReader{
		parent:  r,
		wrapped: r.ds.SnapshotReader(rev),
	}
}

func (r *relationshipIntegrityProxy) ReadWriteTx(ctx context.Context, f datastore.TxUserFunc, opts ...options.RWTOptionsOption) (datastore.Revision, error) {
	return r.ds.ReadWriteTx(ctx, func(ctx context.Context, tx datastore.ReadWriteTransaction) error {
		return f(ctx, &relationshipIntegrityTx{
			ReadWriteTransaction: tx,
			parent:               r,
		})
	}, opts...)
}

func (r *relationshipIntegrityProxy) CheckRevision(ctx context.Context, revision datastore.Revision) error {
	return r.ds.CheckRevision(ctx, revision)
}

func (r *relationshipIntegrityProxy) Close() error {
	return r.ds.Close()
}

func (r *relationshipIntegrityProxy) Features(ctx context.Context) (*datastore.Features, error) {
	return r.ds.Features(ctx)
}

func (r *relationshipIntegrityProxy) OfflineFeatures() (*datastore.Features, error) {
	return r.ds.OfflineFeatures()
}

func (r *relationshipIntegrityProxy) HeadRevision(ctx context.Context) (datastore.Revision, error) {
	return r.ds.HeadRevision(ctx)
}

func (r *relationshipIntegrityProxy) OptimizedRevision(ctx context.Context) (datastore.Revision, error) {
	return r.ds.OptimizedRevision(ctx)
}

func (r *relationshipIntegrityProxy) ReadyState(ctx context.Context) (datastore.ReadyState, error) {
	return r.ds.ReadyState(ctx)
}

func (r *relationshipIntegrityProxy) RevisionFromString(serialized string) (datastore.Revision, error) {
	return r.ds.RevisionFromString(serialized)
}

func (r *relationshipIntegrityProxy) Statistics(ctx context.Context) (datastore.Stats, error) {
	return r.ds.Statistics(ctx)
}

func (r *relationshipIntegrityProxy) validateRelationTuple(tpl *corev1.RelationTuple) error {
	// Ensure the relationship has integrity data.
	if tpl.Integrity == nil || len(tpl.Integrity.Hash) == 0 || tpl.Integrity.KeyId == "" {
		str, err := tuple.String(tpl)
		if err != nil {
			return err
		}

		return fmt.Errorf("relationship %s is missing required integrity data", str)
	}

	hashWithoutByte := tpl.Integrity.Hash[1:]
	if tpl.Integrity.Hash[0] != versionByte || len(hashWithoutByte) != hashLength {
		return fmt.Errorf("relationship %s has invalid integrity data", tpl)
	}

	// Validate the integrity of the relationship.
	key, err := r.lookupKey(tpl.Integrity.KeyId)
	if err != nil {
		return err
	}

	if key.expiredAt != nil && key.expiredAt.Before(tpl.Integrity.HashedAt.AsTime()) {
		return fmt.Errorf("relationship %s is signed by an expired key", tpl)
	}

	computedHash, err := computeRelationshipHash(tpl, key)
	if err != nil {
		return err
	}

	if !hmac.Equal(computedHash, hashWithoutByte) {
		str, err := tuple.String(tpl)
		if err != nil {
			return err
		}

		return fmt.Errorf("relationship %s has invalid integrity hash", str)
	}

	// NOTE: The caller expects the integrity to be nil, so the proxy sets it to nil here.
	tpl.Integrity = nil
	return nil
}

func (r *relationshipIntegrityProxy) Watch(ctx context.Context, afterRevision datastore.Revision, options datastore.WatchOptions) (<-chan *datastore.RevisionChanges, <-chan error) {
	resultsChan, errChan := r.ds.Watch(ctx, afterRevision, options)
	checkedResultsChan := make(chan *datastore.RevisionChanges)
	checkedErrChan := make(chan error, 1)

	go func() {
		defer close(checkedResultsChan)
		defer close(checkedErrChan)

		for {
			select {
			case result := <-resultsChan:
				for _, rel := range result.RelationshipChanges {
					if rel.Operation != corev1.RelationTupleUpdate_DELETE {
						err := r.validateRelationTuple(rel.Tuple)
						if err != nil {
							checkedErrChan <- err
							return
						}
					}
				}
				checkedResultsChan <- result

			case err := <-errChan:
				checkedErrChan <- err
				return
			}
		}
	}()

	return checkedResultsChan, checkedErrChan
}

func (r *relationshipIntegrityProxy) Unwrap() datastore.Datastore {
	return r.ds
}

type relationshipIntegrityReader struct {
	parent  *relationshipIntegrityProxy
	wrapped datastore.Reader
}

func (r relationshipIntegrityReader) QueryRelationships(ctx context.Context, filter datastore.RelationshipsFilter, options ...options.QueryOptionsOption) (datastore.RelationshipIterator, error) {
	it, err := r.wrapped.QueryRelationships(ctx, filter, options...)
	if err != nil {
		return nil, err
	}

	return &relationshipIntegrityIterator{
		parent:  r,
		wrapped: it,
	}, nil
}

func (r relationshipIntegrityReader) ReverseQueryRelationships(ctx context.Context, subjectsFilter datastore.SubjectsFilter, options ...options.ReverseQueryOptionsOption) (datastore.RelationshipIterator, error) {
	it, err := r.wrapped.ReverseQueryRelationships(ctx, subjectsFilter, options...)
	if err != nil {
		return nil, err
	}

	return &relationshipIntegrityIterator{
		parent:  r,
		wrapped: it,
	}, nil
}

func (r relationshipIntegrityReader) CountRelationships(ctx context.Context, name string) (int, error) {
	return r.wrapped.CountRelationships(ctx, name)
}

func (r relationshipIntegrityReader) ListAllCaveats(ctx context.Context) ([]datastore.RevisionedDefinition[*corev1.CaveatDefinition], error) {
	return r.wrapped.ListAllCaveats(ctx)
}

func (r relationshipIntegrityReader) ListAllNamespaces(ctx context.Context) ([]datastore.RevisionedDefinition[*corev1.NamespaceDefinition], error) {
	return r.wrapped.ListAllNamespaces(ctx)
}

func (r relationshipIntegrityReader) LookupCaveatsWithNames(ctx context.Context, names []string) ([]datastore.RevisionedDefinition[*corev1.CaveatDefinition], error) {
	return r.wrapped.LookupCaveatsWithNames(ctx, names)
}

func (r relationshipIntegrityReader) LookupCounters(ctx context.Context) ([]datastore.RelationshipCounter, error) {
	return r.wrapped.LookupCounters(ctx)
}

func (r relationshipIntegrityReader) LookupNamespacesWithNames(ctx context.Context, nsNames []string) ([]datastore.RevisionedDefinition[*corev1.NamespaceDefinition], error) {
	return r.wrapped.LookupNamespacesWithNames(ctx, nsNames)
}

func (r relationshipIntegrityReader) ReadCaveatByName(ctx context.Context, name string) (caveat *corev1.CaveatDefinition, lastWritten datastore.Revision, err error) {
	return r.wrapped.ReadCaveatByName(ctx, name)
}

func (r relationshipIntegrityReader) ReadNamespaceByName(ctx context.Context, nsName string) (ns *corev1.NamespaceDefinition, lastWritten datastore.Revision, err error) {
	return r.wrapped.ReadNamespaceByName(ctx, nsName)
}

type relationshipIntegrityIterator struct {
	parent  relationshipIntegrityReader
	wrapped datastore.RelationshipIterator
	err     error
}

func (r *relationshipIntegrityIterator) Close() {
	r.wrapped.Close()
}

func (r *relationshipIntegrityIterator) Cursor() (options.Cursor, error) {
	return r.wrapped.Cursor()
}

func (r *relationshipIntegrityIterator) Err() error {
	if r.err != nil {
		return r.err
	}

	return r.wrapped.Err()
}

func (r *relationshipIntegrityIterator) Next() *corev1.RelationTuple {
	tpl := r.wrapped.Next()
	if tpl == nil {
		return nil
	}

	err := r.parent.parent.validateRelationTuple(tpl)
	if err != nil {
		r.err = err
		return nil
	}

	return tpl
}

type relationshipIntegrityTx struct {
	datastore.ReadWriteTransaction

	parent *relationshipIntegrityProxy
}

func (r *relationshipIntegrityTx) WriteRelationships(
	ctx context.Context,
	mutations []*corev1.RelationTupleUpdate,
) error {
	// Add integrity data to the relationships.
	key := r.parent.primaryKey
	hashedAt := timestamppb.Now()

	updated := make([]*corev1.RelationTupleUpdate, 0, len(mutations))
	for _, mutation := range mutations {
		if mutation.Tuple.Integrity != nil {
			return spiceerrors.MustBugf("relationship %s already has integrity data", mutation.Tuple)
		}

		hash, err := computeRelationshipHash(mutation.Tuple, key)
		if err != nil {
			return err
		}

		// NOTE: Callers expect to be able to reuse the tuple, so we need to clone it.
		cloned := mutation.CloneVT()
		cloned.Tuple.Integrity = &corev1.RelationshipIntegrity{
			HashedAt: hashedAt,
			Hash:     append([]byte{versionByte}, hash...),
			KeyId:    key.keyID,
		}
		updated = append(updated, cloned)
	}

	return r.ReadWriteTransaction.WriteRelationships(ctx, updated)
}

func (r *relationshipIntegrityTx) BulkLoad(
	ctx context.Context,
	iter datastore.BulkWriteRelationshipSource,
) (uint64, error) {
	wrapped := &integrityAddingBulkLoadInterator{iter, r.parent}
	return r.ReadWriteTransaction.BulkLoad(ctx, wrapped)
}

type integrityAddingBulkLoadInterator struct {
	wrapped datastore.BulkWriteRelationshipSource
	parent  *relationshipIntegrityProxy
}

func (w integrityAddingBulkLoadInterator) Next(ctx context.Context) (*corev1.RelationTuple, error) {
	tpl, err := w.wrapped.Next(ctx)
	if err != nil {
		return nil, err
	}

	if tpl == nil {
		return nil, nil
	}

	key := w.parent.primaryKey
	hashedAt := timestamppb.Now()

	hash, err := computeRelationshipHash(tpl, key)
	if err != nil {
		return nil, err
	}

	if tpl.Integrity != nil {
		return nil, spiceerrors.MustBugf("relationship %s already has integrity data", tpl)
	}

	tpl.Integrity = &corev1.RelationshipIntegrity{
		HashedAt: hashedAt,
		Hash:     append([]byte{versionByte}, hash...),
		KeyId:    key.keyID,
	}

	return tpl, nil
}
