package types

type Guild struct {
	Id         string      `json:"id"`
	Name       string      `json:"name"`
	Icon       interface{} `json:"icon"`
	Configured bool        `json:"configured"`
}
