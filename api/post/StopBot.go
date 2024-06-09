package post

import (
	"fmt"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"net/http"
	"s2wavy/selfbot/bots"
)

type StopBotRequest struct {
	Path string
	App  *pocketbase.PocketBase
}

func (d *StopBotRequest) Execute(c echo.Context) error {
	userId := c.PathParam("user_id")
	selfBot := bots.Bots[userId]
	if selfBot == nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"code":    http.StatusNotFound,
			"message": "Bot not found",
		})
	}
	if !selfBot.Running {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    http.StatusBadRequest,
			"message": "Bot is not running",
		})
	}
	if err := selfBot.Session.Close(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
			"error":   err,
		})
	}
	selfBot.Running = false
	fmt.Println("Bot stopped.", userId)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code":    http.StatusOK,
		"message": "Bot Stopped",
	})
}

var StopBotRequestFunction = StopBotRequest{
	Path: "/stop_bot/:user_id",
}
