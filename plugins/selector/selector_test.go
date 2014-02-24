package selector

import (
	"github.com/dancannon/gofetch/document"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSetup(t *testing.T) {
	var e *SelectorExtractor

	Convey("Subject: Setup Selector extractor", t, func() {
		Convey("When the extractor is setup", func() {
			Convey("With a selector", func() {
				e = &SelectorExtractor{}
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
			Convey("With a selector and attribute", func() {
				e = &SelectorExtractor{}
				err := e.Setup(map[string]interface{}{
					"selector":  ".class",
					"attribute": "attr",
				})

				Convey("No error is returned", func() {
					So(err, ShouldBeNil)
				})
				Convey("The value of the selector field is changed", func() {
					So(e.selector, ShouldEqual, ".class")
				})
				Convey("The value of the attribute field is changed", func() {
					So(e.attribute, ShouldEqual, "attr")
				})
			})
			Convey("With no parameters", func() {
				e = &SelectorExtractor{}
				err := e.Setup(map[string]interface{}{})

				Convey("An error is returned", func() {
					So(err, ShouldNotBeNil)
				})
			})
		})
	})
}

func TestExtract(t *testing.T) {
	var e *SelectorExtractor

	Convey("Subject: Extract values from page", t, func() {
		Convey("When the extractor is setup", func() {
			Convey("With a selector that exists in the document", func() {
				Convey("And no attribute", func() {
					e = &SelectorExtractor{}
					err := e.Setup(map[string]interface{}{
						"selector": ".blog-title",
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
							So(res, ShouldEqual, "The Bootstrap Blog")
						})
					})
				})
				Convey("And an attribute that exists in the resulting element", func() {
					e = &SelectorExtractor{}
					err := e.Setup(map[string]interface{}{
						"selector":  ".blog-title",
						"attribute": "class",
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
							So(res, ShouldEqual, "blog-title")
						})
					})
				})
			})
			Convey("With a selector that does not exists in the document", func() {
				Convey("And no attribute", func() {
					e = &SelectorExtractor{}
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
	})
}
