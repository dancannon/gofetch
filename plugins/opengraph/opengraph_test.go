package opengraph

import (
	"os"
	"testing"

	"github.com/dancannon/gofetch/document"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSetup(t *testing.T) {
	var e *OpengraphExtractor

	Convey("Subject: Setup OpenGraph extractor", t, func() {
		e = &OpengraphExtractor{}

		Convey("When the extractor is setup", func() {
			err := e.Setup(map[string]interface{}{})

			Convey("No error is returned", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}

func TestSupports(t *testing.T) {
	var e *OpengraphExtractor

	Convey("Subject: Check if page supports OpenGraph", t, func() {
		e = &OpengraphExtractor{}

		Convey("Given a document that does not support OpenGraph", func() {
			f, err := os.Open("../../test/simple.html")
			if err != nil {
				t.Fatal(err.Error())
			}
			doc, err := document.NewDocument("", f)
			if err != nil {
				t.Fatal(err.Error())
			}

			Convey("The extractor will not support the document", func() {
				So(e.Supports(*doc), ShouldBeFalse)
			})
		})
		Convey("Given a document that does support OpenGraph", func() {
			f, err := os.Open("../../test/opengraph_text.html")
			if err != nil {
				t.Fatal(err.Error())
			}
			doc, err := document.NewDocument("", f)
			if err != nil {
				t.Fatal(err.Error())
			}

			Convey("The extractor will support the document", func() {
				So(e.Supports(*doc), ShouldBeTrue)
			})
		})
	})
}

func TestExtractValues(t *testing.T) {
	var e *OpengraphExtractor

	Convey("Subject: Extract values from page", t, func() {
		Convey("When the extractor is setup", func() {
			e = &OpengraphExtractor{}
			err := e.Setup(map[string]interface{}{})
			if err != nil {
				t.Fatal(err.Error())
			}

			Convey("And a page that supports opengraph is given", func() {
				Convey("And that page is an article", func() {
					f, err := os.Open("../../test/opengraph_text.html")
					if err != nil {
						t.Fatal(err.Error())
					}
					doc, err := document.NewDocument("", f)
					if err != nil {
						t.Fatal(err.Error())
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
							"title": "Title",
						})
					})
				})
				Convey("And that page is a photo", func() {
					f, err := os.Open("../../test/opengraph_photo.html")
					if err != nil {
						t.Fatal(err.Error())
					}
					doc, err := document.NewDocument("", f)
					if err != nil {
						t.Fatal(err.Error())
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
					f, err := os.Open("../../test/opengraph_other.html")
					if err != nil {
						t.Fatal(err.Error())
					}
					doc, err := document.NewDocument("", f)
					if err != nil {
						t.Fatal(err.Error())
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
