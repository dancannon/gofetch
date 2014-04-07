package gofetch

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/dancannon/gofetch/config"
)

var conf config.Config
var remoteTests = []struct {
	url              string // input URL
	expectedPageType string // expected page type
}{
	{"http://getbootstrap.com/examples/starter-template/", "text"},
	{"http://getbootstrap.com/examples/jumbotron/", "text"},
	{"http://getbootstrap.com/examples/carousel/", "text"},
	{"http://www.theguardian.com/technology/2013/nov/01/caa-easa-electronic-devices-flight-take-off-landing", "general"},
	{"http://www.birmingham.ac.uk/index.aspx", "text"},
	{"https://www.google.co.uk/?gws_rd=cr&ei=IMtzUuLkI-Hb0QX-woD4CA#q=test", "text"},
	{"https://github.com/dancannon/gorethink/issues/51", "general"},
	{"http://www.youtube.com/watch?v=-UUx10KOWIE", "video"},
	{"http://blog.danielcannon.co.uk/2012/07/02/building-a-real-application-with-backbonejs", "text"},
	{"http://www.bbc.co.uk/news/uk-26111598", "text"},
	{"http://imgur.com/7T7MrBc", "image"},
	{"http://www.flickr.com/photos/bees/2341623661/", "image"},
	{"http://stackoverflow.com/questions/7438323/method-requires-pointer-receiver-in-go-programming-language/", "general"},
}
var localTests = []struct {
	inputFile        string      // input URL
	expectedPageType string      // expected page type
	expectedContent  interface{} // expected content
}{
	{
		"test/plain_text.txt",
		"text",
		map[string]interface{}{"title": "", "text": "Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.\n"},
	},
	{
		"test/image.gif",
		"raw",
		map[string]interface{}{
			"mime_type": "image/gif",
		},
	},
	{
		"test/simple.html",
		"text",
		nil,
	},
	{
		"test/blog.html",
		"text",
		nil,
	},
}

func init() {
	var err error
	conf, err = config.LoadConfig("config.json")
	if err != nil {
		panic(err.Error())
	}
}

func TestFetchRemoteUrl(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	var fetcher *Fetcher

	for _, tt := range remoteTests {
		Convey(fmt.Sprintf("Subject: Test fetch content from URL %s", tt.url), t, func() {
			Convey("Given a new fetcher instance", func() {
				fetcher = NewFetcher(conf)

				Convey("When the content is extracted", func() {
					res, err := fetcher.Fetch(tt.url)

					Convey("No error was returned", func() {
						So(err, ShouldBeNil)
					})
					Convey(fmt.Sprintf("The result page type should be %s", tt.expectedPageType), func() {
						So(res.PageType, ShouldEqual, tt.expectedPageType)
					})
					Convey("The result should not be empty", func() {
						So(res.Content, ShouldNotBeNil)
					})
				})
			})
		})
	}
}

