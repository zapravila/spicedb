package proxy

import (
	"testing"
	"time"

	"github.com/zapravila/spicedb/internal/datastore/memdb"
	"github.com/zapravila/spicedb/pkg/datastore"
	"github.com/zapravila/spicedb/pkg/datastore/test"
)

type observableTest struct{}

func (obs observableTest) New(revisionQuantization, _, gcWindow time.Duration, watchBufferLength uint16) (datastore.Datastore, error) {
	db, err := memdb.NewMemdbDatastore(watchBufferLength, revisionQuantization, gcWindow)
	if err != nil {
		return nil, err
	}
	return NewObservableDatastoreProxy(db), nil
}

func TestObservableProxy(t *testing.T) {
	test.All(t, observableTest{}, true)
}

func (p *observableProxy) ExampleRetryableError() error {
	return memdb.ErrSerialization
}
