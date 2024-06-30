package delete

import (
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
	name := c.PathParam("name")
	result, err := tr.App.Dao().DB().Delete("tags", dbx.NewExp("user_id = {:user_id} and name = {:name}", dbx.Params{
		"user_id": userId,
		"name":    name,
	})).Execute()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": err.Error(),
			"error":   err,
		})
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": err.Error(),
		})
	}
	if rowsAffected == 0 {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"message": "tag not found",
		})
	}
	delete(bots.CommandHandlers, name)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Successfully created tag",
		"tag":     name,
	})
}

var TagRequestFunction = TagRequest{
	Path: "/tag/:user_id/:name",
}
