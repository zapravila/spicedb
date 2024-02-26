package migrations

import (
	"context"

	"github.com/jackc/pgx/v5"
)

var sqlStmts = []string{
	`ALTER TABLE relation_tuple
		ADD COLUMN r_id bigserial NOT NULL`,
	`ALTER TABLE relation_tuple ADD CONSTRAINT r_id_unique UNIQUE (r_id)`,
	`CREATE TABLE public.relation_tuple_add_info (
		id bigserial NOT NULL,
		relation_tuple_id int8 NULL,
		description varchar(4096) NULL,
		"comment" text NULL,
		CONSTRAINT relation_tuple_add_info_pk PRIMARY KEY (id),
		CONSTRAINT relation_tuple_add_info_relation_tuple_fk FOREIGN KEY (relation_tuple_id) REFERENCES public.relation_tuple(r_id)		
	)`,
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
