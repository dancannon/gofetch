package opengraph

import (
	"github.com/dancannon/gofetch/document"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSetup(t *testing.T) {
	Convey("Subject: Setup OpenGraph extractor", t, func() {
		e := &OpengraphExtractor{}

		Convey("When the extractor is setup", func() {
			err := e.Setup(map[string]interface{}{})

			Convey("No error is returned", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}

func TestSupports(t *testing.T) {
	Convey("Subject: Check if page supports OpenGraph", t, func() {
		e := &OpengraphExtractor{}

		Convey("Given a document that does not support OpenGraph", func() {
			f, err := os.Open("../../test/data/simple.html")
			if err != nil {
				panic(err.Error())
			}
			doc, err := document.NewDocument("", f)
			if err != nil {
				panic(err.Error())
			}

			Convey("The extractor will not support the document", func() {
				So(e.Supports(*doc), ShouldBeFalse)
			})
		})
		Convey("Given a document that does support OpenGraph", func() {
			f, err := os.Open("../../test/data/opengraph_text.html")
			if err != nil {
				panic(err.Error())
			}
			doc, err := document.NewDocument("", f)
			if err != nil {
				panic(err.Error())
			}

			Convey("The extractor will support the document", func() {
				So(e.Supports(*doc), ShouldBeTrue)
			})
		})
	})

}

func TestExtractValues(t *testing.T) {
	Convey("Subject: Extract values from page", t, func() {
		e := &OpengraphExtractor{}

		Convey("When the extractor is setup", func() {
			err := e.Setup(map[string]interface{}{})
			if err != nil {
				panic(err.Error())
			}

			Convey("And a page that supports opengraph is given", func() {
				Convey("And that page is an article", func() {
					f, err := os.Open("../../test/data/opengraph_text.html")
					if err != nil {
						panic(err.Error())
					}
					doc, err := document.NewDocument("", f)
					if err != nil {
						panic(err.Error())
					}

					res, typ, err := e.ExtractValues(*doc)
					Convey("No error was returned", func() {
						So(err, ShouldBeNil)
					})
					Convey("The type should be 'text'", func() {
						So(typ, ShouldEqual, "text")
					})
					Convey("The result should be valid", func() {
						So(res, ShouldResemble, map[string]interface{}{
							"title": "Title",
						})
					})
				})
				Convey("And that page is a photo", func() {
					f, err := os.Open("../../test/data/opengraph_photo.html")
					if err != nil {
						panic(err.Error())
					}
					doc, err := document.NewDocument("", f)
					if err != nil {
						panic(err.Error())
					}

					res, typ, err := e.ExtractValues(*doc)
					Convey("No error was returned", func() {
						So(err, ShouldBeNil)
					})
					Convey("The type should be 'image'", func() {
						So(typ, ShouldEqual, "image")
					})
					Convey("The result should be valid", func() {
						So(res, ShouldResemble, map[string]interface{}{
							"title":  "Title",
							"url":    "url",
							"width":  "640",
							"height": "478",
						})
					})
				})
				Convey("And that page is an unrecognised type", func() {
					f, err := os.Open("../../test/data/opengraph_other.html")
					if err != nil {
						panic(err.Error())
					}
					doc, err := document.NewDocument("", f)
					if err != nil {
						panic(err.Error())
					}

					res, typ, err := e.ExtractValues(*doc)
					Convey("No error was returned", func() {
						So(err, ShouldBeNil)
					})
					Convey("The type should be 'general'", func() {
						So(typ, ShouldEqual, "general")
					})
					Convey("The result should be valid", func() {
						So(res, ShouldResemble, map[string]interface{}{
							"title":   "Title",
							"content": "Description",
						})
					})
				})
			})
			Convey("And a page that does not support opengraph is given", func() {

			})
		})
	})
}

func TestCreateMapFromProps(t *testing.T) {
	Convey("Subject: Create map from properties", t, func() {
		var props map[string]interface{}

		Convey("Given a map of properties", func() {
			props = map[string]interface{}{
				"hello": "world",
				"foo":   "bar",
				"baz":   "baz",
			}

			Convey("When I create a map using a key that does not exist in the map", func() {
				m := createMapFromProps(props, map[string]string{
					"Test": "Test",
				})
				Convey("An empty map is returned", func() {
					So(m, ShouldResemble, map[string]interface{}{})
				})
			})
			Convey("When I create a map using a key that does exist in the map", func() {
				m := createMapFromProps(props, map[string]string{
					"hello": "hello",
				})
				Convey("An map with the key 'hello' is returned", func() {
					So(m, ShouldResemble, map[string]interface{}{
						"hello": "world",
					})
				})
			})
			Convey("When I create a map using multiple keys that do exist in the map", func() {
				m := createMapFromProps(props, map[string]string{
					"hello": "hello",
					"foo":   "foo",
				})
				Convey("An map with the keys 'hello' and 'foo' is returned", func() {
					So(m, ShouldResemble, map[string]interface{}{
						"hello": "world",
						"foo":   "bar",
					})
				})
			})
			Convey("When I create a map and rename a field", func() {
				m := createMapFromProps(props, map[string]string{
					"foo": "baz",
				})
				Convey("A map with the key 'foo' is returned", func() {
					So(m, ShouldResemble, map[string]interface{}{
						"foo": "baz",
					})
				})
			})
		})
	})
}
