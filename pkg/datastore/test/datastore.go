package test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	core "github.com/zapravila/spicedb/pkg/proto/core/v1"

	"github.com/zapravila/spicedb/pkg/datastore"
	"github.com/zapravila/spicedb/pkg/namespace"
)

// veryLargeGCWindow is a very large time duration, which when passed to a constructor should
// effectively disable garbage collection.
const (
	veryLargeGCWindow   = 90000 * time.Second
	veryLargeGCInterval = 90000 * time.Second
)

// DatastoreTester provides a generic datastore suite a means of initializing
// a particular datastore.
type DatastoreTester interface {
	// New creates a new datastore instance for a single test.
	New(revisionQuantization, gcInterval, gcWindow time.Duration, watchBufferLength uint16) (datastore.Datastore, error)
}

type DatastoreTesterFunc func(revisionQuantization, gcInterval, gcWindow time.Duration, watchBufferLength uint16) (datastore.Datastore, error)

func (f DatastoreTesterFunc) New(revisionQuantization, gcInterval, gcWindow time.Duration, watchBufferLength uint16) (datastore.Datastore, error) {
	return f(revisionQuantization, gcInterval, gcWindow, watchBufferLength)
}

type TestableDatastore interface {
	datastore.Datastore

	ExampleRetryableError() error
}

type Categories map[string]struct{}

func (c Categories) GC() bool {
	_, ok := c[GCCategory]
	return ok
}

func (c Categories) Stats() bool {
	_, ok := c[StatsCategory]
	return ok
}

func (c Categories) Watch() bool {
	_, ok := c[WatchCategory]
	return ok
}

func (c Categories) WatchSchema() bool {
	_, ok := c[WatchSchemaCategory]
	return ok
}

func (c Categories) WatchCheckpoints() bool {
	_, ok := c[WatchCheckpointsCategory]
	return ok
}

var noException = Categories{}

const (
	GCCategory               = "GC"
	WatchCategory            = "Watch"
	WatchSchemaCategory      = "WatchSchema"
	WatchCheckpointsCategory = "WatchCheckpoints"
	StatsCategory            = "Stats"
)

func WithCategories(cats ...string) Categories {
	c := Categories{}
	for _, cat := range cats {
		c[cat] = struct{}{}
	}
	return c
}

func parallel(tester DatastoreTester, tt func(t *testing.T, tester DatastoreTester)) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		tt(t, tester)
	}
}

func serial(tester DatastoreTester, tt func(t *testing.T, tester DatastoreTester)) func(t *testing.T) {
	return func(t *testing.T) {
		tt(t, tester)
	}
}

