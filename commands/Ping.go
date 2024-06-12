package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

type PingCommand struct {
	Command *discordgo.ApplicationCommand
}

var Ping = &PingCommand{
	Command: &discordgo.ApplicationCommand{
		Name:        "ping",
		Description: "Pings the server",
	},
}

func (command *PingCommand) Execute(s *discordgo.Session, m *discordgo.Message, args string) {
	_, err := s.ChannelMessageSend(m.ChannelID, "Pong!")
	if err != nil {
		fmt.Println("Error sending pong:", err.Error())
		return
	}
}
