package config

type Type struct {
	Id       string      `json:"id"`
	Validate bool        `json:"validate"`
	Values   interface{} `json:"values"`
}
