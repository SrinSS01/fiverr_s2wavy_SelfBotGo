package get

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v5"
)

type GuildIconRequest struct {
	Path string
}

func (gir *GuildIconRequest) Execute(c echo.Context) error {
	guildId := c.PathParam("guild_id")
	icon := c.PathParam("icon")

	var url string
	if icon != "null" {
		url = fmt.Sprintf("https://cdn.discordapp.com/icons/%s/%s.png?size=256", guildId, icon)
	} else {
		url = "https://cdn.discordapp.com/embed/avatars/0.png"
	}
	request := resty.New().R()
	response, err := request.Get(url)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   err,
			"message": err.Error(),
		})
	}
	if response.StatusCode() != http.StatusOK {
		return c.JSON(response.StatusCode(), response.Body())
	}
	return c.Blob(http.StatusOK, "image/png", response.Body())
}

var GuildIconRequestFunction = GuildIconRequest{
	Path: "/guild_icons/:guild_id/:icon",
}
