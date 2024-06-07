package types

type SelfBotUsers struct {
	UserID string `db:"user_id" json:"user_id"`
	// Token  string `db:"token" json:"token"`
	Name   string `db:"name" json:"name"`
	Avatar string `db:"avatar" json:"avatar"`
}
