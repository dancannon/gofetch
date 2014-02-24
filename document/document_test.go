package document

import (
	"code.google.com/p/go.net/html"
	"os"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLoadConfig(t *testing.T) {
	Convey("Subject: New document", t, func() {
		Convey("Given a valid HTML document", func() {
			Convey("With only single node", func() {
				r := strings.NewReader("<div>Test</div>")

				doc, err := NewDocument("http://test/single", r)

				Convey("No error was returned", func() {
					So(err, ShouldBeNil)
				})
				Convey("The title is empty", func() {
					So(doc.Title, ShouldBeBlank)
				})
				Convey("There are no meta tags", func() {
					So(len(doc.Meta), ShouldEqual, 0)
				})
				Convey("There body field contains the body node", func() {
					So(doc.Body, ShouldNotBeNil)
					So(doc.Body.Type, ShouldEqual, html.ElementNode)
					So(doc.Body.Data, ShouldEqual, "body")
				})
				Convey("There doc field contains the html node", func() {
					So(doc.Doc, ShouldNotBeNil)
					So(doc.Doc.Type, ShouldEqual, html.DocumentNode)
				})
			})
			Convey("With multiple nodes", func() {
				// Load html from file
				r, err := os.Open("../test/simple.html")
				if err != nil {
					panic(err)
				}

				doc, err := NewDocument("http://test/multiple", r)

				Convey("No error was returned", func() {
					So(err, ShouldBeNil)
				})
				Convey("The title is 'Starter Template for Bootstrap'", func() {
					So(doc.Title, ShouldEqual, "Starter Template for Bootstrap")
				})
				Convey("There are 3 meta tags", func() {
					So(len(doc.Meta), ShouldEqual, 3)
				})
				Convey("There body field contains the body node", func() {
					So(doc.Body, ShouldNotBeNil)
					So(doc.Body.Type, ShouldEqual, html.ElementNode)
					So(doc.Body.Data, ShouldEqual, "body")
				})
				Convey("There doc field contains the html node", func() {
					So(doc.Doc, ShouldNotBeNil)
					So(doc.Doc.Type, ShouldEqual, html.DocumentNode)
				})
			})
		})
	})
}
