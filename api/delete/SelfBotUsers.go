package delete

import (
	"fmt"
	"net/http"
	"s2wavy/selfbot/bots"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
)

type SelfBotUsersRequest struct {
	Path string
	App  *pocketbase.PocketBase
}

func (d *SelfBotUsersRequest) Execute(c echo.Context) error {
	userId := c.PathParam("user_id")
	result, err := d.App.Dao().DB().Delete("self_bot_users", dbx.NewExp("user_id = {:user_id}", dbx.Params{
		"user_id": userId,
	})).Execute()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
			"error":   err,
		})
	}
	if rows, err := result.RowsAffected(); rows == 0 || err != nil {
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message":       err.Error(),
				"error":         err,
				"rows_affected": rows,
			})
		}
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"rows_affected": rows,
		})
	}
	selfBot := bots.Bots[userId]
	if selfBot == nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Successfully deleted user but unable to stop the bot",
			"user_id": userId,
		})
	}
	if !selfBot.Running {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Successfully deleted user but no running instance of bot found!",
			"user_id": userId,
		})
	}
	if user, err := selfBot.Session.User("@me"); err == nil {
		fmt.Println("deleting user", user.Username)
		if err := selfBot.Session.Close(); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": err.Error(),
				"error":   err,
			})
		}
		for _, timer := range selfBot.Timers {
			timer.Stop()
		}
		selfBot.Timers = []*time.Ticker{}
		delete(bots.Bots, userId)
		fmt.Println(user.Username, "Stopped and deleted.")
	} else {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Successfully deleted user but unable to stop the bot",
			"user_id": userId,
			"error":   err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]string{
		"message": "successfully deleted user",
		"user_id": userId,
	})
}

var SelfBotUsersFunction = SelfBotUsersRequest{
	Path: "/self_bot_users/:user_id",
}
