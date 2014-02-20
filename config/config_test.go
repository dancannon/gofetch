package config

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLoadConfig(t *testing.T) {
	Convey("Subject: Load config from file 'test/data/config.json'", t, func() {
		Convey("When the config is loaded from file 'test/data/config.json'", func() {
			conf, err := LoadConfig("../test/data/config.json")

			Convey("No error was returned", func() {
				So(err, ShouldBeNil)
			})
			Convey("There are 7 types", func() {
				So(len(conf.Types), ShouldEqual, 7)
			})
			Convey("There are 2 rules", func() {
				So(len(conf.Rules), ShouldEqual, 2)
			})
		})
		Convey("When the config is loaded from non-existant file", func() {
			_, err := LoadConfig("../test/data/not-here-config.json")

			Convey("An error was returned", func() {
				So(err, ShouldNotBeNil)
			})
		})
		Convey("When the config is loaded from invalid file", func() {
			_, err := LoadConfig("../test/data/invalid-config.json")

			Convey("An error was returned", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}
