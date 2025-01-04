package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	comp "potblog/components"
)

func ServeHomePage(c echo.Context) error {
	language := c.Param("language")
	fmt.Println("language:", language)

	return Render(c, http.StatusOK, comp.Root(comp.Home(), "Home"))
}

func ServeLanguageSelector(c echo.Context) error {
	return Render(c, http.StatusOK, comp.LanguageSelectorMenu())
}
