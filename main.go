package main

import (
	"log"
	"s2wavy/selfbot/api/delete"
	"s2wavy/selfbot/api/get"
	"s2wavy/selfbot/api/post"
	"s2wavy/selfbot/api/types"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

var (
	GetApiFunctions = map[string]func(c echo.Context) error{
		get.SelfBotUsersFunction.Path:    get.SelfBotUsersFunction.Execute,
		get.ServersRequestFunction.Path:  get.ServersRequestFunction.Execute,
		get.ChannelsRequestFunction.Path: get.ChannelsRequestFunction.Execute,
		get.ScheduleRequestFunction.Path: get.ScheduleRequestFunction.Execute,
	}
	PostApiFunctions = map[string]func(c echo.Context) error{
		post.SelfBotUsersFunction.Path:    post.SelfBotUsersFunction.Execute,
		post.StartBotRequestFunction.Path: post.StartBotRequestFunction.Execute,
		post.StopBotRequestFunction.Path:  post.StopBotRequestFunction.Execute,
		post.ScheduleRequestFunction.Path: post.ScheduleRequestFunction.Execute,
	}
	DeleteApiFunctions = map[string]func(c echo.Context) error{
		delete.SelfBotUsersFunction.Path: delete.SelfBotUsersFunction.Execute,
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

	get.SelfBotUsersFunction.App = app
	get.ServersRequestFunction.App = app
	get.ChannelsRequestFunction.App = app
	get.ScheduleRequestFunction.App = app

	delete.SelfBotUsersFunction.App = app

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
		selfBotUsersRowsAffected, err := selfBotUsersResult.RowsAffected()
		if err != nil {
			return err
		}
		messageSchedulingRowsAffected, err := messageSchedulingResult.RowsAffected()
		if err != nil {
			return err
		}
		if selfBotUsersRowsAffected == 0 && messageSchedulingRowsAffected == 0 {
			var users []struct {
				UserId string `db:"user_id"`
				Token  string `db:"token"`
			}
			if err := db.Select("user_id", "token").From("self_bot_users").All(&users); err != nil {
				return err
			}
			for _, user := range users {
				if err := post.SetBotSessionCache(user.Token, user.UserId); err != nil {
					return err
				}
			}
			var messageSchedulings []types.MessageScheduling
			err := db.Select("*").From("message_schedulings").All(&messageSchedulings)
			if err != nil {
				return err
			}
			//for _, scheduling := range messageSchedulings {
			//	atoi, _ := strconv.Atoi(scheduling.InitiateTime)
			//	scheduleTime := time.Unix(int64(atoi), 0)
			//}
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
