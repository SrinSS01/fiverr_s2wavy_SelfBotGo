package post

import (
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"net/http"
	"s2wavy/selfbot/bots"
)

type StartBotRequest struct {
	Path string
	App  *pocketbase.PocketBase
}

func (d *StartBotRequest) Execute(c echo.Context) error {
	userId := c.PathParam("user_id")
	selfBot := bots.Bots[userId]
	if selfBot == nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"code":    http.StatusNotFound,
			"message": "Bot not found",
		})
	}
	if err := selfBot.Session.Open(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
			"error":   err,
		})
	}
	selfBot.Running = true
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code":    http.StatusOK,
		"message": "Bot started",
	})
}

var StartBotRequestFunction = StartBotRequest{
	Path: "/start_bot/:user_id",
}
