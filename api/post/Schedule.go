package post

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"net/http"
	"s2wavy/selfbot/api/types"
	"s2wavy/selfbot/bots"
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
	bot := bots.Bots[schedule.SelfbotUserID]
	if bot == nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"message": "bot not found",
		})
	}
	if permissions, _ := bot.Session.UserChannelPermissions(schedule.SelfbotUserID, schedule.ChannelID); permissions&discordgo.PermissionViewChannel == 0 || permissions&discordgo.PermissionSendMessages == 0 {
		return c.JSON(http.StatusForbidden, map[string]interface{}{
			"message": "You do not have permission to send messages in this channel.",
		})
	}
	scheduledTime, _ := strconv.Atoi(schedule.InitiateTime)
	initiateUnixTime := time.UnixMilli(int64(scheduledTime))
	now := time.Now()
	if now.After(initiateUnixTime) {
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

	if bot.Running {
		time.AfterFunc(initiateUnixTime.Sub(now), func() {
			_, err := bot.Session.ChannelMessageSend(schedule.ChannelID, schedule.MessageContent)
			if err != nil {
				fmt.Println(err.Error())
			}
			ticker := time.NewTicker(time.Duration(schedule.Interval) * time.Second)
			bot.Timers = append(bot.Timers, ticker)
			for range ticker.C {
				_, err := bot.Session.ChannelMessageSend(schedule.ChannelID, schedule.MessageContent)
				if err != nil {
					fmt.Println(err.Error())
				}
			}
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Scheduled",
	})
}

var ScheduleRequestFunction = ScheduleRequest{
	Path: "/schedule",
}
