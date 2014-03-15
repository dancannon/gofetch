package js

import (
	"github.com/dancannon/gofetch/document"
	"github.com/dancannon/gofetch/sandbox"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSandbox(t *testing.T) {
	var sb sandbox.Sandbox

	Convey("Subject: JavaScript sandbox", t, func() {
		Convey("Given a valid script", func() {
			sbc := sandbox.SandboxConfig{
				Script: `function processMessage() {
					setPageType("unknown");
					setValue("test");

					return 0;
				}`,
			}

			Convey("When the sandbox is created", func() {
				var err error
				sb, err = NewSandbox(sbc)
				Convey("No error was returned", func() {
					So(err, ShouldBeNil)
				})
				err = sb.Init()
				Convey("No error was returned", func() {
					So(err, ShouldBeNil)
				})

				Convey("And when a message is processed", func() {
					msg := sandbox.SandboxMessage{}
					err = sb.ProcessMessage(&msg)

					Convey("The result should be correct", func() {
						So(msg.PageType, ShouldEqual, "unknown")
						So(msg.Value, ShouldEqual, "test")
					})
					Convey("No error was returned", func() {
						So(err, ShouldBeNil)
					})
				})
			})
		})
		Convey("Given a script with no ProcessMessage function", func() {
			Convey("When the sandbox is created", func() {
				Convey("An error was returned", func() {
					sbc := sandbox.SandboxConfig{
						Script: ``,
					}

					Convey("When the sandbox is created", func() {
						var err error
						sb, err = NewSandbox(sbc)
						Convey("No error was returned", func() {
							So(err, ShouldBeNil)
						})
						err = sb.Init()
						Convey("No error was returned", func() {
							So(err, ShouldBeNil)
						})

						Convey("And when a message is processed", func() {
							msg := sandbox.SandboxMessage{}
							err = sb.ProcessMessage(&msg)

							Convey("An error was returned", func() {
								So(err, ShouldNotBeNil)
							})
						})
					})
				})
			})
		})
		Convey("Given a script with syntax errors", func() {
			Convey("When the sandbox is created", func() {
				Convey("An error was returned", func() {
					sbc := sandbox.SandboxConfig{
						Script: `function processMessage() {
							setPageType("unknown");
							setValue("test);

							return 0;
						}`,
					}

					Convey("When the sandbox is created", func() {
						var err error
						sb, err = NewSandbox(sbc)
						Convey("An error was returned", func() {
							So(err, ShouldNotBeNil)
						})
						err = sb.Init()
						Convey("An error was returned", func() {
							So(err, ShouldNotBeNil)
						})
					})
				})
			})
		})
	})
}

func TestGetValue(t *testing.T) {
	var sb sandbox.Sandbox

	Convey("Subject: JavaScript sandbox", t, func() {
		Convey("Given a valid document", func() {
			f, err := os.Open("../../test/simple.html")
			if err != nil {
				t.Fatal(err.Error())
			}
			doc, err := document.NewDocument("url", f)
			if err != nil {
				t.Fatal(err.Error())
			}

			Convey("And a script that tests getValue('PageType')", func() {
				sbc := sandbox.SandboxConfig{
					Script: `function processMessage() {
							setValue(getValue('PageType') === 'pagetype');

							return 0;
						}`,
				}

				Convey("When the sandbox is created", func() {
					var err error
					sb, err = NewSandbox(sbc)
					Convey("No error was returned", func() {
						So(err, ShouldBeNil)
					})
					err = sb.Init()
					Convey("No error was returned", func() {
						So(err, ShouldBeNil)
					})

					Convey("And when a message is processed", func() {
						msg := sandbox.SandboxMessage{
							PageType: "pagetype",
							Value:    "value",
							Document: *doc,
						}
						err = sb.ProcessMessage(&msg)

						Convey("The result should be correct", func() {
							So(msg.Value, ShouldEqual, true)
						})
						Convey("No error was returned", func() {
							So(err, ShouldBeNil)
						})
					})
				})
			})
			Convey("And a script that tests getValue('Value')", func() {
				sbc := sandbox.SandboxConfig{
					Script: `function processMessage() {
							setValue(getValue('Value') === 'value');

							return 0;
						}`,
				}

				Convey("When the sandbox is created", func() {
					var err error
					sb, err = NewSandbox(sbc)
					Convey("No error was returned", func() {
						So(err, ShouldBeNil)
					})
					err = sb.Init()
					Convey("No error was returned", func() {
						So(err, ShouldBeNil)
					})

					Convey("And when a message is processed", func() {
						msg := sandbox.SandboxMessage{
							PageType: "pagetype",
							Value:    "value",
							Document: *doc,
						}
						err = sb.ProcessMessage(&msg)

						Convey("The result should be correct", func() {
							So(msg.Value, ShouldEqual, true)
						})
						Convey("No error was returned", func() {
							So(err, ShouldBeNil)
						})
					})
				})
			})
			Convey("And a script that tests getValue('Document.*')", func() {
				sbc := sandbox.SandboxConfig{
					Script: `function processMessage() {
							setValue({
								"meta": getValue('Document.Meta'),
								"doc": getValue('Document.Doc') !== null,
								"body": getValue('Document.Body') !== null,
								"url": getValue('Document.URL'),
								"title": getValue('Document.Title'),
							});

							return 0;
						}
						`,
				}

				Convey("When the sandbox is created", func() {
					var err error
					sb, err = NewSandbox(sbc)
					Convey("No error was returned", func() {
						So(err, ShouldBeNil)
					})
					err = sb.Init()
					Convey("No error was returned", func() {
						So(err, ShouldBeNil)
					})

					Convey("And when a message is processed", func() {
						msg := sandbox.SandboxMessage{
							PageType: "pagetype",
							Value:    "value",
							Document: *doc,
						}
						err = sb.ProcessMessage(&msg)

						Convey("The result should be correct", func() {
							So(msg.Value, ShouldResemble, map[string]interface{}{
								"url":   "url",
								"title": "Starter Template for Bootstrap",
								"meta": []map[string]string{
									map[string]string{
										"charset": "utf-8",
									},
									map[string]string{
										"name":    "description",
										"content": "description",
									},
									map[string]string{
										"name":    "author",
										"content": "author",
									},
								},
								"doc":  true,
								"body": true,
							})
						})
						Convey("No error was returned", func() {
							So(err, ShouldBeNil)
						})
					})
				})
				Convey("And a script that tests getValue('unknown') is undefined", func() {
					sbc := sandbox.SandboxConfig{
						Script: `function processMessage() {
							setValue(typeof getValue('unknown') === 'undefined');

							return 0;
						}`,
					}

					Convey("When the sandbox is created", func() {
						var err error
						sb, err = NewSandbox(sbc)
						Convey("No error was returned", func() {
							So(err, ShouldBeNil)
						})
						err = sb.Init()
						Convey("No error was returned", func() {
							So(err, ShouldBeNil)
						})

						Convey("And when a message is processed", func() {
							msg := sandbox.SandboxMessage{
								PageType: "pagetype",
								Value:    "value",
								Document: *doc,
							}
							err = sb.ProcessMessage(&msg)

							Convey("The result should be correct", func() {
								So(msg.Value, ShouldEqual, true)
							})
							Convey("No error was returned", func() {
								So(err, ShouldBeNil)
							})
						})
					})
				})
				Convey("And a script that contains an infinite loop", func() {
					sbc := sandbox.SandboxConfig{
						Script: `
						while(true) {

						}`,
					}

					Convey("When the sandbox is created", func() {
						var err error
						sb, err = NewSandbox(sbc)
						Convey("An error was returned", func() {
							So(err, ShouldNotBeNil)
						})
						err = sb.Init()
						Convey("An error was returned", func() {
							So(err, ShouldNotBeNil)
						})
					})
				})
				Convey("And a script that has a processMessage which contains an infinite loop", func() {
					sbc := sandbox.SandboxConfig{
						Script: `function processMessage() {
							while(true) {

							}
						}`,
					}

					Convey("When the sandbox is created", func() {
						var err error
						sb, err = NewSandbox(sbc)
						Convey("No error was returned", func() {
							So(err, ShouldBeNil)
						})
						err = sb.Init()
						Convey("No error was returned", func() {
							So(err, ShouldBeNil)
						})

						Convey("And when a message is processed", func() {
							msg := sandbox.SandboxMessage{
								PageType: "pagetype",
								Value:    "value",
								Document: *doc,
							}
							err = sb.ProcessMessage(&msg)
							Convey("An error was returned", func() {
								So(err, ShouldNotBeNil)
							})
						})
					})
				})
			})
		})
	})
}
