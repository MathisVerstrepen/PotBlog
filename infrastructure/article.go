package infrastructure

import (
	"context"
	"strings"

	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

type Metadata struct {
	Title       string
	Description string
	Date        string
	Tags        []string
	Author      string
}

func (db *DB) GetArticle(name string) (Metadata, error) {
	var metadata Metadata
	query := `
        SELECT name, title, description, date, tags, author
        FROM articles
        WHERE name = ?;
    `

	conn, err := Database.Pool.Take(context.Background())
	if err != nil {
		return metadata, err
	}
	defer Database.Pool.Put(conn)

	err = sqlitex.Execute(conn, strings.TrimSpace(query), &sqlitex.ExecOptions{
		Args: []interface{}{name},
		ResultFunc: func(stmt *sqlite.Stmt) error {
			metadata.Title = stmt.ColumnText(1)
			metadata.Description = stmt.ColumnText(2)
			metadata.Date = stmt.ColumnText(3)
			metadata.Tags = strings.Split(stmt.ColumnText(4), ",")
			metadata.Author = stmt.ColumnText(5)
			return nil
		},
	})

	if err != nil {
		return metadata, err
	}

	return metadata, nil
}

func (db *DB) SaveArticle(metadata Metadata, name string) error {
	conn, err := Database.Pool.Take(context.Background())
	if err != nil {
		return err
	}
	defer Database.Pool.Put(conn)

	query := `
		INSERT OR REPLACE INTO articles (name, title, description, date, tags, author)
		VALUES ($name, $title, $description, $date, $tags, $author);
	`

	stmt := conn.Prep(strings.TrimSpace(query))
	stmt.SetText("$name", name)
	stmt.SetText("$title", metadata.Title)
	stmt.SetText("$description", metadata.Description)
	stmt.SetText("$date", metadata.Date)
	stmt.SetText("$tags", strings.Join(metadata.Tags, ","))
	stmt.SetText("$author", metadata.Author)

	_, err = stmt.Step()
	if err != nil {
		return err
	}

	return nil
}
