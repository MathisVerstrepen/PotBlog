package infrastructure

import (
	"context"
	"strings"

	"zombiezen.com/go/sqlite/sqlitex"
)

type DB struct {
	Pool *sqlitex.Pool
}

var Database *DB

func Open(path string) error {
	dbpool, err := sqlitex.NewPool(path, sqlitex.PoolOptions{
		PoolSize: 10,
	})
	if err != nil {
		return err
	}

	Database = &DB{Pool: dbpool}

	err = initTables()
	if err != nil {
		return err
	}

	return nil
}

func initTables() error {
	query := `
        CREATE TABLE IF NOT EXISTS articles (
			name TEXT PRIMARY KEY,
            title TEXT,
			description TEXT,
			date TEXT,
			tags TEXT,
			author TEXT
        );`

	conn, err := Database.Pool.Take(context.Background())
	if err != nil {
		return err
	}
	defer Database.Pool.Put(conn)

	return sqlitex.Execute(conn, strings.TrimSpace(query), nil)
}
