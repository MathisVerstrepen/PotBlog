package services

import (
	"os"
	"path"
	"strings"
)

const (
	MarkdownArticlesPathBase = "assets/articles/markdown"
	HTMLArticlesPath         = "assets/articles/html"
)

func MarkdownArticlesPath() string {
	root := os.Getenv("ROOT_DIR")
	if strings.Compare(os.Getenv("ENV"), "dev") == 0 {
		return path.Join(root, MarkdownArticlesPathBase, "dev")
	}
	return path.Join(root, MarkdownArticlesPathBase, "prod")
}

func RetrieveLocalMdArticles() ([]string, error) {
	articles := []string{}

	files, err := os.ReadDir(MarkdownArticlesPath())
	if err != nil {
		return articles, err
	}

	for _, file := range files {
		if strings.Contains(file.Name(), ".md") {
			articles = append(articles, file.Name())
		}
	}

	return articles, nil
}

func RetriveLocalHtmlArticle(article string) (string, error) {
	article += ".html"
	file, err := os.ReadFile(path.Join(os.Getenv("ROOT_DIR"), HTMLArticlesPath, article))
	if err != nil {
		return "", err
	}

	return string(file), nil
}

func PersistHtmlArticle(article string, content string) error {
	article = strings.Replace(article, ".md", ".html", 1)

	file, err := os.Create(path.Join(os.Getenv("ROOT_DIR"), HTMLArticlesPath, article))
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
