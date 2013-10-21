package gofetch

import (
	"github.com/davecgh/go-spew/spew"
	"testing"
)

func TestRequest(t *testing.T) {
	res, err := Fetch("http://getbootstrap.com/examples/starter-template/")
	// res, err := Fetch("http://www.youtube.com/watch?v=C0DPdy98e4c")
	doc := prepareDocument(res)
	// response, err := r.Send("http://hn.meteor.com")

	if err != nil {
		t.Errorf("Error was returned(%s)", err)
	}
	spew.Dump(doc)
}
