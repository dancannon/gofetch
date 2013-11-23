package gofetch

import (
	"github.com/dancannon/gofetch/config"
	"github.com/dancannon/gofetch/plugin/selector"
	"github.com/dancannon/gofetch/plugin/selector_text"
	"github.com/dancannon/gofetch/plugin/text"
	"github.com/davecgh/go-spew/spew"
	"testing"
)

var c = config.LoadConfig("config.json")

func TestFetch(t *testing.T) {
	fetcher := NewFetcher(c)

	// Register all plugins
	RegisterExtractor(new(text.Extractor))
	RegisterExtractor(new(selector.Extractor))
	RegisterExtractor(new(selector_text.Extractor))

	// res, err := fetcher.Fetch("http://getbootstrap.com/examples/starter-template/")
	// res, err := fetcher.Fetch("http://getbootstrap.com/examples/jumbotron/")
	// res, err := fetcher.Fetch("http://getbootstrap.com/examples/carousel/")
	// res, err := fetcher.Fetch("http://www.theguardian.com/technology/2013/nov/01/caa-easa-electronic-devices-flight-take-off-landing")
	// res, err := fetcher.Fetch("http://www.birmingham.ac.uk/index.aspx")
	// res, err := fetcher.Fetch("http://www.birmingham.ac.uk/university/index.aspx")
	// res, err := fetcher.Fetch("https://www.google.co.uk/?gws_rd=cr&ei=IMtzUuLkI-Hb0QX-woD4CA#q=test")
	res, err := fetcher.Fetch("http://www.techradar.com/news/phone-and-communications/mobile-phones/blackberry-takeover-plan-abandoned-as-thorsten-heins-steps-down-1196359")
	// res, err := fetcher.Fetch("http://www.bbc.co.uk/news/business-24815793")
	// res, err := fetcher.Fetch("http://www.bbc.co.uk/news/technology-25042563")
	// res, err := fetcher.Fetch("http://imgur.com")
	// res, err := fetcher.Fetch("http://imgur.com/7T7MrBc")
	if err != nil {
		t.Errorf("Error was returned(%s)", err)
	}

	// spew.Dump(res.Content)
	spew.Print(res.Content)
}

func TestConfig(t *testing.T) {
	// spew.Dump(config)
}
