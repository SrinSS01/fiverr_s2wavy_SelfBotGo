package post

import (
	"encoding/json"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"net/http"
	"s2wavy/selfbot/api/types"
)

type ScheduleRequest struct {
	Path string
	App  *pocketbase.PocketBase
}

func (d *ScheduleRequest) Execute(c echo.Context) error {
	var schedule types.MessageScheduling
	if err := json.NewDecoder(c.Request().Body).Decode(&schedule); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
			"error":   err,
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": schedule,
	})
}

var ScheduleRequestFunction = ScheduleRequest{
	Path: "/schedule",
}
