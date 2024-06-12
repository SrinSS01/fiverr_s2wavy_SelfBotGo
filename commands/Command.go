package commands

import "github.com/bwmarrin/discordgo"

type Command interface {
	Execute(session *discordgo.Session, messageCreate *discordgo.Message, args string)
}