// AllWithExceptions runs all generic datastore tests on a DatastoreTester, except
// those specified test categories
func AllWithExceptions(t *testing.T, tester DatastoreTester, except Categories, concurrent bool) {
	runner := serial
	if concurrent {
		runner = parallel
	}

	t.Run("TestUseAfterClose", runner(tester, UseAfterCloseTest))

	t.Run("TestNamespaceNotFound", runner(tester, NamespaceNotFoundTest))
	t.Run("TestNamespaceWrite", runner(tester, NamespaceWriteTest))
	t.Run("TestNamespaceDelete", runner(tester, NamespaceDeleteTest))
	t.Run("TestNamespaceMultiDelete", runner(tester, NamespaceMultiDeleteTest))
	t.Run("TestEmptyNamespaceDelete", runner(tester, EmptyNamespaceDeleteTest))
	t.Run("TestStableNamespaceReadWrite", runner(tester, StableNamespaceReadWriteTest))

	t.Run("TestSimple", runner(tester, SimpleTest))
	t.Run("TestObjectIDs", runner(tester, ObjectIDsTest))
	t.Run("TestDeleteRelationships", runner(tester, DeleteRelationshipsTest))
	t.Run("TestDeleteNonExistant", runner(tester, DeleteNotExistantTest))
	t.Run("TestDeleteAlreadyDeleted", runner(tester, DeleteAlreadyDeletedTest))
	t.Run("TestRecreateRelationshipsAfterDeleteWithFilter", runner(tester, RecreateRelationshipsAfterDeleteWithFilter))
	t.Run("TestWriteDeleteWrite", runner(tester, WriteDeleteWriteTest))
	t.Run("TestCreateAlreadyExisting", runner(tester, CreateAlreadyExistingTest))
	t.Run("TestTouchAlreadyExistingWithoutCaveat", runner(tester, TouchAlreadyExistingTest))
	t.Run("TestCreateDeleteTouch", runner(tester, CreateDeleteTouchTest))
	t.Run("TestDeleteOneThousandIndividualInOneCall", runner(tester, DeleteOneThousandIndividualInOneCallTest))
	t.Run("TestCreateTouchDeleteTouch", runner(tester, CreateTouchDeleteTouchTest))
	t.Run("TestTouchAlreadyExistingCaveated", runner(tester, TouchAlreadyExistingCaveatedTest))
	t.Run("TestBulkDeleteRelationships", runner(tester, BulkDeleteRelationshipsTest))
	t.Run("TestDeleteCaveatedTuple", runner(tester, DeleteCaveatedTupleTest))
	t.Run("TestDeleteWithLimit", runner(tester, DeleteWithLimitTest))
	t.Run("TestQueryRelationshipsWithVariousFilters", runner(tester, QueryRelationshipsWithVariousFiltersTest))
	t.Run("TestDeleteRelationshipsWithVariousFilters", runner(tester, DeleteRelationshipsWithVariousFiltersTest))
	t.Run("TestTouchTypedAlreadyExistingWithoutCaveat", runner(tester, TypedTouchAlreadyExistingTest))
	t.Run("TestTouchTypedAlreadyExistingWithCaveat", runner(tester, TypedTouchAlreadyExistingWithCaveatTest))

	t.Run("TestMultipleReadsInRWT", runner(tester, MultipleReadsInRWTTest))
	t.Run("TestConcurrentWriteSerialization", runner(tester, ConcurrentWriteSerializationTest))

	t.Run("TestOrdering", runner(tester, OrderingTest))
	t.Run("TestLimit", runner(tester, LimitTest))
	t.Run("TestOrderedLimit", runner(tester, OrderedLimitTest))
	t.Run("TestResume", runner(tester, ResumeTest))
	t.Run("TestCursorErrors", runner(tester, CursorErrorsTest))
	t.Run("TestReverseQueryCursor", runner(tester, ReverseQueryCursorTest))

	t.Run("TestRevisionQuantization", runner(tester, RevisionQuantizationTest))
	t.Run("TestRevisionSerialization", runner(tester, RevisionSerializationTest))
	t.Run("TestSequentialRevisions", runner(tester, SequentialRevisionsTest))
	t.Run("TestConcurrentRevisions", runner(tester, ConcurrentRevisionsTest))
	t.Run("TestCheckRevisions", runner(tester, CheckRevisionsTest))

	if !except.GC() {
		t.Run("TestRevisionGC", runner(tester, RevisionGCTest))
		t.Run("TestInvalidReads", runner(tester, InvalidReadsTest))
	}

	t.Run("TestBulkUpload", runner(tester, BulkUploadTest))
	t.Run("TestBulkUploadErrors", runner(tester, BulkUploadErrorsTest))
	t.Run("TestBulkUploadAlreadyExistsError", runner(tester, BulkUploadAlreadyExistsErrorTest))
	t.Run("TestBulkUploadAlreadyExistsSameCallError", runner(tester, BulkUploadAlreadyExistsSameCallErrorTest))
	t.Run("BulkUploadEditCaveat", runner(tester, BulkUploadEditCaveat))

	if !except.Stats() {
		t.Run("TestStats", runner(tester, StatsTest))
	}

	t.Run("TestRetries", runner(tester, RetryTest))

	t.Run("TestCaveatNotFound", runner(tester, CaveatNotFoundTest))
	t.Run("TestWriteReadDeleteCaveat", runner(tester, WriteReadDeleteCaveatTest))
	t.Run("TestWriteCaveatedRelationship", runner(tester, WriteCaveatedRelationshipTest))
	t.Run("TestCaveatedRelationshipFilter", runner(tester, CaveatedRelationshipFilterTest))
	t.Run("TestCaveatSnapshotReads", runner(tester, CaveatSnapshotReadsTest))

	if !except.Watch() {
		t.Run("TestWatchBasic", runner(tester, WatchTest))
		t.Run("TestWatchCancel", runner(tester, WatchCancelTest))
		t.Run("TestCaveatedRelationshipWatch", runner(tester, CaveatedRelationshipWatchTest))
		t.Run("TestWatchWithTouch", runner(tester, WatchWithTouchTest))
		t.Run("TestWatchWithDelete", runner(tester, WatchWithDeleteTest))
		t.Run("TestWatchWithMetadata", runner(tester, WatchWithMetadataTest))
	}

	if !except.Watch() && !except.WatchSchema() {
		t.Run("TestWatchSchema", runner(tester, WatchSchemaTest))
		t.Run("TestWatchAll", runner(tester, WatchAllTest))
	}

	if !except.Watch() && !except.WatchCheckpoints() {
		t.Run("TestWatchCheckpoints", runner(tester, WatchCheckpointsTest))
	}

	t.Run("TestRelationshipCounters", runner(tester, RelationshipCountersTest))
	t.Run("TestUpdateRelationshipCounter", runner(tester, UpdateRelationshipCounterTest))
	t.Run("TestDeleteAllData", runner(tester, DeleteAllDataTest))
}

// All runs all generic datastore tests on a DatastoreTester.
func All(t *testing.T, tester DatastoreTester, concurrent bool) {
	AllWithExceptions(t, tester, noException, concurrent)
}

var testResourceNS = namespace.Namespace(
	testResourceNamespace,
	namespace.MustRelation(testReaderRelation, nil),
)

var testGroupNS = namespace.Namespace(
	testGroupNamespace,
	namespace.MustRelation(testMemberRelation, nil),
)

var testUserNS = namespace.Namespace(testUserNamespace)

func makeTestTuple(resourceID, userID string) *core.RelationTuple {
	return &core.RelationTuple{
		ResourceAndRelation: &core.ObjectAndRelation{
			Namespace: testResourceNamespace,
			ObjectId:  resourceID,
			Relation:  testReaderRelation,
		},
		Subject: &core.ObjectAndRelation{
			Namespace: testUserNamespace,
			ObjectId:  userID,
			Relation:  ellipsis,
		},
	}
}

func setupDatastore(ds datastore.Datastore, require *require.Assertions) datastore.Revision {
	ctx := context.Background()

	revision, err := ds.ReadWriteTx(ctx, func(ctx context.Context, rwt datastore.ReadWriteTransaction) error {
		return rwt.WriteNamespaces(ctx, testGroupNS, testResourceNS, testUserNS)
	})
	require.NoError(err)

	return revision
}
