package post

import (
	"encoding/json"
	"net/http"
	"s2wavy/selfbot/api/types"
	"s2wavy/selfbot/bots"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/go-resty/resty/v2"
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
	token := jsonMap["token"].(string)
	if token == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": "Token cannot be empty"})
	}
	request := resty.New().R()
	request.SetHeader("Authorization", token)
	response, err := request.Execute("GET", "https://discord.com/api/v9/users/@me")
	if err != nil || response.StatusCode() != http.StatusOK {
		if err == nil {
			var body map[string]interface{}
			_ = json.Unmarshal(response.Body(), &body)
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": body,
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": err.Error(),
			"error":   err,
		})
	}
	user := types.User{}
	if err := json.Unmarshal(response.Body(), &user); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": err.Error(),
			"error":   err,
		})
	}
	if _, err := d.App.Dao().DB().Insert("self_bot_users", dbx.Params{
		"user_id": user.Id,
		"token":   token,
		"name":    user.Username,
		"avatar":  user.Avatar,
	}).Execute(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": err.Error(),
			"error":   err,
		})
	}

	if err := SetBotSessionCache(token, user.Id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "Failed to create discord session",
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "successfully added user!",
		"user":    user,
	})
}

func SetBotSessionCache(token, id string) error {
	discord, err := discordgo.New(token)
	if err != nil {
		return err
	}
	bots.Bots[id] = &bots.SelfBot{
		Session: discord,
		Running: false,
		Timers:  []*time.Ticker{},
	}
	discord.AddHandler(bots.ReadyEvent)
	discord.AddHandler(bots.MessageCreateEvent)
	return nil
}

var SelfBotUsersFunction = SelfBotUsersRequest{
	Path: "/self_bot_users",
}
