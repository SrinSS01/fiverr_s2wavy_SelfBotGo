package post

import (
	"encoding/json"
	"net/http"
	"s2wavy/selfbot/api/types"
	"strconv"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
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
	scheduledTime, _ := strconv.Atoi(schedule.InitiateTime)
	if time.Now().After(time.UnixMilli(int64(scheduledTime))) {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Please enter a valid time",
		})
	}

	if _, err := d.App.Dao().DB().Insert("message_schedulings", dbx.Params{
		"guild_id":        schedule.GuildID,
		"channel_id":      schedule.ChannelID,
		"selfbot_user_id": schedule.SelfbotUserID,
		"message_content": schedule.MessageContent,
		"initiate_time":   schedule.InitiateTime,
		"interval":        schedule.Interval,
		"expired":         schedule.Expired,
	}).Execute(); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
			"error":   err,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Scheduled",
	})
}

var ScheduleRequestFunction = ScheduleRequest{
	Path: "/schedule",
}
