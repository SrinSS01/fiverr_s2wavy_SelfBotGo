package get

import (
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"net/http"
	"s2wavy/selfbot/api/types"
	"s2wavy/selfbot/bots"
)

type ChannelsRequest struct {
	Path string
	App  *pocketbase.PocketBase
}

func (d *ChannelsRequest) Execute(c echo.Context) error {
	userId := c.PathParam("user_id")
	guildId := c.PathParam("guild_id")
	bot := bots.Bots[userId]
	if bot == nil {
		return c.JSON(http.StatusNotFound, "Bot not found")
	}
	token := bot.Session.Token
	request := resty.New().R()
	request.SetHeader("Authorization", token)
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
	var filteredChannels []*types.Channel
	for _, channel := range channels {
		if channel.Type != discordgo.ChannelTypeGuildText {
			continue
		}
		//permissions, _ := bot.Session.UserChannelPermissions(userId, channel.Id)
		//if permissions&discordgo.PermissionViewChannel == 0 || permissions&discordgo.PermissionSendMessages == 0 {
		//	continue
		//}
		filteredChannels = append(filteredChannels, channel)
		for _, cChannel := range configuredChannelIds {
			if channel.Id == cChannel.ChannelId {
				channel.Configured = true
				break
			}
		}
	}
	return c.JSON(http.StatusOK, filteredChannels)
}

var ChannelsRequestFunction = ChannelsRequest{
	Path: "/channels/:user_id/:guild_id",
}
