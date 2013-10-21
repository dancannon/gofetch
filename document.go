package gofetch

type Document struct {
	Title string
	Mime  string
	Meta  map[string]interface{}
	Body  string
}

func NewDocument() *Document {
	doc := &Document{}
	doc.Meta = map[string]interface{}{}

	return doc
}
