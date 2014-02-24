package plugins

import (
	"github.com/dancannon/gofetch/document"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type TestPlugin struct{}

func (t *TestPlugin) Setup(_ interface{}) error {
	return nil
}
func (t *TestPlugin) Extract(doc document.Document) (interface{}, error) {
	return nil, nil
}
func (t *TestPlugin) ExtractValues(doc document.Document) (interface{}, string, error) {
	return nil, "", nil
}

func TestSetup(t *testing.T) {
	Convey("Subject: Register plugin", t, func() {
		Convey("Given an instance of a plugin", func() {
			plugin := new(TestPlugin)

			Convey("Registering that plugin will add it to the plugins slice", func() {
				RegisterPlugin("base", plugin)
				So(GetExtractor("base"), ShouldNotBeNil)
				So(GetMultiExtractor("base"), ShouldNotBeNil)
			})
			Convey("Getting an extractor that doesnt exist will return nil", func() {
				So(GetExtractor("unknown"), ShouldBeNil)
				So(GetMultiExtractor("unknown"), ShouldBeNil)
			})
		})
	})
}
