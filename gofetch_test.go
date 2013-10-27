package gofetch

import (
	"code.google.com/p/go.net/html"
	"os"
	"testing"
)

func TestRequest(t *testing.T) {
	res, err := Fetch("http://getbootstrap.com/examples/starter-template/")
	// res, err := Fetch("http://www.birmingham.ac.uk/index.aspx")
	// res, err := Fetch("http://www.youtube.com/watch?v=C0DPdy98e4c")
	doc := prepareDocument(res)
	// response, err := r.Send("http://hn.meteor.com")

	if err != nil {
		t.Errorf("Error was returned(%s)", err)
	}
	html.Render(os.Stderr, doc.Body)
}
