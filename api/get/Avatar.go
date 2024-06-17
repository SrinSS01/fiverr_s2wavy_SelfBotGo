package get

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v5"
)

type AvatarRequest struct {
	Path string
}

func (r *AvatarRequest) Execute(c echo.Context) error {
	userId := c.PathParam("user_id")
	avatar := c.PathParam("avatar")

	var url string
	if avatar != "null" {
		url = fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png?size=256", userId, avatar)
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

var AvatarRequestFunction = AvatarRequest{
	Path: "/avatars/:user_id/:avatar",
}
