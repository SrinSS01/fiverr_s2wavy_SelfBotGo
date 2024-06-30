package post

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"net/http"
	"s2wavy/selfbot/bots"
)

type TagRequest struct {
	Path string
	App  *pocketbase.PocketBase
}

type Tag struct {
	Name  string `json:"name"`
	Reply string `json:"reply"`
}

func (tr *TagRequest) Execute(c echo.Context) error {
	userId := c.PathParam("user_id")
	var tag Tag
	if err := json.NewDecoder(c.Request().Body).Decode(&tag); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
			"error":   err,
		})
	}
	_, err := tr.App.Dao().DB().Insert("tags", dbx.Params{
		"user_id": userId,
		"name":    tag.Name,
		"reply":   tag.Reply,
	}).Execute()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": err.Error(),
			"error":   err,
		})
	}
	bots.CommandHandlers[tag.Name] = func(session *discordgo.Session, message *discordgo.Message, _ string) {
		if message.Author.ID != userId {
			return
		}
		_, err := session.ChannelMessageSend(message.ChannelID, tag.Reply)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Successfully created tag",
		"tag":     tag.Name,
	})
}

var TagRequestFunction = TagRequest{
	Path: "/tag/:user_id",
}
