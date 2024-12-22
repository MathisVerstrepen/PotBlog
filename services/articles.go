package services

import (
	"os"
	"path"
	"strings"
)

const (
	MarkdownArticlesPath = "assets/articles/markdown"
	HTMLArticlesPath     = "assets/articles/html"
)

func GetArticles() ([]string, error) {
	articles := []string{}

	files, err := os.ReadDir(path.Join(Root, MarkdownArticlesPath))
	if err != nil {
		return articles, err
	}

	for _, file := range files {
		articles = append(articles, file.Name())
	}

	return articles, nil
}

func GetArticle(article string) (string, error) {
	article += ".html"
	file, err := os.ReadFile(path.Join(Root, HTMLArticlesPath, article))
	if err != nil {
		return "", err
	}

	return string(file), nil
}

func SaveArticle(article string, content string) error {
	article = strings.Replace(article, ".md", ".html", 1)

	file, err := os.Create(path.Join(Root, HTMLArticlesPath, article))
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}
