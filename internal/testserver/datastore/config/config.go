//go:build docker
// +build docker

package config

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	testdatastore "github.com/zapravila/spicedb/internal/testserver/datastore"
	dsconfig "github.com/zapravila/spicedb/pkg/cmd/datastore"
	"github.com/zapravila/spicedb/pkg/datastore"
)

// DatastoreConfigInitFunc returns a InitFunc that constructs a ds
// with the top-level cmd/datastore machinery.
// It can't be used everywhere due to import cycles, but makes it easy to write
// an independent test with CLI-like config where possible.
func DatastoreConfigInitFunc(t testing.TB, options ...dsconfig.ConfigOption) testdatastore.InitFunc {
	return func(engine, uri string) datastore.Datastore {
		ds, err := dsconfig.NewDatastore(context.Background(),
			append(options,
				dsconfig.WithEngine(engine),
				dsconfig.WithEnableDatastoreMetrics(false),
				dsconfig.WithURI(uri),
			)...)
		require.NoError(t, err)
		return ds
	}
}
