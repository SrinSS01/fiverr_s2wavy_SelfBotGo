package get

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"net/http"
	"s2wavy/selfbot/api/types"
)

type ChannelsRequest struct {
	Path string
	App  *pocketbase.PocketBase
}

func (d *ChannelsRequest) Execute(c echo.Context) error {
	userId := c.PathParam("user_id")
	guildId := c.PathParam("guild_id")
	var token struct {
		Token string `json:"token"`
	}
	// fetch user with the user_id
	err := d.App.Dao().DB().NewQuery("select token from self_bot_users where user_id = {:user_id}").Bind(dbx.Params{
		"user_id": userId,
	}).One(&token)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	request := resty.New().R()
	request.SetHeader("Authorization", token.Token)
	response, err := request.Execute("GET", "https://discord.com/api/v9/guilds/"+guildId+"/channels")
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
	var channels []*types.Channel
	if err := json.Unmarshal(response.Body(), &channels); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": err.Error(),
			"error":   err,
		})
	}
	var configuredChannelIds []struct {
		ChannelId string `db:"channel_id"`
	}
	err = d.App.Dao().
	DB().
	Select("channel_id").
	From("message_schedulings").
	Where(dbx.NewExp("selfbot_user_id = {:selfbot_user_id}", dbx.Params{
		"selfbot_user_id": userId,
	})).
	GroupBy("channel_id").All(&configuredChannelIds)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
			"error":   err,
		})
	}
	for _, channel := range channels {
		for _, cChannel := range configuredChannelIds {
			if channel.Id == cChannel.ChannelId {
				channel.Configured = true
				break
			}
		}
	}
	return c.JSON(http.StatusOK, channels)
}

var ChannelsRequestFunction = ChannelsRequest{
	Path: "/channels/:user_id/:guild_id",
}
