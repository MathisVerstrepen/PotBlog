package handlers

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"

	comp "potblog/components"
	infra "potblog/infrastructure"
	"potblog/services"
)

func ServeArticle(c echo.Context) error {
	article := c.Param("article")
	language := c.Param("language")

	fmt.Println("language:", language)

	articleHtmlContent, err := services.RetriveLocalHtmlArticle(article)
	if err != nil {
		return Render(c, http.StatusOK, comp.Root(comp.ArticleNotFound(), fmt.Sprintf("%s - Not found", article)))
	}

	articleMetadata, err := infra.Database.GetArticle(article)
	if err != nil {
		return Render(c, http.StatusOK, comp.Root(comp.ArticleNotFound(), fmt.Sprintf("%s - Not found", article)))
	}

	return Render(c, http.StatusOK, comp.Root(comp.Article(articleMetadata, articleHtmlContent), article))
}

func ServeArticles(c echo.Context) error {
	language := c.Param("language")
	fmt.Println("language:", language)

	articles, err := infra.Database.GetArticles(infra.ArticleSortingCriteria{}.Default())
	if err != nil {
		return Render(c, http.StatusOK, comp.Root(comp.ServerError(), "Server Error"))
	}

	return Render(c, http.StatusOK, comp.Root(comp.Articles(articles), language))
}

func mapFormDataToSortingCriteria(formData url.Values) infra.ArticleSortingCriteria {
	var sortingCriteria infra.ArticleSortingCriteria
	for key, value := range formData {
		if key == "sort" {
			sortingCriteria.SortBy = value[0]
		} else {
			sortingCriteria.FilterBy = append(sortingCriteria.FilterBy, key)
		}
	}

	return sortingCriteria
}

func ServeSortedArticles(c echo.Context) error {
	language := c.Param("language")
	fmt.Println("language:", language)

	formParams, err := c.FormParams()
	if err != nil {
		return Render(c, http.StatusOK, comp.Root(comp.ServerError(), "Server Error"))
	}

	sortingCriteria := mapFormDataToSortingCriteria(formParams)

	articles, err := infra.Database.GetArticles(sortingCriteria)
	if err != nil {
		return Render(c, http.StatusOK, comp.Root(comp.ServerError(), "Server Error"))
	}

	return Render(c, http.StatusOK, comp.ArticleGrid(articles))
}
