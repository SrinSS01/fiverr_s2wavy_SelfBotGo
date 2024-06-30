package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"log"
	"s2wavy/selfbot/api/delete"
	"s2wavy/selfbot/api/get"
	"s2wavy/selfbot/api/post"
	"s2wavy/selfbot/bots"
)

var (
	GetApiFunctions = map[string]func(c echo.Context) error{
		get.SelfBotUsersFunction.Path:     get.SelfBotUsersFunction.Execute,
		get.ServersRequestFunction.Path:   get.ServersRequestFunction.Execute,
		get.ChannelsRequestFunction.Path:  get.ChannelsRequestFunction.Execute,
		get.ScheduleRequestFunction.Path:  get.ScheduleRequestFunction.Execute,
		get.AvatarRequestFunction.Path:    get.AvatarRequestFunction.Execute,
		get.GuildIconRequestFunction.Path: get.GuildIconRequestFunction.Execute,
		get.TagRequestFunction.Path:       get.TagRequestFunction.Execute,
	}
	PostApiFunctions = map[string]func(c echo.Context) error{
		post.SelfBotUsersFunction.Path:    post.SelfBotUsersFunction.Execute,
		post.StartBotRequestFunction.Path: post.StartBotRequestFunction.Execute,
		post.StopBotRequestFunction.Path:  post.StopBotRequestFunction.Execute,
		post.ScheduleRequestFunction.Path: post.ScheduleRequestFunction.Execute,
		post.TagRequestFunction.Path:      post.TagRequestFunction.Execute,
	}
	DeleteApiFunctions = map[string]func(c echo.Context) error{
		delete.SelfBotUsersFunction.Path: delete.SelfBotUsersFunction.Execute,
		delete.TagRequestFunction.Path:   delete.TagRequestFunction.Execute,
	}
)

// func init() {
// 	os.Args = append(os.Args, "serve")
// }

func main() {
	app := pocketbase.New()

	post.SelfBotUsersFunction.App = app
	post.StartBotRequestFunction.App = app
	post.StopBotRequestFunction.App = app
	post.ScheduleRequestFunction.App = app
	post.TagRequestFunction.App = app

	get.SelfBotUsersFunction.App = app
	get.ServersRequestFunction.App = app
	get.ChannelsRequestFunction.App = app
	get.ScheduleRequestFunction.App = app
	get.TagRequestFunction.App = app

	delete.SelfBotUsersFunction.App = app
	delete.TagRequestFunction.App = app

	app.OnAfterBootstrap().Add(func(e *core.BootstrapEvent) error {
		db := e.App.Dao().DB()
		selfBotUsersResult, err := db.NewQuery("create table if not exists self_bot_users(user_id text primary key, token text not null, name text, avatar text)").Execute()
		if err != nil {
			return err
		}
		messageSchedulingResult, err := db.NewQuery("create table if not exists message_schedulings(guild_id text, channel_id text, selfbot_user_id text, message_content text, initiate_time text, interval int not null, expired boolean, primary key (initiate_time, guild_id, channel_id, selfbot_user_id))").Execute()
		if err != nil {
			return err
		}
		tagResult, err := db.NewQuery("create table if not exists tags(user_id text, name text, reply text, primary key (user_id, name))").Execute()
		if err != nil {
			return err
		}
		selfBotUsersRowsAffected, err := selfBotUsersResult.RowsAffected()
		if err != nil {
			return err
		}
		messageSchedulingRowsAffected, err := messageSchedulingResult.RowsAffected()
		if err != nil {
			return err
		}
		tagRowsAffected, err := tagResult.RowsAffected()
		if err != nil {
			return err
		}
		if selfBotUsersRowsAffected == 0 && messageSchedulingRowsAffected == 0 {
			var users []struct {
				UserId string `db:"user_id"`
				Token  string `db:"token"`
			}
			if err := db.Select("user_id", "token").
				From("self_bot_users").
				All(&users); err != nil {
				return err
			}
			for _, user := range users {
				if err := post.SetBotSessionCache(user.Token, user.UserId); err != nil {
					return err
				}
				if tagRowsAffected == 0 {
					var tags []post.Tag
					err := db.Select("name", "reply").
						From("tags").
						Where(dbx.NewExp("user_id = {:user_id}", dbx.Params{"user_id": user.UserId})).
						All(&tags)
					if err != nil {
						return err
					}
					for _, tag := range tags {
						bots.CommandHandlers[tag.Name] = func(session *discordgo.Session, message *discordgo.Message, _ string) {
							if message.Author.ID != user.UserId {
								return
							}
							_, err := session.ChannelMessageSend(message.ChannelID, tag.Reply)
							if err != nil {
								fmt.Println(err.Error())
								return
							}
						}
					}
				}
			}
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
		for path, function := range DeleteApiFunctions {
			e.Router.DELETE(path, function)
		}
		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
