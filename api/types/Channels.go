package types

type Channel struct {
	Id       string `json:"id"`
	Type     int    `json:"type"`
	GuildId  string `json:"guild_id"`
	Name     string `json:"name"`
	ParentId string `json:"parent_id"`
}
