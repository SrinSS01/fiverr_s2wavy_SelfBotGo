package post

import (
	"fmt"
	"net/http"
	"s2wavy/selfbot/api/types"
	"s2wavy/selfbot/bots"
	"strconv"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
)

type StartBotRequest struct {
	Path string
	App  *pocketbase.PocketBase
}

func (d *StartBotRequest) Execute(c echo.Context) error {
	userId := c.PathParam("user_id")
	selfBot := bots.Bots[userId]
	if selfBot == nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"code":    http.StatusNotFound,
			"message": "Bot not found",
		})
	}
	if err := selfBot.Session.Open(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
			"error":   err,
		})
	}
	selfBot.Running = true
	var messageSchedulings []types.MessageScheduling
	if err := d.App.Dao().DB().
		Select("*").
		From("message_schedulings").
		Where(dbx.NewExp("selfbot_user_id = {:selfbot_user_id}", dbx.Params{
			"selfbot_user_id": userId,
		})).All(&messageSchedulings); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
			"error":   err,
		})
	}

	go func() {
		now := time.Now()
		for _, schedule := range messageSchedulings {
			scheduleCopy := schedule
			scheduledTime, _ := strconv.Atoi(schedule.InitiateTime)
			initiateUnixTime := time.UnixMilli(int64(scheduledTime))
			var duration time.Duration

			schedule.Expired = now.After(initiateUnixTime)
			if schedule.Expired {
				duration = 0
			} else {
				duration = initiateUnixTime.Sub(now)
			}
			fmt.Println("content before", scheduleCopy)
			time.AfterFunc(duration, func() {
				fmt.Println("content after", scheduleCopy)
				ticker := time.NewTicker(time.Duration(scheduleCopy.Interval) * time.Second)
				selfBot.Timers = append(selfBot.Timers, ticker)
				for range ticker.C {
					_, err := selfBot.Session.ChannelMessageSend(scheduleCopy.ChannelID, scheduleCopy.MessageContent)
					if err != nil {
						fmt.Println("Unable to send message {", err.Error(), "}", scheduleCopy.MessageContent)
					}
				}
			})
		}
	}()
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code":    http.StatusOK,
		"message": "Bot started",
	})
}

var StartBotRequestFunction = StartBotRequest{
	Path: "/start_bot/:user_id",
}
