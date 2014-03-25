package selector_text

import (
	"os"
	"testing"

	"github.com/dancannon/gofetch/document"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSetup(t *testing.T) {
	var e *SelectorTextExtractor

	Convey("Subject: Setup Selector extractor", t, func() {
		Convey("When the extractor is setup", func() {
			Convey("With a selector", func() {
				e = &SelectorTextExtractor{}
				err := e.Setup(map[string]interface{}{
					"selector": ".class",
				})

				Convey("No error is returned", func() {
					So(err, ShouldBeNil)
				})
				Convey("The value of the selector field is changed", func() {
					So(e.selector, ShouldEqual, ".class")
				})
			})
			Convey("With no parameters", func() {
				e = &SelectorTextExtractor{}
				err := e.Setup(map[string]interface{}{})

				Convey("An error is returned", func() {
					So(err, ShouldNotBeNil)
				})
			})
		})
	})
}

func TestExtract(t *testing.T) {
	var e *SelectorTextExtractor

	Convey("Subject: Extract values from page", t, func() {
		Convey("When the extractor is setup", func() {
			Convey("With a selector that exists in the document", func() {
				e = &SelectorTextExtractor{}
				err := e.Setup(map[string]interface{}{
					"selector": ".blog-description",
				})
				if err != nil {
					t.Fatal(err.Error())
				}

				Convey("And the extractor is run", func() {
					f, err := os.Open("../../test/blog.html")
					if err != nil {
						t.Fatal(err.Error())
					}
					doc, err := document.NewDocument("", f)
					if err != nil {
						t.Fatal(err.Error())
					}

					res, err := e.Extract(*doc)
					Convey("No error was returned", func() {
						So(err, ShouldBeNil)
					})
					Convey("The result should be valid", func() {
						So(res, ShouldEqual, "<p>The official example template of creating a blog with Bootstrap. </p>")
					})
				})
			})
			Convey("With a selector that does not exists in the document", func() {
				e = &SelectorTextExtractor{}
				err := e.Setup(map[string]interface{}{
					"selector": ".not-here",
				})
				if err != nil {
					t.Fatal(err.Error())
				}

				Convey("And the extractor is run", func() {
					f, err := os.Open("../../test/blog.html")
					if err != nil {
						t.Fatal(err.Error())
					}
					doc, err := document.NewDocument("", f)
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
