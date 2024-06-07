package delete

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
)

type SelfBotUsersRequest struct {
	Path string
	App  *pocketbase.PocketBase
}

func (d *SelfBotUsersRequest) Execute(c echo.Context) error {
	jsonMap := make(map[string]interface{})
	if err := json.NewDecoder(c.Request().Body).Decode(&jsonMap); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
			"error":   err,
		})
	}
	userId := jsonMap["user_id"].(string)
	result, err := d.App.Dao().DB().NewQuery("delete from self_bot_users where user_id = {:user_id}").Bind(dbx.Params{
		"user_id": userId,
	}).Execute()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
			"error":   err,
		})
	}
	if rows, err := result.RowsAffected(); rows == 0 || err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message":       err.Error(),
			"error":         err,
			"rows_affected": rows,
		})
	}
	return c.JSON(http.StatusOK, map[string]string{
		"message": "successfully deleted user",
		"user_id": userId,
	})
}

var SelfBotUsersFunction = SelfBotUsersRequest{
	Path: "/self_bot_users",
}
