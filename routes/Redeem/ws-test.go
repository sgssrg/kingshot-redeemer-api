package redeem

import (
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v5"
)

var (
	upgrader = websocket.Upgrader{}
)

// WSTest is a simple WebSocket echo test endpoint.
// @Summary WebSocket test
// @Description Upgrades the connection to a WebSocket and greets the client.
// @Tags Redeem
// @Success 101
// @Router /redeem/ws [get]
func WSTest(c *echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	for {
		// Write
		err := ws.WriteMessage(websocket.TextMessage, []byte("Hello, Client!"))
		if err != nil {
			c.Logger().Error("failed to write WS message", "error", err)
		}

		// Read
		_, msg, err := ws.ReadMessage()
		if err != nil {
			c.Logger().Error("failed to read WS message", "error", err)
		}
		fmt.Printf("%s\n", msg)
	}
}
