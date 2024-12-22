package handlers

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"

	"github.com/a-h/templ"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

func Init() {
	fmt.Println("Startup sequence starting...")

	err := godotenv.Load(".env")
	if err != nil {
		log.Println("WARNING : Failed to load .env file")
	}

	fmt.Println("Startup sequence done.")
}

func Render(ctx echo.Context, statusCode int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := t.Render(ctx.Request().Context(), buf); err != nil {
		return err
	}

	return ctx.HTML(statusCode, buf.String())
}
