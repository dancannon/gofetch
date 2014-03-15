package js

import (
	"github.com/dancannon/gofetch/document"
	"github.com/dancannon/gofetch/sandbox"
	"github.com/davecgh/go-spew/spew"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSandbox(t *testing.T) {
	var sb sandbox.Sandbox

	Convey("Subject: JavaScript sandbox", t, func() {
		Convey("Given a valid script", func() {
			sbc := sandbox.SandboxConfig{
				Script: `
					setPageType("unknown");
					setValue("test");
				`,
			}

			Convey("and an instance of a sandbox", func() {
				sb = NewSandbox(sbc)

				Convey("When the sandbox is initialized", func() {
					err := sb.Init()

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
		})
		Convey("Given a script with syntax errors", func() {
			Convey("When the sandbox is created", func() {
				sbc := sandbox.SandboxConfig{
					Script: `
						setPageType("unknown");
						setValue("test); // No closing quote
					`,
				}

				Convey("and an instance of a sandbox", func() {
					sb = NewSandbox(sbc)

					Convey("When the sandbox is initialized", func() {
						err := sb.Init()

						Convey("No error was returned", func() {
							So(err, ShouldBeNil)
						})

						Convey("And when a message is processed", func() {
							msg := sandbox.SandboxMessage{}
							err = sb.ProcessMessage(&msg)

							Convey("Ab error was returned", func() {
								So(err, ShouldNotBeNil)
							})
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

			Convey("And a script that tests getPageType()", func() {
				sbc := sandbox.SandboxConfig{
					Script: `
						setValue(getPageType() === 'pagetype');
					`,
				}

				Convey("and an instance of a sandbox", func() {
					sb = NewSandbox(sbc)

					Convey("When the sandbox is initialized", func() {
						err := sb.Init()

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
			})
			Convey("And a script that tests getValue()", func() {
				sbc := sandbox.SandboxConfig{
					Script: `
						setValue(getValue() === 'value');
					`,
				}

				Convey("and an instance of a sandbox", func() {
					sb = NewSandbox(sbc)

					Convey("When the sandbox is initialized", func() {
						err := sb.Init()

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
			})
			Convey("And a script that tests getValue('Document.*')", func() {
				sbc := sandbox.SandboxConfig{
					Script: `
							setValue({
								"meta": document.Meta,
								"doc": document.Doc !== null,
								"body": document.Body !== null,
								"url": document.URL.String(),
								"title": document.Title,
							});
						`,
				}

				Convey("and an instance of a sandbox", func() {
					sb = NewSandbox(sbc)

					Convey("When the sandbox is initialized", func() {
						err := sb.Init()

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
								spew.Dump(msg.Value)
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
				})
			})
		})
	})
}

func TestInterrupt(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	var sb sandbox.Sandbox

	Convey("Subject: JavaScript sandbox", t, func() {
		Convey("Given a script that contains an infinite loop", func() {
			sbc := sandbox.SandboxConfig{
				Script: `
						while(true) {

						}
					`,
			}
			Convey("and an instance of a sandbox", func() {
				sb = NewSandbox(sbc)

				Convey("When the sandbox is initialized", func() {
					err := sb.Init()

					Convey("No error was returned", func() {
						So(err, ShouldBeNil)
					})

					Convey("And when a message is processed", func() {
						msg := sandbox.SandboxMessage{}
						err = sb.ProcessMessage(&msg)

						Convey("Ab error was returned", func() {
							So(err, ShouldNotBeNil)
						})
					})
				})
			})
		})
	})
}
