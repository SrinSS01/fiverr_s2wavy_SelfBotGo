package types

import "github.com/bwmarrin/discordgo"

type Channel struct {
	Id         string                `json:"id"`
	Type       discordgo.ChannelType `json:"type"`
	GuildId    string                `json:"guild_id"`
	Name       string                `json:"name"`
	ParentId   string                `json:"parent_id"`
	Configured bool                  `json:"configured"`
}
