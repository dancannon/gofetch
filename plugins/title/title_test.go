package title

import (
	"github.com/dancannon/gofetch/document"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSetup(t *testing.T) {
	var e *TitleExtractor

	Convey("Subject: Setup UrlMapper extractor", t, func() {
		Convey("When the extractor is setup", func() {
			e = &TitleExtractor{}
			err := e.Setup(map[string]interface{}{})

			Convey("No error is returned", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}

func TestExtract(t *testing.T) {
	var e *TitleExtractor

	Convey("Subject: Extract value from page", t, func() {
		Convey("When the extractor is setup", func() {
			e = &TitleExtractor{}
			err := e.Setup(map[string]interface{}{})

			Convey("No error is returned", func() {
				So(err, ShouldBeNil)
			})

			Convey("And the extractor is run", func() {
				Convey("on a document with a title element and no header elements", func() {
					doc, err := document.NewDocument("", strings.NewReader(`
						<html>
							<head>
								<title>Test</title>
							</head>
						</html>
					`))
					if err != nil {
						t.Fatal(err.Error())
					}

					res, err := e.Extract(*doc)
					Convey("No error was returned", func() {
						So(err, ShouldBeNil)
					})
					Convey("The result should be valid", func() {
						So(res, ShouldEqual, "Test")
					})
				})
				Convey("on a document with a title element and a header element", func() {
					Convey("that is a substring of the title", func() {
						doc, err := document.NewDocument("", strings.NewReader(`
								<html>
									<head>
										<title>Test Title - Site</title>
									</head>
									<body>
										<h1>Test Title</h1>
									</body>
								</html>
							`))
						if err != nil {
							t.Fatal(err.Error())
						}

						res, err := e.Extract(*doc)
						Convey("No error was returned", func() {
							So(err, ShouldBeNil)
						})
						Convey("The result should be valid", func() {
							So(res, ShouldEqual, "Test Title")
						})
					})
					Convey("that is not substring of the title", func() {
						doc, err := document.NewDocument("", strings.NewReader(`
								<html>
									<head>
										<title>Test Title</title>
									</head>
									<body>
										<h1>Other</h1>
									</body>
								</html>
							`))
						if err != nil {
							t.Fatal(err.Error())
						}

						res, err := e.Extract(*doc)
						Convey("No error was returned", func() {
							So(err, ShouldBeNil)
						})
						Convey("The result should be valid", func() {
							So(res, ShouldEqual, "Test Title")
						})
					})
				})
				Convey("on a document with a title element and multiple header element", func() {
					Convey("that is a substring of the title", func() {
						doc, err := document.NewDocument("", strings.NewReader(`
								<html>
									<head>
										<title>Test Title</title>
									</head>
									<body>
										<h1>Test Title</h1>
										<h1>Other</h1>
										<h2>Test Title - Site</h2>
									</body>
								</html>
							`))
						if err != nil {
							t.Fatal(err.Error())
						}

						res, err := e.Extract(*doc)
						Convey("No error was returned", func() {
							So(err, ShouldBeNil)
						})
						Convey("The result should be valid", func() {
							So(res, ShouldEqual, "Test Title")
						})
					})
				})
			})
		})
	})
}
