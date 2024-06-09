package get

import (
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"net/http"
	"s2wavy/selfbot/api/types"
	"s2wavy/selfbot/bots"
)

type SelfBotUsersRequest struct {
	Path string
	App  *pocketbase.PocketBase
}

func (d *SelfBotUsersRequest) Execute(c echo.Context) error {
	var users []*types.SelfBotUsers
	if err := d.App.Dao().DB().Select("*").From("self_bot_users").All(&users); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	for _, user := range users {
		user.BotRunning = bots.Bots[user.UserID] != nil && bots.Bots[user.UserID].Running
	}
	return c.JSON(http.StatusOK, users)
}

var SelfBotUsersFunction = SelfBotUsersRequest{
	Path: "/self_bot_users",
}
