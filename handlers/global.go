package handlers

import (
	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
)

func ServePing(c echo.Context) error {
	return c.String(200, "pong")
}

func InitHotReloadWebSocket(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		for {
			msg := ""
			err := websocket.Message.Receive(ws, &msg)
			if err != nil {
				return
			}
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}
