package db

import (
	"fmt"

	log "github.com/inconshreveable/log15"
	"github.com/jmoiron/sqlx"
	"github.com/knadh/goyesql"
	_ "github.com/lib/pq"
)

type DB struct {
	db      *sqlx.DB
	schema  goyesql.Queries
	queries goyesql.Queries
}

func NewDBClient(dbUrl string) (*DB, error) {
	db, err := sqlx.Connect("postgres", dbUrl)
	if err != nil {
		return nil, err
	}

	schema := goyesql.MustParseFile("./sql/schema.sql")
	queries := goyesql.MustParseFile("./sql/queries.sql")

	return &DB{
		db:      db,
		schema:  schema,
		queries: queries,
	}, nil
}

func (d *DB) ApplySchema() error {
	sq, ok := d.schema["schema"]
	if !ok {
		log.Error("schema not found")
		return fmt.Errorf("schemanot not found")
	}

	_, err := d.db.Exec(sq.Query)
	if err != nil {
		return err
	}

	return nil
}
