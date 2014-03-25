package base

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/dancannon/gofetch/document"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSetup(t *testing.T) {
	var e *BaseExtractor

	Convey("Subject: Setup Base ExtractValuesor", t, func() {
		Convey("When the ExtractValuesor is setup", func() {
			e = &BaseExtractor{}
			err := e.Setup(map[string]interface{}{})

			Convey("No error is returned", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}

func TestExtractValues(t *testing.T) {
	var e *BaseExtractor

	Convey("Subject: ExtractValues values from page", t, func() {
		Convey("When the ExtractValuesor is setup", func() {
			Convey("With a page that support oembed", func() {
				e = &BaseExtractor{}
				err := e.Setup(map[string]interface{}{})
				if err != nil {
					t.Fatal(err.Error())
				}

				Convey("And the ExtractValuesor is run", func() {
					ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						f, err := os.Open("../../test/oembed_photo.json")
						if err != nil {
							http.Error(w, err.Error(), http.StatusInternalServerError)
						}

						_, err = io.Copy(w, f)
						if err != nil {
							http.Error(w, err.Error(), http.StatusInternalServerError)
						}
					}))

					doc, err := document.NewDocument("", strings.NewReader(`
						<html><head><link rel="alternate" type="application/json+oembed" href="`+ts.URL+`/?" /></head></html>
					`))
					if err != nil {
						t.Fatal(err.Error())
					}

					_, _, err = e.ExtractValues(*doc)
					Convey("No error was returned", func() {
						So(err, ShouldBeNil)
					})
				})
			})
			Convey("With a page that supports opengraph", func() {
				e = &BaseExtractor{}
				err := e.Setup(map[string]interface{}{})
				if err != nil {
					t.Fatal(err.Error())
				}

				Convey("And the ExtractValuesor is run", func() {
					f, err := os.Open("../../test/opengraph_text.html")
					if err != nil {
						t.Fatal(err.Error())
					}
					doc, err := document.NewDocument("", f)
					if err != nil {
						t.Fatal(err.Error())
					}

					_, _, err = e.ExtractValues(*doc)
					Convey("No error was returned", func() {
						So(err, ShouldBeNil)
					})
				})
			})
			Convey("With a page that does not support OEmbed or OpenGraph", func() {
				e = &BaseExtractor{}
				err := e.Setup(map[string]interface{}{})
				if err != nil {
					t.Fatal(err.Error())
				}

				Convey("And the ExtractValuesor is run", func() {
					f, err := os.Open("../../test/simple.html")
					if err != nil {
						t.Fatal(err.Error())
					}
					doc, err := document.NewDocument("", f)
					if err != nil {
						t.Fatal(err.Error())
					}

					_, _, err = e.ExtractValues(*doc)
					Convey("An error was returned", func() {
						So(err, ShouldBeNil)
					})
				})
			})
		})
	})
}
