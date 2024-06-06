package get

import (
	// "github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
)

type SelfBotUsers struct {
	Path string
	App  *pocketbase.PocketBase
}

// func (d *SelfBotUsers) Execute(c echo.Context) error {
// 	d.App.Dao().FindRecordsByExpr("SelfBotUsers")
// }

var SelfBotUsersFunction = SelfBotUsers{
	Path: "/users",
}
