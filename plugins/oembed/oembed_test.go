package oembed

import (
	"github.com/dancannon/gofetch/document"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSetup(t *testing.T) {
	var e *OEmbedExtractor

	Convey("Subject: Setup OEmbed extractor", t, func() {
		e = &OEmbedExtractor{}

		Convey("When the extractor is setup", func() {
			Convey("With a map containing an endpoint URL", func() {
				err := e.Setup(map[string]interface{}{
					"endpoint": "url?url=%s",
				})

				Convey("No error is returned", func() {
					So(err, ShouldBeNil)
				})
				Convey("The extractor should contain a field named extractor which has the value 'url'", func() {
					So(e.endpointFormat, ShouldEqual, "url?url=%s")
				})
			})
			Convey("And an endpoint has already been provided", func() {
				e.endpoint = "url"
				err := e.Setup(map[string]interface{}{})

				Convey("No error is returned", func() {
					So(err, ShouldBeNil)
				})
				Convey("The extractor should contain a field named extractor which has the value 'url'", func() {
					So(e.endpoint, ShouldEqual, "url")
				})
			})
			Convey("with an empty map", func() {
				err := e.Setup(map[string]interface{}{})

				Convey("An error is returned", func() {
					So(err, ShouldNotBeNil)
				})
			})
		})
	})
}

func TestSupports(t *testing.T) {
	var e *OEmbedExtractor
	Convey("Subject: Check if page supports OEmbed", t, func() {
		e = &OEmbedExtractor{}

		Convey("Given a document that does not support OEmbed", func() {
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
		Convey("Given a document that does support OEmbed", func() {
			f, err := os.Open("../../test/oembed.html")
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
	var e *OEmbedExtractor

	Convey("Subject: Extract values from page", t, func() {
		e = &OEmbedExtractor{}

		Convey("And a JSON endpoint has been started", func() {
			// Start a server that loads the content from the input file
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var inputFile string

				switch r.URL.Path {
				case "/oembed_link.html":
					inputFile = "../../test/oembed_link.json"
				case "/oembed_photo.html":
					inputFile = "../../test/oembed_photo.json"
				case "/oembed_video.html":
					inputFile = "../../test/oembed_video.json"
				default:
					http.NotFound(w, r)
					return
				}

				f, err := os.Open(inputFile)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}

				_, err = io.Copy(w, f)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}))

			Convey("and the extractor is setup", func() {
				err := e.Setup(map[string]interface{}{
					"endpoint": ts.URL + "/%s",
				})
				if err != nil {
					t.Fatal(err.Error())
				}

				Convey("And a page that supports oembed is given", func() {
					Convey("And that page is a photo", func() {
						doc, err := document.NewDocument("oembed_photo.html", strings.NewReader(""))
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
								"width":  240,
								"height": 160,
								"url":    "Url",
								"author": map[string]interface{}{
									"name": "Author",
									"url":  "Author Url",
								},
							})
						})
					})
					Convey("And that page is a video", func() {
						doc, err := document.NewDocument("oembed_video.html", strings.NewReader(""))
						if err != nil {
							t.Fatal(err.Error())
						}

						res, typ, err := e.ExtractValues(*doc)
						Convey("No error was returned", func() {
							So(err, ShouldBeNil)
						})
						Convey("The type should be 'video'", func() {
							So(typ, ShouldEqual, "video")
						})
						Convey("The result should be valid", func() {
							So(res, ShouldResemble, map[string]interface{}{
								"title": "Title",
								"author": map[string]interface{}{
									"name": "Author",
									"url":  "Author Url",
								},
								"html": "HTML",
							})
						})
					})
					Convey("And that page is an unrecognised type", func() {
						doc, err := document.NewDocument("oembed_link.html", strings.NewReader(""))
						if err != nil {
							t.Fatal(err.Error())
						}

						res, typ, err := e.ExtractValues(*doc)
						Convey("No error was returned", func() {
							So(err, ShouldBeNil)
						})
						Convey("The type should be 'unknown'", func() {
							So(typ, ShouldEqual, "unknown")
						})
						Convey("The result should be valid", func() {
							So(res, ShouldResemble, map[string]interface{}{
								"version":       "1.0",
								"provider_name": "Provider",
								"provider_url":  "Provider Url",
								"author_name":   "Author",
								"author_url":    "Author Url",
								"title":         "Test",
								"type":          "link",
								"html":          "Content",
							})
						})
					})
				})
			})
			Convey("And the endpoint returns a 404 error", func() {
				doc, err := document.NewDocument("oembed_not_found.html", strings.NewReader(""))
				if err != nil {
					t.Fatal(err.Error())
				}

				_, _, err = e.ExtractValues(*doc)
				Convey("No error was returned", func() {
					So(err, ShouldNotBeNil)
				})
			})

			Reset(func() {
				ts.Close()
			})
		})
		Convey("And an XML endpoint has been started", func() {
			// Start a server that loads the content from the input file
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var inputFile string

				switch r.URL.Path {
				case "/oembed_photo.html":
					inputFile = "../../test/oembed_photo.xml"
				default:
					http.NotFound(w, r)
					return
				}

				f, err := os.Open(inputFile)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}

				_, err = io.Copy(w, f)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}))

			Convey("and the extractor is setup", func() {
				err := e.Setup(map[string]interface{}{
					"endpoint": ts.URL + "/%s",
					"format":   "xml",
				})
				if err != nil {
					t.Fatal(err.Error())
				}

				Convey("And a page that supports oembed is given", func() {
					Convey("And that page is a photo", func() {
						doc, err := document.NewDocument("oembed_photo.html", strings.NewReader(""))
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
								"width":  240,
								"height": 160,
								"url":    "Url",
								"author": map[string]interface{}{
									"name": "Author",
									"url":  "Author Url",
								},
								"thumbnail": map[string]interface{}{
									"url":    "Thumbnail Url",
									"width":  75,
									"height": 75,
								},
							})
						})
					})

					Reset(func() {
						ts.Close()
					})
				})
			})
		})
	})
}
