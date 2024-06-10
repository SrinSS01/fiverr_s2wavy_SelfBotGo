package get

import (
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
	guild_id := c.PathParam("guild_id")
	selfbot_user_id := c.PathParam("selfbot_user_id")
	channel_id := c.PathParam("channel_id")
	var schedules []*types.MessageScheduling
	if err := d.App.Dao().DB().
		Select("*").
		From("message_schedulings").
		Where(dbx.NewExp("guild_id = {:guild_id} and selfbot_user_id = {:selfbot_user_id} and channel_id = {:channel_id}", dbx.Params{
			"guild_id":        guild_id,
			"selfbot_user_id": selfbot_user_id,
			"channel_id":      channel_id,
		})).
		All(&schedules); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
			"error":   err,
		})
	}
	for _, schedule := range schedules {
		scheduledTime, _ := strconv.Atoi(schedule.InitiateTime)
		schedule.Expired = time.Now().After(time.UnixMilli(int64(scheduledTime)))
	}

	return c.JSON(http.StatusOK, schedules)
}

var ScheduleRequestFunction = ScheduleRequest{
	Path: "/schedules/:selfbot_user_id/:guild_id/:channel_id",
}
