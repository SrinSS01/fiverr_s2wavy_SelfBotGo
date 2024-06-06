package get

import (
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"net/http"
)

type SelfBotUsers struct {
	Path   string
	App    *pocketbase.PocketBase
	UserID string `db:"user_id" json:"user_id"`
	Token  string `db:"token" json:"token"`
	Name   string `db:"name" json:"name"`
	Avatar string `db:"avatar" json:"avatar"`
}

func (d *SelfBotUsers) Execute(c echo.Context) error {
	var users []*SelfBotUsers
	if err := d.App.Dao().DB().NewQuery("select * from self_bot_users").All(&users); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, users)
}

var SelfBotUsersFunction = SelfBotUsers{
	Path: "/self_bot_users",
}
