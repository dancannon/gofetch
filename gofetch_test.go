package gofetch

import (
	"github.com/dancannon/gofetch/plugin/article"
	"github.com/dancannon/gofetch/plugin/selector"
	"github.com/davecgh/go-spew/spew"
	"testing"
)

func TestRequest(t *testing.T) {
	// Register all plugins
	RegisterExtractor(new(article.Extractor))
	RegisterExtractor(new(selector.Extractor))

	// res, err := Fetch("http://getbootstrap.com/examples/starter-template/")
	// res, err := Fetch("http://getbootstrap.com/examples/jumbotron/")
	// res, err := Fetch("http://getbootstrap.com/examples/carousel/")
	// res, err := Fetch("http://www.theguardian.com/technology/2013/nov/01/caa-easa-electronic-devices-flight-take-off-landing")
	// res, err := Fetch("http://www.birmingham.ac.uk/index.aspx")
	// res, err := Fetch("http://www.birmingham.ac.uk/university/index.aspx")
	// res, err := Fetch("https://www.google.co.uk/?gws_rd=cr&ei=IMtzUuLkI-Hb0QX-woD4CA#q=test")
	// res, err := Fetch("http://www.techradar.com/news/phone-and-communications/mobile-phones/blackberry-takeover-plan-abandoned-as-thorsten-heins-steps-down-1196359")
	// res, err := Fetch("http://www.bbc.co.uk/news/business-24815793")
	// res, err := Fetch("http://imgur.com")
	res, err := Fetch("http://imgur.com/rXmjOMe")
	if err != nil {
		t.Errorf("Error was returned(%s)", err)
	}

	spew.Dump(res.Content)
}

func TestConfig(t *testing.T) {
	// spew.Dump(config)
}
