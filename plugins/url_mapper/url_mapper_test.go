package url_mapper

import (
	"github.com/dancannon/gofetch/document"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSetup(t *testing.T) {
	var e *UrlMapperExtractor

	Convey("Subject: Setup UrlMapper extractor", t, func() {
		Convey("When the extractor is setup", func() {
			Convey("With a pattern and replacement", func() {
				e = &UrlMapperExtractor{}
				err := e.Setup(map[string]interface{}{
					"pattern":     "pattern",
					"replacement": "replacement",
				})

				Convey("No error is returned", func() {
					So(err, ShouldBeNil)
				})
				Convey("The value of the selector field is changed", func() {
					So(e.pattern, ShouldEqual, "pattern")
				})
				Convey("The value of the attribute field is changed", func() {
					So(e.replacement, ShouldEqual, "replacement")
				})
			})
			Convey("With no parameters", func() {
				e = &UrlMapperExtractor{}
				err := e.Setup(map[string]interface{}{})

				Convey("An error is returned", func() {
					So(err, ShouldNotBeNil)
				})
			})
		})
	})
}

func TestExtract(t *testing.T) {
	var e *UrlMapperExtractor

	Convey("Subject: Extract values from page", t, func() {
		Convey("When the extractor is setup", func() {
			Convey("With a pattern that matches the document URL", func() {
				e = &UrlMapperExtractor{}
				err := e.Setup(map[string]interface{}{
					"pattern":     "^http://example.com/(.*)$",
					"replacement": "http://example.com/api/$1.json",
				})
				if err != nil {
					t.Fatal(err.Error())
				}

				Convey("And the extractor is run", func() {
					doc, err := document.NewDocument("http://example.com/resource", strings.NewReader(""))
					if err != nil {
						t.Fatal(err.Error())
					}

					res, err := e.Extract(*doc)
					Convey("No error was returned", func() {
						So(err, ShouldBeNil)
					})
					Convey("The result should be valid", func() {
						So(res, ShouldEqual, "http://example.com/api/resource.json")
					})
				})
			})
			Convey("With a pattern that does not match the document URL", func() {
				e = &UrlMapperExtractor{}
				err := e.Setup(map[string]interface{}{
					"pattern":     "^http://example.com/([a-z]+)$",
					"replacement": "http://example.com/api/$1.json",
				})
				if err != nil {
					t.Fatal(err.Error())
				}

				Convey("And the extractor is run", func() {
					doc, err := document.NewDocument("http://example.com/resource1", strings.NewReader(""))
					if err != nil {
						t.Fatal(err.Error())
					}

					res, err := e.Extract(*doc)
					Convey("No error was returned", func() {
						So(err, ShouldBeNil)
					})
					Convey("The result should be the same as the input", func() {
						So(res, ShouldEqual, "http://example.com/resource1")
					})
				})
			})
			Convey("With an invalid pattern", func() {
				e = &UrlMapperExtractor{}
				err := e.Setup(map[string]interface{}{
					"pattern":     "^ttp://example.com/((($",
					"replacement": "http://example.com/api/$1.json",
				})
				if err != nil {
					t.Fatal(err.Error())
				}

				Convey("And the extractor is run", func() {
					doc, err := document.NewDocument("http://example.com/api/resource.json", strings.NewReader(""))
					if err != nil {
						t.Fatal(err.Error())
					}

					_, err = e.Extract(*doc)
					Convey("An error was returned", func() {
						So(err, ShouldNotBeNil)
					})
				})
			})
		})
	})
}
