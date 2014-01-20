package config

type Type struct {
	Id         string      `json:"id"`
	AllowExtra bool        `json:"allow_extra"`
	Values     interface{} `json:"values"`
}
