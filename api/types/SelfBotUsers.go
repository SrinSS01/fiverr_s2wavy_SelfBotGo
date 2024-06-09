package types

type SelfBotUsers struct {
	UserID     string `db:"user_id" json:"user_id"`
	Name       string `db:"name" json:"name"`
	Avatar     string `db:"avatar" json:"avatar"`
	BotRunning bool   `db:"bot_running" json:"bot_running"`
}
