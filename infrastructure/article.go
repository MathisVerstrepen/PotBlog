package infrastructure

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

type Metadata struct {
	Name        string
	Title       string
	Description string
	Date        string
	Tags        []string
	Author      string
}

func (db *DB) GetArticle(name string) (Metadata, error) {
	var metadata Metadata
	query := `
        SELECT articles.name, title, description, date, group_concat(tag), author
        FROM articles
		LEFT JOIN tags ON articles.name = tags.name
        WHERE articles.name = ?;
    `

	conn, err := Database.Pool.Take(context.Background())
	if err != nil {
		return metadata, err
	}
	defer Database.Pool.Put(conn)

	err = sqlitex.Execute(conn, strings.TrimSpace(query), &sqlitex.ExecOptions{
		Args: []interface{}{name},
		ResultFunc: func(stmt *sqlite.Stmt) error {
			metadata.Name = stmt.ColumnText(0)
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

type SortAndFilter struct {
	SortBy   string
	FilterBy []string
}

func (sf SortAndFilter) Default() SortAndFilter {
	sf.SortBy = "date_desc"
	return sf
}

func (sf SortAndFilter) HasFilter() bool {
	return len(sf.FilterBy) > 0
}

func (sf SortAndFilter) OrderClause() (string, string) {
	switch sf.SortBy {
	case "date_asc":
		return "date", "ASC"
	case "date_desc":
		return "date", "DESC"
	case "title_asc":
		return "title", "ASC"
	case "title_desc":
		return "title", "DESC"
	default:
		return "date", "DESC"
	}
}

func (sf SortAndFilter) FilterClause(dbTags []string) string {
	var validTags []string
	for _, tag := range sf.FilterBy {
		if !slices.Contains(dbTags, tag) {
			continue
		}

		validTags = append(validTags, fmt.Sprintf("'%s'", tag))
	}

	if len(validTags) == 0 {
		return "1"
	}

	return fmt.Sprintf("tag IN (%s)", strings.Join(validTags, ","))
}

func (db *DB) GetArticles(sorter SortAndFilter) ([]Metadata, error) {
	var metadata Metadata
	var articles []Metadata
	var filterClause string

	if sorter.HasFilter() {
		dbTags, err := db.GetArticlesTags()
		if err != nil {
			fmt.Println("error on get articles tags:", err)
			return articles, err
		}
		filterClause = sorter.FilterClause(dbTags)
	} else {
		filterClause = "1"
	}

	orderColumn, orderDirection := sorter.OrderClause()
	query := fmt.Sprintf(`
		SELECT articles.name, title, description, date, group_concat(tag) as tags, author
		FROM articles
		LEFT JOIN tags ON articles.name = tags.name
		WHERE %s
		GROUP BY articles.name
        ORDER BY %s %s;
    `, filterClause, orderColumn, orderDirection)

	conn, err := Database.Pool.Take(context.Background())
	if err != nil {
		return articles, err
	}
	defer Database.Pool.Put(conn)

	err = sqlitex.Execute(conn, strings.TrimSpace(query), &sqlitex.ExecOptions{
		ResultFunc: func(stmt *sqlite.Stmt) error {
			metadata.Name = stmt.ColumnText(0)
			metadata.Title = stmt.ColumnText(1)
			metadata.Description = stmt.ColumnText(2)
			metadata.Date = stmt.ColumnText(3)
			metadata.Tags = strings.Split(stmt.ColumnText(4), ",")
			metadata.Author = stmt.ColumnText(5)
			articles = append(articles, metadata)
			return nil
		},
	})

	if err != nil {
		return articles, err
	}

	return articles, nil
}

func (db *DB) SaveArticle(metadata Metadata) error {
	conn, err := Database.Pool.Take(context.Background())
	if err != nil {
		return err
	}
	defer Database.Pool.Put(conn)

	query := `
		INSERT OR REPLACE INTO articles (name, title, description, date, author)
		VALUES ($name, $title, $description, $date, $author);
	`

	stmt := conn.Prep(strings.TrimSpace(query))
	stmt.SetText("$name", metadata.Name)
	stmt.SetText("$title", metadata.Title)
	stmt.SetText("$description", metadata.Description)
	stmt.SetText("$date", metadata.Date)
	stmt.SetText("$author", metadata.Author)

	_, err = stmt.Step()
	if err != nil {
		return err
	}

	query = `
		INSERT OR REPLACE INTO tags (name, tag)
		VALUES ($name, $tag);
	`

	stmt = conn.Prep(strings.TrimSpace(query))
	for _, tag := range metadata.Tags {
		stmt.SetText("$name", metadata.Name)
		stmt.SetText("$tag", tag)

		_, err = stmt.Step()
		if err != nil {
			fmt.Println("error on insert tag:", err)
			return err
		}

		stmt.Reset()
	}

	return nil
}

func (db *DB) GetArticlesTags() ([]string, error) {
	var tags []string
	query := `
		SELECT DISTINCT tag
		FROM tags;
	`

	conn, err := Database.Pool.Take(context.Background())
	if err != nil {
		return tags, err
	}
	defer Database.Pool.Put(conn)

	err = sqlitex.Execute(conn, strings.TrimSpace(query), &sqlitex.ExecOptions{
		ResultFunc: func(stmt *sqlite.Stmt) error {
			tags = append(tags, stmt.ColumnText(0))
			return nil
		},
	})

	if err != nil {
		fmt.Println("error on get tags:", err)
		return tags, err
	}

	return tags, nil
}
