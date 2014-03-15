package selector

import (
	"github.com/dancannon/gofetch/document"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSetup(t *testing.T) {
	var e *TemplateExtractor

	Convey("Subject: Setup Selector extractor", t, func() {
		Convey("When the extractor is setup", func() {
			Convey("With a template", func() {
				e = &TemplateExtractor{}
				err := e.Setup(map[string]interface{}{
					"template": "template",
				})

				Convey("No error is returned", func() {
					So(err, ShouldBeNil)
				})
				Convey("The value of the template field is changed", func() {
					So(e.template, ShouldEqual, "template")
				})
			})
			Convey("With no parameters", func() {
				e = &TemplateExtractor{}
				err := e.Setup(map[string]interface{}{})

				Convey("An error is returned", func() {
					So(err, ShouldNotBeNil)
				})
			})
		})
	})
}

func TestExtract(t *testing.T) {
	var e *TemplateExtractor

	Convey("Subject: Extract values from page", t, func() {
		Convey("When the extractor is setup", func() {
			Convey("With a valid template", func() {
				Convey("That has some variables", func() {
					e = &TemplateExtractor{}
					err := e.Setup(map[string]interface{}{
						"template": "<iframe width=\"100%\" height=\"240\" scrolling=\"no\" frameborder=\"no\" src=\"https://service/embed/?url={{.URL.RequestURI}}\"></iframe>",
					})
					if err != nil {
						t.Fatal(err.Error())
					}

					Convey("And the extractor is run", func() {
						f, err := os.Open("../../test/blog.html")
						if err != nil {
							t.Fatal(err.Error())
						}
						doc, err := document.NewDocument("http://blog/name?q=1", f)
						if err != nil {
							t.Fatal(err.Error())
						}

						res, err := e.Extract(*doc)
						Convey("No error was returned", func() {
							So(err, ShouldBeNil)
						})
						Convey("The result should be valid", func() {
							So(res, ShouldEqual, "<iframe width=\"100%\" height=\"240\" scrolling=\"no\" frameborder=\"no\" src=\"https://service/embed/?url=%2fname%3fq%3d1\"></iframe>")
						})
					})
				})
				Convey("That has no variables", func() {
					e = &TemplateExtractor{}
					err := e.Setup(map[string]interface{}{
						"template": "template",
					})
					if err != nil {
						t.Fatal(err.Error())
					}

					Convey("And the extractor is run", func() {
						f, err := os.Open("../../test/blog.html")
						if err != nil {
							t.Fatal(err.Error())
						}
						doc, err := document.NewDocument("http://blog/name", f)
						if err != nil {
							t.Fatal(err.Error())
						}

						res, err := e.Extract(*doc)
						Convey("No error was returned", func() {
							So(err, ShouldBeNil)
						})
						Convey("The result should be valid", func() {
							So(res, ShouldEqual, "template")
						})
					})
				})
			})
			Convey("With an invalid template", func() {
				Convey("And no attribute", func() {
					e = &TemplateExtractor{}
					err := e.Setup(map[string]interface{}{
						"template": "template {{",
					})
					if err != nil {
						t.Fatal(err.Error())
					}

					Convey("And the extractor is run", func() {
						f, err := os.Open("../../test/blog.html")
						if err != nil {
							t.Fatal(err.Error())
						}
						doc, err := document.NewDocument("http://blog/name", f)
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
	})
}
