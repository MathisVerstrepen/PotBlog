package infrastructure

import (
	"strings"

	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

type DB struct {
	*sqlite.Conn
}

var Database *DB

func Open(path string) error {
	conn, err := sqlite.OpenConn(path)
	if err != nil {
		return err
	}

	Database = &DB{conn}

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
	return sqlitex.ExecuteTransient(Database.Conn, strings.TrimSpace(query), nil)
}

func (db *DB) Close() error {
	return db.Conn.Close()
}
