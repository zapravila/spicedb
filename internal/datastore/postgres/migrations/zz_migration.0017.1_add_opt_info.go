package migrations

import (
	"context"

	"github.com/jackc/pgx/v5"
)

var sqlStmts = []string{
	`ALTER TABLE relation_tuple
		ADD COLUMN description varchar(4096) NULL,
		ADD COLUMN comment text NULL`,
}

func init() {
	if err := DatabaseMigrations.Register("add_opt_info", "add-rel-by-alive-resource-relation-subject",
		noNonatomicMigration,
		func(ctx context.Context, tx pgx.Tx) error {
			for _, stmt := range sqlStmts {
				if _, err := tx.Exec(ctx, stmt); err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
		panic("failed to register migration: " + err.Error())
	}
}
