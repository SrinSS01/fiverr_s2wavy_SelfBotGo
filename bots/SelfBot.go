package bots

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
	"regexp"
	"s2wavy/selfbot/commands"
	"time"

	"github.com/bwmarrin/discordgo"
)

type SelfBot struct {
	Session *discordgo.Session
	Running bool
	Timers  []*time.Ticker
}

var (
	CommandHandlers = map[string]func(*discordgo.Session, *discordgo.Message, string){
		commands.Ping.Command.Name: commands.Ping.Execute,
	}
	Bots         = map[string]*SelfBot{}
	commandRegex = regexp.MustCompile("-(?P<name>[\\w-]+)(?:\\s+(?P<args>.+))?")
)

func ReadyEvent(s *discordgo.Session, r *discordgo.Ready) {
	fmt.Println(r.User.Username, "is ready")
}

func MessageCreateEvent(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	request := resty.New().R()
	request.SetHeader("Authorization", s.Token)
	response, err := request.Execute("GET", "https://discord.com/api/v9/channels/"+m.ChannelID+"/messages?limit=1")
	if err != nil {
		fmt.Println(err)
		return
	}
	if response.StatusCode() != http.StatusOK {
		fmt.Println(string(response.Body()))
		return
	}
	var message []*discordgo.Message
	err = json.Unmarshal(response.Body(), &message)
	if err != nil {
		fmt.Println(err)
		return
	}
	matches := commandRegex.FindStringSubmatch(message[0].Content)
	if len(matches) == 0 {
		return
	}
	name := matches[commandRegex.SubexpIndex("name")]
	args := matches[commandRegex.SubexpIndex("args")]
	f := CommandHandlers[name]
	if f == nil {
		return
	}
	f(s, message[0], args)
}