func TestFetchMockedServer(t *testing.T) {
	for _, tt := range localTests {
		Convey(fmt.Sprintf("Subject: Test fetch content using content from %s", tt.inputFile), t, func() {
			var fetcher *Fetcher

			Convey("Given a new fetcher instance", func() {
				fetcher = NewFetcher(conf)

				Convey("and the HTTP server has been started", func() {
					// Start a server that loads the content from the input file
					ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						f, err := os.Open(tt.inputFile)
						if err != nil {
							http.Error(w, err.Error(), http.StatusInternalServerError)
						}

						_, err = io.Copy(w, f)
						if err != nil {
							http.Error(w, err.Error(), http.StatusInternalServerError)
						}
					}))

					Convey("When the content is extracted", func() {
						res, err := fetcher.Fetch(ts.URL)

						Convey("No error was returned", func() {
							So(err, ShouldBeNil)
						})
						Convey(fmt.Sprintf("The result page type should be %s", tt.expectedPageType), func() {
							So(res.PageType, ShouldEqual, tt.expectedPageType)
						})

						if tt.expectedContent != nil {
							Convey("The result should be equal to the expected value", func() {
								So(res.Content, ShouldResemble, tt.expectedContent)
							})
						}
					})

					Reset(func() {
						ts.Close()
					})
				})
			})
		})
	}

	Convey("Subject: Test fetch content using a rule that sets simple values", t, func() {
		var fetcher *Fetcher

		Convey("Given a new fetcher instance", func() {
			fetcher = NewFetcher(conf)

			Convey("and the custom rule has been added", func() {
				rule := config.Rule{
					Type:       "unknown",
					UrlPattern: ".*",
					Priority:   100,
					Values: []interface{}{
						map[string]interface{}{
							"name":  "a",
							"type":  "value",
							"value": "A",
						},
						map[string]interface{}{
							"name": "b",
							"type": "values",
							"value": []interface{}{
								map[string]interface{}{
									"name":  "1",
									"type":  "value",
									"value": "B1",
								},
								map[string]interface{}{
									"name":  "2",
									"type":  "value",
									"value": "B2",
								},
							},
						},
					},
				}
				fetcher.Config.Rules = []config.Rule{rule}

				Convey("and the HTTP server has been started", func() {
					// Start a server that loads the content from the input file
					ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						f, err := os.Open("test/blog.html")
						if err != nil {
							http.Error(w, err.Error(), http.StatusInternalServerError)
						}

						_, err = io.Copy(w, f)
						if err != nil {
							http.Error(w, err.Error(), http.StatusInternalServerError)
						}
					}))

					Convey("When the content is extracted", func() {
						res, err := fetcher.Fetch(ts.URL)

						Convey("No error was returned", func() {
							So(err, ShouldBeNil)
						})
						Convey(fmt.Sprintf("The result page type should be unknown"), func() {
							So(res.PageType, ShouldEqual, "unknown")
						})
						Convey("The result should be equal to the expected value", func() {
							So(res.Content, ShouldResemble, map[string]interface{}{
								"a": "A",
								"b": map[string]interface{}{
									"1": "B1",
									"2": "B2",
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

	Convey("Subject: Test fetch content using a rule that uses extractors", t, func() {
		var fetcher *Fetcher

		Convey("Given a new fetcher instance", func() {
			fetcher = NewFetcher(conf)

			Convey("and the custom rule has been added", func() {
				rule := config.Rule{
					Type:       "text",
					UrlPattern: ".*",
					Priority:   100,
					Values: []interface{}{
						map[string]interface{}{
							"name": "title",
							"id":   "selector",
							"type": "extractor",
							"params": map[string]interface{}{
								"selector": ".blog-post-title",
							},
						},
						map[string]interface{}{
							"name": "text",
							"id":   "selector",
							"type": "extractor",
							"params": map[string]interface{}{
								"selector": ".blog-post",
							},
						},
					},
				}
				fetcher.Config.Rules = []config.Rule{rule}

				Convey("and the HTTP server has been started", func() {
					// Start a server that loads the content from the input file
					ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						f, err := os.Open("test/blog.html")
						if err != nil {
							http.Error(w, err.Error(), http.StatusInternalServerError)
						}

						_, err = io.Copy(w, f)
						if err != nil {
							http.Error(w, err.Error(), http.StatusInternalServerError)
						}
					}))

					Convey("When the content is extracted", func() {
						res, err := fetcher.Fetch(ts.URL)

						Convey("No error was returned", func() {
							So(err, ShouldBeNil)
						})
						Convey(fmt.Sprintf("The result page type should be text"), func() {
							So(res.PageType, ShouldEqual, "text")
						})
						Convey("The result should be equal to the expected value", func() {
							So(res.Content.(map[string]interface{})["title"], ShouldEqual, "Sample blog post\nAnother blog post\nNew feature")
						})
					})

					Reset(func() {
						ts.Close()
					})
				})
			})
		})
	})

	Convey("Subject: Test fetch content from non-existant file", t, func() {
		var fetcher *Fetcher

		Convey("Given a new fetcher instance", func() {
			fetcher = NewFetcher(conf)

			Convey("and the HTTP server has been started", func() {
				// Start a server that loads the content from the input file
				ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					http.NotFound(w, r)
				}))

				Convey("When the content is extracted", func() {
					res, err := fetcher.Fetch(ts.URL)

					Convey("No error was returned", func() {
						So(err, ShouldBeNil)
					})
					Convey("A 404 status code was returned", func() {
						So(res.Response.StatusCode, ShouldEqual, 404)
						So(res.Response.Status, ShouldEqual, "404 Not Found")
					})
				})

				Reset(func() {
					ts.Close()
				})
			})
		})
	})

	Convey("Subject: Test fetch content from a page returning an error", t, func() {
		var fetcher *Fetcher

		Convey("Given a new fetcher instance", func() {
			fetcher = NewFetcher(conf)

			Convey("and the HTTP server has been started", func() {
				// Start a server that loads the content from the input file
				ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					http.Error(w, "Test error", http.StatusInternalServerError)
				}))

				Convey("When the content is extracted", func() {
					res, err := fetcher.Fetch(ts.URL)

					Convey("No error was returned", func() {
						So(err, ShouldBeNil)
					})
					Convey("A 500 status code was returned", func() {
						So(res.Response.StatusCode, ShouldEqual, 500)
						So(res.Response.Status, ShouldEqual, "500 Internal Server Error")
					})
				})

				Reset(func() {
					ts.Close()
				})
			})
		})
	})

	Convey("Subject: Test fetch content from an invalid URL", t, func() {
		var fetcher *Fetcher

		Convey("Given a new fetcher instance", func() {
			fetcher = NewFetcher(conf)

			Convey("When the content is extracted", func() {
				_, err := fetcher.Fetch("not-a-url")

				Convey("An error was returned", func() {
					So(err, ShouldNotBeNil)
				})
			})
		})
	})
}

func TestValidation(t *testing.T) {
	Convey("Subject: Test validate using a type that allows extra fields", t, func() {
		var fetcher *Fetcher
		var result Result

		Convey("Given a new fetcher instance", func() {
			fetcher = NewFetcher(conf)

			Convey("and an empty map", func() {
				result = Result{
					PageType: "unknown",
					Content:  map[string]interface{}{},
				}

				Convey("When the the input is validated", func() {
					err := fetcher.validateResult(result)

					Convey("An error was returned", func() {
						So(err, ShouldBeNil)
					})
				})
			})
			Convey("and a map that contains extra fields", func() {
				result = Result{
					PageType: "unknown",
					Content: map[string]interface{}{
						"title": "title",
						"test":  "test",
					},
				}

				Convey("When the the input is validated", func() {
					err := fetcher.validateResult(result)

					Convey("An error was returned", func() {
						So(err, ShouldBeNil)
					})
				})
			})
			Convey("and a map that contains a child map with extra fields", func() {
				result = Result{
					PageType: "unknown",
					Content: map[string]interface{}{
						"title": "title",
						"test": map[string]interface{}{
							"test2": "test2",
						},
					},
				}

				Convey("When the the input is validated", func() {
					err := fetcher.validateResult(result)

					Convey("An error was returned", func() {
						So(err, ShouldBeNil)
					})
				})
			})
			Convey("and a valid result", func() {
				result = Result{
					PageType: "unknown",
					Content: map[string]interface{}{
						"title":   "title",
						"content": "content",
					},
				}

				Convey("When the the input is validated", func() {
					err := fetcher.validateResult(result)

					Convey("An error was returned", func() {
						So(err, ShouldBeNil)
					})
				})
			})
		})
	})
	Convey("Subject: Test validate using a type that does not allow extra fields", t, func() {
		var fetcher *Fetcher
		var result Result

		Convey("Given a new fetcher instance", func() {
			fetcher = NewFetcher(conf)

			Convey("and a nil result", func() {
				result = Result{
					PageType: "general",
					Content:  nil,
				}

				Convey("When the the input is validated", func() {
					err := fetcher.validateResult(result)

					Convey("An error was returned", func() {
						So(err, ShouldNotBeNil)
						So(err.Error(), ShouldEqual, "The result is not of the correct type")
					})
				})
			})
			Convey("and an empty map", func() {
				result = Result{
					PageType: "general",
					Content:  map[string]interface{}{},
				}

				Convey("When the the input is validated", func() {
					err := fetcher.validateResult(result)

					Convey("An error was returned", func() {
						So(err, ShouldNotBeNil)
					})
				})
			})
			Convey("and a map that contains extra fields", func() {
				result = Result{
					PageType: "general",
					Content: map[string]interface{}{
						"title": "title",
						"test":  "test",
					},
				}

				Convey("When the the input is validated", func() {
					err := fetcher.validateResult(result)

					Convey("An error was returned", func() {
						So(err, ShouldNotBeNil)
					})
				})
			})
			Convey("and a map that contains a child map", func() {
				Convey("with a missing fields", func() {
					result = Result{
						PageType: "text",
						Content: map[string]interface{}{
							"title": "title",
						},
					}

					Convey("When the the input is validated", func() {
						err := fetcher.validateResult(result)

						Convey("An error was returned", func() {
							So(err, ShouldNotBeNil)
							So(err.Error(), ShouldEqual, "The type text requires the field text")
						})
					})
				})
				Convey("with extra fields", func() {
					result = Result{
						PageType: "text",
						Content: map[string]interface{}{
							"title": "title",
							"author": map[string]interface{}{
								"name": "name",
								"test": "test",
							},
							"text": "text",
						},
					}

					Convey("When the the input is validated", func() {
						err := fetcher.validateResult(result)

						Convey("An error was returned", func() {
							So(err, ShouldNotBeNil)
							So(err.Error(), ShouldEqual, "The type text does not contain the field test")
						})
					})
				})
			})
			Convey("and a valid result", func() {
				result = Result{
					PageType: "text",
					Content: map[string]interface{}{
						"title": "title",
						"author": map[string]interface{}{
							"name": "name",
						},
						"text": "text",
					},
				}

				Convey("When the the input is validated", func() {
					err := fetcher.validateResult(result)

					Convey("An error was returned", func() {
						So(err, ShouldBeNil)
					})
				})
			})
		})
	})
}
