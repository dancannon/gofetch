package gofetch

import (
	"github.com/dancannon/gofetch/request"
	"testing"
)

func TestRequest(t *testing.T) {
	var r request.Requester = &request.PhantomRequest{}

	content, err := r.Send("http://google.com")

	if err != nil {
		t.Error("Error was returned(%s)", err)
	}
	t.Log(content)
}
