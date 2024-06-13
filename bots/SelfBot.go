package bots

import (
	"fmt"
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
	// To add a new command just put the relevant name and execute function values as shown below
	CommandHandlers = map[string]func(*discordgo.Session, *discordgo.Message, string){
		commands.Ping.Command.Name: commands.Ping.Execute,
	}


	Bots         = map[string]*SelfBot{}
	commandRegex = regexp.MustCompile(`-(?P<name>[\w-]+)(?:\s+(?P<args>.+))?`)
)

func ReadyEvent(s *discordgo.Session, r *discordgo.Ready) {
	fmt.Println(r.User.Username, "is ready")
}

func MessageCreateEvent(session *discordgo.Session, event *discordgo.MessageCreate) {
	if event.Author.ID != session.State.User.ID {
		return
	}
	/* 
	The commented out section is used to fetch message content of other users, by default message contents are not available 
	for self-bots in the default message struct supplied by Message create event.
	
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
	} */
	message := event.Message
	matches := commandRegex.FindStringSubmatch(message.Content)
	if len(matches) == 0 {
		return
	}
	name := matches[commandRegex.SubexpIndex("name")]
	args := matches[commandRegex.SubexpIndex("args")]
	f := CommandHandlers[name]
	if f == nil {
		return
	}
	f(session, message, args)
}
