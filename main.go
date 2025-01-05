package main

import (
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"potblog/handlers"
)

func main() {
	e := echo.New()
	if os.Getenv("ENV") == "dev" {
		e.Use(middleware.Logger())
	}
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
	e.Static("/assets", "assets")

	handlers.Init()

	// ---- Home Routes ---- //
	e.GET("/", handlers.ServeHomePage)
	e.GET("/:language", handlers.ServeHomePage)

	e.GET("/languages", handlers.ServeLanguageSelector)

	// ---- Article Routes ---- //
	e.GET("/:language/article/:article", handlers.ServeArticle)
	e.GET("/:language/articles", handlers.ServeArticles)
	e.POST("/:language/articles", handlers.ServeSortedArticles)

	// ---- Global Routes ---- //
	e.GET("/ping", handlers.ServePing)
	if os.Getenv("ENV") != "prod" {
		e.GET("/ws", handlers.InitHotReloadWebSocket)
	}

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", os.Getenv("PORT"))))
}
