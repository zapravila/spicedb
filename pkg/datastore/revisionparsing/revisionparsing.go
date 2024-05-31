package revisionparsing

import (
	"github.com/zapravila/spicedb/internal/datastore/memdb"
	"github.com/zapravila/spicedb/internal/datastore/postgres"
	"github.com/zapravila/spicedb/internal/datastore/revisions"
	"github.com/zapravila/spicedb/pkg/datastore"
)

// ParsingFunc is a function that can parse a string into a revision.
type ParsingFunc func(revisionStr string) (rev datastore.Revision, err error)

// ParseRevisionStringByDatastoreEngineID defines a map from datastore engine ID to its associated
// revision parsing function.
var ParseRevisionStringByDatastoreEngineID = map[string]ParsingFunc{
	memdb.Engine:    ParsingFunc(memdb.ParseRevisionString),
	postgres.Engine: ParsingFunc(postgres.ParseRevisionString),
}

// MustParseRevisionForTest is a convenience ParsingFunc that can be used in tests and panics when parsing an error.
func MustParseRevisionForTest(revisionStr string) (rev datastore.Revision) {
	rev, err := testParser(revisionStr)
	if err != nil {
		panic(err)
	}

	return rev
}

var testParser = revisions.RevisionParser(revisions.HybridLogicalClock)
