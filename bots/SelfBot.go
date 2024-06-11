package bots

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

type SelfBot struct {
	Session *discordgo.Session
	Running bool
	Timers  []*time.Ticker
}

var Bots = map[string]*SelfBot{}

func ReadyEvent(s *discordgo.Session, r *discordgo.Ready) {
	fmt.Println(r.User.Username, "is ready")
}

func MessageCreateEvent(s *discordgo.Session, m *discordgo.MessageCreate) {
	// fmt.Println(m.Author.GlobalName, m.Content)
}
