package datastore_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zapravila/spicedb/internal/datastore/postgres"
	"github.com/zapravila/spicedb/pkg/datastore"
)

func createEngine(engineID string, uri string) error {
	ctx := context.Background()

	switch engineID {
	case "postgres":
		_, err := postgres.NewPostgresDatastore(ctx, uri)
		return err
	default:
		panic(fmt.Sprintf("missing create implementation for engine %s", engineID))
	}
}

func TestDatastoreURIErrors(t *testing.T) {
	tcs := map[string]string{
		"some-wrong-uri":                                  "wrong",
		"postgres://foo:bar:baz@someurl":                  "bar",
		"postgres://spicedb:somepassword":                 "somepassword",
		"postgres://spicedb:somepassword#@foo":            "somepassword",
		"username=foo password=somepassword dsn=whatever": "somepassword",
	}

	for _, engineID := range datastore.Engines {
		t.Run(engineID, func(t *testing.T) {
			for tc, check := range tcs {
				t.Run(tc, func(t *testing.T) {
					err := createEngine(engineID, tc)
					require.Error(t, err)
					require.NotContains(t, err.Error(), check)
				})
			}
		})
	}
}
