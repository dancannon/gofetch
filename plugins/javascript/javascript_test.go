package javascript

import (
	"github.com/dancannon/gofetch/document"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSetup(t *testing.T) {
	var e *JavaScriptExtractor

	Convey("Subject: Setup Selector extractor", t, func() {
		Convey("When the extractor is setup", func() {
			Convey("With a valid script", func() {
				e = &JavaScriptExtractor{}
				err := e.Setup(map[string]interface{}{
					"script": `
						setPageType("unknown");
						setValue("test");

						return 0;
					`,
				})

				Convey("No error is returned", func() {
					So(err, ShouldBeNil)
				})
			})
			Convey("With no script", func() {
				e = &JavaScriptExtractor{}
				err := e.Setup(map[string]interface{}{})

				Convey("An error is returned", func() {
					So(err, ShouldNotBeNil)
				})
			})
		})
	})
}

func TestExtract(t *testing.T) {
	var e *JavaScriptExtractor

	Convey("Subject: Extract value from page", t, func() {
		Convey("When the extractor is setup", func() {
			Convey("With a valid script", func() {
				e = &JavaScriptExtractor{}
				err := e.Setup(map[string]interface{}{
					"script": `
						setValue("test");
					`,
				})
				Convey("No error was returned", func() {
					So(err, ShouldBeNil)
				})

				Convey("When the extractor is run", func() {
					doc, err := document.NewDocument("", strings.NewReader(""))
					if err != nil {
						t.Fatal(err.Error())
					}

					res, err := e.Extract(*doc)
					Convey("No error was returned", func() {
						So(err, ShouldBeNil)
					})
					Convey("The returned value was correct", func() {
						So(res, ShouldEqual, "test")
					})
				})
			})
			Convey("With a script with a syntax error", func() {
				e = &JavaScriptExtractor{}
				err := e.Setup(map[string]interface{}{
					"script": "setValue(\"test);",
				})
				Convey("No error was returned", func() {
					So(err, ShouldBeNil)
				})

				Convey("When the extractor is run", func() {
					doc, err := document.NewDocument("", strings.NewReader(""))
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

func TestExtractValues(t *testing.T) {
	var e *JavaScriptExtractor

	Convey("Subject: Extract values from page", t, func() {
		Convey("When the extractor is setup", func() {
			Convey("With a valid script", func() {
				e = &JavaScriptExtractor{}
				err := e.Setup(map[string]interface{}{
					"script": `
						setValue("test");
						setPageType("type");
					`,
				})
				Convey("No error was returned", func() {
					So(err, ShouldBeNil)
				})

				Convey("When the extractor is run", func() {
					doc, err := document.NewDocument("", strings.NewReader(""))
					if err != nil {
						t.Fatal(err.Error())
					}

					res, typ, err := e.ExtractValues(*doc)
					Convey("No error was returned", func() {
						So(err, ShouldBeNil)
					})
					Convey("The returned value was correct", func() {
						So(res, ShouldEqual, "test")
						So(typ, ShouldEqual, "type")
					})
				})
			})
			Convey("With a script with a syntax error", func() {
				e = &JavaScriptExtractor{}
				err := e.Setup(map[string]interface{}{
					"script": "setValue(\"test);",
				})
				Convey("No error was returned", func() {
					So(err, ShouldBeNil)
				})

				Convey("When the extractor is run", func() {
					doc, err := document.NewDocument("", strings.NewReader(""))
					if err != nil {
						t.Fatal(err.Error())
					}

					_, _, err = e.ExtractValues(*doc)
					Convey("An error was returned", func() {
						So(err, ShouldNotBeNil)
					})
				})
			})
		})
	})
}
