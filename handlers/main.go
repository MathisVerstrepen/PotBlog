package handlers

import (
	"fmt"
	"log"

	"github.com/a-h/templ"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"

	"potblog/infrastructure"
	"potblog/services"
)

func Init() {
	fmt.Println("Startup sequence starting...")

	err := godotenv.Load(".env")
	if err != nil {
		log.Println("WARNING : Failed to load .env file")
	}

	err = infrastructure.Open("potblog.db")
	if err != nil {
		log.Fatalf("Failed to open database: %s", err)
	}

	err = generateStaticArticles()
	if err != nil {
		log.Fatalf("Failed to generate static articles: %s", err)
	}

	fmt.Println("Startup sequence done.")
}

func generateStaticArticles() error {
	articles, err := services.GetArticles()

	if err != nil {
		return err
	}

	for _, article := range articles {
		md := services.ReadMarkdownFile(fmt.Sprintf("assets/articles/markdown/%s", article))
		articleData, err := services.MarkdownToHTML(&md)
		if err != nil {
			return err
		}

		err = services.SaveArticle(article, articleData.RawHTML)
		if err != nil {
			return err
		}

		articleName := article[:len(article)-3]
		err = infrastructure.Database.SaveArticle(articleData.Metadata, articleName)
		if err != nil {
			return err
		}
	}

	return nil
}

func Render(ctx echo.Context, statusCode int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := t.Render(ctx.Request().Context(), buf); err != nil {
		return err
	}

	return ctx.HTML(statusCode, buf.String())
}
