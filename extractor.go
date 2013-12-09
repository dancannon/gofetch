package gofetch

type Extractor interface {
	Id() string
	Setup(map[string]interface{}) error
	Extract(*Document, *Result) (interface{}, error)
}
