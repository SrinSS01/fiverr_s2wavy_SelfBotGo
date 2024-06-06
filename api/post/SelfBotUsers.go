package post

import (
	"encoding/json"
	"errors"
	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"net/http"
)

type SelfBotUsers struct {
	Path   string
	App    *pocketbase.PocketBase
	UserID string `db:"user_id" json:"user_id"`
	Token  string `db:"token" json:"token"`
	Name   string `db:"name" json:"name"`
	Avatar string `db:"avatar" json:"avatar"`
}
type User struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	//Discriminator        string        `json:"discriminator"`
	//PublicFlags          int           `json:"public_flags"`
	//Flags                int           `json:"flags"`
	Banner      interface{} `json:"banner"`
	AccentColor int         `json:"accent_color"`
	//GlobalName           string        `json:"global_name"`
	//AvatarDecorationData interface{}   `json:"avatar_decoration_data"`
	BannerColor string `json:"banner_color"`
	//Clan                 interface{}   `json:"clan"`
	//MfaEnabled           bool          `json:"mfa_enabled"`
	//Locale               string        `json:"locale"`
	//PremiumType          int           `json:"premium_type"`
	//Email                string        `json:"email"`
	//Verified             bool          `json:"verified"`
	//Phone                string        `json:"phone"`
	//NsfwAllowed          bool          `json:"nsfw_allowed"`
	//LinkedUsers          []interface{} `json:"linked_users"`
	//PurchasedFlags     int    `json:"purchased_flags"`
	Bio string `json:"bio"`
	//AuthenticatorTypes []int  `json:"authenticator_types"`
}

func (d *SelfBotUsers) Execute(c echo.Context) error {
	jsonMap := make(map[string]interface{})
	if err := json.NewDecoder(c.Request().Body).Decode(&jsonMap); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	token := jsonMap["token"].(string)
	if token == "" {
		return c.JSON(http.StatusBadRequest, errors.New("token is empty"))
	}
	request := resty.New().R()
	request.SetHeader("Authorization", token)
	response, err := request.Execute("GET", "https://discord.com/api/v9/users/@me")
	if err != nil || response.StatusCode() != http.StatusOK {
		return c.JSON(http.StatusInternalServerError, err)
	}
	user := User{}
	if err := json.Unmarshal(response.Body(), &user); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	if _, err := d.App.Dao().DB().NewQuery("insert into self_bot_users values ({:user_id}, {:token}, {:name}, {:avatar})").Bind(dbx.Params{
		"user_id": user.Id,
		"token":   token,
		"name":    user.Username,
		"avatar":  user.Avatar,
	}).Execute(); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "successfully added user!",
		"user":    user,
	})
}

var SelfBotUsersFunction = SelfBotUsers{
	Path: "/self_bot_users",
}
