package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	comp "potblog/components"
	"potblog/services"
)

func ServeArticle(c echo.Context) error {
	article := c.Param("article")
	language := c.Param("language")

	fmt.Println("language:", language)

	html, err := services.GetArticle(article)
	if err != nil {
		return c.String(http.StatusNotFound, "Article not found")
	}

	return Render(c, http.StatusOK, comp.Root(comp.Article(html), article))
}
