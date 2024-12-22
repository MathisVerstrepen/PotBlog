package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	comp "potblog/components"
	"potblog/services"
)

func ServeArticle(c echo.Context) error {
	article := c.Param("article")

	html, err := services.GetArticle(article)
	if err != nil {
		return c.String(http.StatusNotFound, "Article not found")
	}

	return Render(c, http.StatusOK, comp.Root(comp.Article(html), article))
}
