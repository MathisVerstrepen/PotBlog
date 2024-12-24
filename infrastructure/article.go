package infrastructure

import (
	"fmt"
	"strings"
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
		WHERE name = $name;
	`

	stmt := db.Conn.Prep(strings.TrimSpace(query))
	stmt.SetText("$name", name)

	hasRow, err := stmt.Step()
	if err != nil {
		return metadata, err
	}

	if !hasRow {
		return metadata, fmt.Errorf("article not found")
	}

	metadata.Title = stmt.GetText("title")
	metadata.Description = stmt.GetText("description")
	metadata.Date = stmt.GetText("date")
	metadata.Tags = strings.Split(stmt.GetText("tags"), ",")
	metadata.Author = stmt.GetText("author")

	return metadata, err
}

func (db *DB) SaveArticle(metadata Metadata, name string) error {
	query := `
		INSERT OR REPLACE INTO articles (name, title, description, date, tags, author)
		VALUES ($name, $title, $description, $date, $tags, $author);
	`

	stmt := db.Conn.Prep(strings.TrimSpace(query))
	stmt.SetText("$name", name)
	stmt.SetText("$title", metadata.Title)
	stmt.SetText("$description", metadata.Description)
	stmt.SetText("$date", metadata.Date)
	stmt.SetText("$tags", strings.Join(metadata.Tags, ","))
	stmt.SetText("$author", metadata.Author)

	_, err := stmt.Step()
	if err != nil {
		return err
	}

	return nil
}
