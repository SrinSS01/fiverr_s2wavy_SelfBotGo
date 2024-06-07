package get

import (
	"net/http"
	"s2wavy/selfbot/api/types"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
)

type SelfBotUsersRequest struct {
	Path string
	App  *pocketbase.PocketBase
}

func (d *SelfBotUsersRequest) Execute(c echo.Context) error {
	var users []*types.SelfBotUsers
	if err := d.App.Dao().DB().NewQuery("select * from self_bot_users").All(&users); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, users)
}

var SelfBotUsersFunction = SelfBotUsersRequest{
	Path: "/self_bot_users",
}
