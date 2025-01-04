package handlers

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"

	comp "potblog/components"
	"potblog/infrastructure"
	"potblog/services"
)

func ServeArticle(c echo.Context) error {
	article := c.Param("article")
	language := c.Param("language")

	fmt.Println("language:", language)

	html, err := services.GetArticle(article)
	if err != nil {
		return Render(c, http.StatusOK, comp.Root(comp.ArticleNotFound(), fmt.Sprintf("%s - Not found", article)))
	}

	metadata, err := infrastructure.Database.GetArticle(article)
	if err != nil {
		return Render(c, http.StatusOK, comp.Root(comp.ArticleNotFound(), fmt.Sprintf("%s - Not found", article)))
	}

	return Render(c, http.StatusOK, comp.Root(comp.Article(metadata, html), article))
}

func ServeArticles(c echo.Context) error {
	language := c.Param("language")
	fmt.Println("language:", language)

	articles, err := infrastructure.Database.GetArticles(infrastructure.SortAndFilter{}.Default())
	if err != nil {
		return Render(c, http.StatusOK, comp.Root(comp.ServerError(), "Server Error"))
	}

	return Render(c, http.StatusOK, comp.Root(comp.Articles(articles), language))
}

func formDataMapper(formData url.Values) infrastructure.SortAndFilter {
	var sortAndFilter infrastructure.SortAndFilter
	for key, value := range formData {
		if key == "sort" {
			sortAndFilter.SortBy = value[0]
		} else {
			sortAndFilter.FilterBy = append(sortAndFilter.FilterBy, key)
		}
	}

	return sortAndFilter
}

func ServeArticlesSortAndFilter(c echo.Context) error {
	language := c.Param("language")
	fmt.Println("language:", language)

	formData, err := c.FormParams()
	if err != nil {
		return Render(c, http.StatusOK, comp.Root(comp.ServerError(), "Server Error"))
	}

	sortAndFilter := formDataMapper(formData)

	articles, err := infrastructure.Database.GetArticles(sortAndFilter)
	if err != nil {
		return Render(c, http.StatusOK, comp.Root(comp.ServerError(), "Server Error"))
	}

	return Render(c, http.StatusOK, comp.ArticleGrid(articles))
}
