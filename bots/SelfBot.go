package bots

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

type SelfBot struct {
	Session *discordgo.Session
	Running bool
}

var Bots = map[string]*SelfBot{}

func ReadyEvent(s *discordgo.Session, r *discordgo.Ready) {
	fmt.Println(r.User.Username, "is ready")
}
