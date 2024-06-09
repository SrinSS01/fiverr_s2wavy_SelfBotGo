package types

type MessageScheduling struct {
	GuildID        string `json:"guild_id" db:"guild_id"`
	ChannelID      string `json:"channel_id" db:"channel_id"`
	SelfbotUserID  string `json:"selfbot_user_id" db:"selfbot_user_id"`
	MessageContent string `json:"message_content" db:"message_content"`
	InitiateTime   string `json:"initiate_time" db:"initiate_time"`
	Interval       int    `json:"interval" db:"interval"`
	Expired        bool   `json:"expired" db:"expired"`
}
