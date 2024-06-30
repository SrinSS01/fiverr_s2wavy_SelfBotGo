package get

import (
	"encoding/json"
	"net/http"
	"s2wavy/selfbot/api/types"
	"s2wavy/selfbot/bots"

	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
)

type ServersRequest struct {
	Path string
	App  *pocketbase.PocketBase
}

func (d *ServersRequest) Execute(c echo.Context) error {
	userId := c.PathParam("user_id")
	bot := bots.Bots[userId]
	if bot == nil {
		return c.JSON(http.StatusNotFound, "Bot not found")
	}
	token := bot.Session.Token
	//var token struct {
	//	Token string `json:"token"`
	//}
	//// fetch user with the user_id
	//err := d.App.Dao().DB().
	//	Select("token").
	//	From("self_bot_users").
	//	Where(dbx.NewExp("user_id = {:user_id}", dbx.Params{
	//		"user_id": userId,
	//	})).
	//	One(&token)
	//if err != nil {
	//	return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	//}
	request := resty.New().R()
	request.SetHeader("Authorization", token)
	response, err := request.Execute("GET", "https://discord.com/api/v9/users/@me/guilds")
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

	var guilds []*types.Guild
	if err := json.Unmarshal(response.Body(), &guilds); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": err.Error(),
			"error":   err,
		})
	}
	var configuredGuildIds []struct {
		GuildId string `db:"guild_id"`
	}
	err = d.App.Dao().DB().
		Select("guild_id").
		From("message_schedulings").
		Where(dbx.NewExp("selfbot_user_id = {:selfbot_user_id}", dbx.Params{
			"selfbot_user_id": userId,
		})).
		GroupBy("guild_id").
		All(&configuredGuildIds)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
			"error":   err,
		})
	}
	for _, guild := range guilds {
		for _, cGuild := range configuredGuildIds {
			if guild.Id == cGuild.GuildId {
				guild.Configured = true
				break
			}
		}
	}
	return c.JSON(http.StatusOK, guilds)
}

var ServersRequestFunction = ServersRequest{
	Path: "/servers/:user_id",
}
