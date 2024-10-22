package v1

import (
	"context"

	grpcvalidate "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/validator"
	v1 "github.com/zapravila/authzed-go/proto/authzed/api/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	log "github.com/zapravila/spicedb/internal/logging"
	"github.com/zapravila/spicedb/internal/middleware"
	datastoremw "github.com/zapravila/spicedb/internal/middleware/datastore"
	"github.com/zapravila/spicedb/internal/middleware/usagemetrics"
	"github.com/zapravila/spicedb/internal/services/shared"
	"github.com/zapravila/spicedb/pkg/datastore"
	"github.com/zapravila/spicedb/pkg/genutil"
	dispatchv1 "github.com/zapravila/spicedb/pkg/proto/dispatch/v1"
	"github.com/zapravila/spicedb/pkg/schemadsl/compiler"
	"github.com/zapravila/spicedb/pkg/schemadsl/generator"
	"github.com/zapravila/spicedb/pkg/schemadsl/input"
	"github.com/zapravila/spicedb/pkg/zedtoken"
)

// NewSchemaServer creates a SchemaServiceServer instance.
func NewSchemaServer(additiveOnly bool) v1.SchemaServiceServer {
	return &schemaServer{
		WithServiceSpecificInterceptors: shared.WithServiceSpecificInterceptors{
			Unary: middleware.ChainUnaryServer(
				grpcvalidate.UnaryServerInterceptor(),
				usagemetrics.UnaryServerInterceptor(),
			),
			Stream: middleware.ChainStreamServer(
				grpcvalidate.StreamServerInterceptor(),
				usagemetrics.StreamServerInterceptor(),
			),
		},
		additiveOnly: additiveOnly,
	}
}

type schemaServer struct {
	v1.UnimplementedSchemaServiceServer
	shared.WithServiceSpecificInterceptors

	additiveOnly bool
}

func (ss *schemaServer) rewriteError(ctx context.Context, err error) error {
	return shared.RewriteError(ctx, err, nil)
}

func (ss *schemaServer) ReadSchema(ctx context.Context, _ *v1.ReadSchemaRequest) (*v1.ReadSchemaResponse, error) {
	// Schema is always read from the head revision.
	ds := datastoremw.MustFromContext(ctx)
	headRevision, err := ds.HeadRevision(ctx)
	if err != nil {
		return nil, ss.rewriteError(ctx, err)
	}

	reader := ds.SnapshotReader(headRevision)

	nsDefs, err := reader.ListAllNamespaces(ctx)
	if err != nil {
		return nil, ss.rewriteError(ctx, err)
	}

	caveatDefs, err := reader.ListAllCaveats(ctx)
	if err != nil {
		return nil, ss.rewriteError(ctx, err)
	}

	if len(nsDefs) == 0 {
		return nil, status.Errorf(codes.NotFound, "No schema has been defined; please call WriteSchema to start")
	}

	schemaDefinitions := make([]compiler.SchemaDefinition, 0, len(nsDefs)+len(caveatDefs))
	for _, caveatDef := range caveatDefs {
		schemaDefinitions = append(schemaDefinitions, caveatDef.Definition)
	}

	for _, nsDef := range nsDefs {
		schemaDefinitions = append(schemaDefinitions, nsDef.Definition)
	}

	schemaText, _, err := generator.GenerateSchema(schemaDefinitions)
	if err != nil {
		return nil, ss.rewriteError(ctx, err)
	}

	dispatchCount, err := genutil.EnsureUInt32(len(nsDefs) + len(caveatDefs))
	if err != nil {
		return nil, ss.rewriteError(ctx, err)
	}

	usagemetrics.SetInContext(ctx, &dispatchv1.ResponseMeta{
		DispatchCount: dispatchCount,
	})

	return &v1.ReadSchemaResponse{
		SchemaText: schemaText,
		ReadAt:     zedtoken.MustNewFromRevision(headRevision),
	}, nil
}

func (ss *schemaServer) WriteSchema(ctx context.Context, in *v1.WriteSchemaRequest) (*v1.WriteSchemaResponse, error) {
	log.Ctx(ctx).Trace().Str("schema", in.GetSchema()).Msg("requested Schema to be written")

	ds := datastoremw.MustFromContext(ctx)

	// Compile the schema into the namespace definitions.
	compiled, err := compiler.Compile(compiler.InputSchema{
		Source:       input.Source("schema"),
		SchemaString: in.GetSchema(),
	}, compiler.AllowUnprefixedObjectType())
	if err != nil {
		return nil, ss.rewriteError(ctx, err)
	}
	log.Ctx(ctx).Trace().Int("objectDefinitions", len(compiled.ObjectDefinitions)).Int("caveatDefinitions", len(compiled.CaveatDefinitions)).Msg("compiled namespace definitions")

	// Do as much validation as we can before talking to the datastore.
	validated, err := shared.ValidateSchemaChanges(ctx, compiled, ss.additiveOnly)
	if err != nil {
		return nil, ss.rewriteError(ctx, err)
	}

	// Update the schema.
	revision, err := ds.ReadWriteTx(ctx, func(ctx context.Context, rwt datastore.ReadWriteTransaction) error {
		applied, err := shared.ApplySchemaChanges(ctx, rwt, validated)
		if err != nil {
			return err
		}

		dispatchCount, err := genutil.EnsureUInt32(applied.TotalOperationCount)
		if err != nil {
			return err
		}

		usagemetrics.SetInContext(ctx, &dispatchv1.ResponseMeta{
			DispatchCount: dispatchCount,
		})
		return nil
	})
	if err != nil {
		return nil, ss.rewriteError(ctx, err)
	}

	return &v1.WriteSchemaResponse{
		WrittenAt: zedtoken.MustNewFromRevision(revision),
	}, nil
}
