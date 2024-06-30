package get

import (
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"net/http"
)

type TagRequest struct {
	Path string
	App  *pocketbase.PocketBase
}

type Tag struct {
	Name  string `json:"name" db:"name"`
	Reply string `json:"reply" db:"reply"`
}

func (tr *TagRequest) Execute(c echo.Context) error {
	userId := c.PathParam("user_id")
	var tags []Tag
	err := tr.App.Dao().DB().
		Select("name", "reply").
		From("tags").
		Where(dbx.NewExp("user_id = {:user_id}", dbx.Params{"user_id": userId})).
		All(&tags)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": err.Error(),
			"error":   err,
		})
	}
	return c.JSON(http.StatusOK, tags)
}

var TagRequestFunction = TagRequest{
	Path: "/tag/:user_id",
}
