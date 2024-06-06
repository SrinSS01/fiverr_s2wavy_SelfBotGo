package main

import (
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
)

var (
	discord          *discordgo.Session
	GetApiFunctions  = map[string]func(c echo.Context) error{}
	PostApiFunctions = map[string]func(c echo.Context) error{}
)

// func onReady(s *discordgo.Session, r *discordgo.Ready) {
// 	log.Println(s.State.User.Username + " is ready")
// }
// func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
// 	if m.Author.ID == s.State.User.ID {
// 		return
// 	}
// 	fmt.Println(m.Content)
// }

func init() {
	os.Args = append(os.Args, "serve")
}

func main() {
	app := pocketbase.New()
	app.OnAfterBootstrap().Add(func(e *core.BootstrapEvent) error {
		selfbotUsers := &models.Collection{
			Name: "SelfBotUsers",
			Type: models.CollectionTypeBase,
			Schema: schema.NewSchema(
				&schema.SchemaField{
					Name:     "Token",
					Type:     schema.FieldTypeText,
					Required: true,
				},
				&schema.SchemaField{
					Name:     "Name",
					Type:     schema.FieldTypeText,
					Required: false,
				},
				&schema.SchemaField{
					Name:     "Avatar",
					Type:     schema.FieldTypeText,
					Required: false,
				},
			),
		}
		if err := e.App.Dao().SaveCollection(selfbotUsers); err != nil {
			return err
		}

		messageSchedulingConfigurations := &models.Collection{
			Name: "MessageSchedulings",
			Type: models.CollectionTypeBase,
			Schema: schema.NewSchema(
				&schema.SchemaField{
					Name:     "GuildID",
					Type:     schema.FieldTypeText,
					Required: true,
				},
				&schema.SchemaField{
					Name:     "SelfBotUserID",
					Type:     schema.FieldTypeText,
					Required: true,
				},
				&schema.SchemaField{
					Name:     "ChannelID",
					Type:     schema.FieldTypeText,
					Required: true,
				},
				&schema.SchemaField{
					Name:     "MessageContent",
					Type:     schema.FieldTypeText,
					Required: true,
				},
				&schema.SchemaField{
					Name:     "InitiateTime",
					Type:     schema.FieldTypeDate,
					Required: true,
				},
				&schema.SchemaField{
					Name:     "IntervalInSeconds",
					Type:     schema.FieldTypeNumber,
					Required: true,
				},
			),
		}

		if err := e.App.Dao().SaveCollection(messageSchedulingConfigurations); err != nil {
			return err
		}

		return nil
	})
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		return nil
	})
}
