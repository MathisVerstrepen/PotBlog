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

func InitializeDatabase(path string) error {
	dbpool, err := sqlitex.NewPool(path, sqlitex.PoolOptions{
		PoolSize: 10,
	})
	if err != nil {
		return err
	}

	Database = &DB{Pool: dbpool}

	err = createDatabaseSchema()
	if err != nil {
		return err
	}

	return nil
}

func createDatabaseSchema() error {
	createArticles := `
        CREATE TABLE IF NOT EXISTS articles (
            name TEXT PRIMARY KEY,
            title TEXT,
            description TEXT,
            date TEXT,
            author TEXT
        );`

	createTags := `
        CREATE TABLE IF NOT EXISTS tags (
            name TEXT,
            tag TEXT,
            FOREIGN KEY(name) REFERENCES articles(name),
            PRIMARY KEY(name, tag)
        );`

	conn, err := Database.Pool.Take(context.Background())
	if err != nil {
		return err
	}
	defer Database.Pool.Put(conn)

	if err := sqlitex.Execute(conn, strings.TrimSpace(createArticles), nil); err != nil {
		return err
	}

	return sqlitex.Execute(conn, strings.TrimSpace(createTags), nil)
}
