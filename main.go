package main

import (
	"log"
	"os"
	"s2wavy/selfbot/api/get"
	"s2wavy/selfbot/api/post"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

var (
	// discord          *discordgo.Session
	GetApiFunctions = map[string]func(c echo.Context) error{
		get.SelfBotUsersFunction.Path: get.SelfBotUsersFunction.Execute,
	}
	PostApiFunctions = map[string]func(c echo.Context) error{
		post.SelfBotUsersFunction.Path: post.SelfBotUsersFunction.Execute,
	}
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

	post.SelfBotUsersFunction.App = app
	get.SelfBotUsersFunction.App = app

	app.OnAfterBootstrap().Add(func(e *core.BootstrapEvent) error {
		db := e.App.Dao().DB()
		if _, err := db.NewQuery("create table if not exists self_bot_users(user_id text primary key, token text not null, name text, avatar text)").Execute(); err != nil {
			return err
		}
		if _, err := db.NewQuery("create table if not exists message_schedulings(guild_id text, channel_id text, selfbot_user_id text, message_content text, initiate_time text not null, interval int not null, primary key (guild_id, channel_id, selfbot_user_id))").Execute(); err != nil {
			return err
		}
		return nil
	})
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		for path, function := range GetApiFunctions {
			e.Router.GET(path, function)
		}
		for path, function := range PostApiFunctions {
			e.Router.POST(path, function)
		}
		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
