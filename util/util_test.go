package util

import (
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCreateMapFromProps(t *testing.T) {
	Convey("Subject: Create map from properties", t, func() {
		var props map[string]interface{}

		Convey("Given a map of properties", func() {
			props = map[string]interface{}{
				"hello":     "world",
				"foo":       "bar",
				"baz":       "baz",
				"stringint": "123",
				"int":       123,
				"float":     float64(123),
			}

			Convey("When I create a map using a key that does not exist in the map", func() {
				m := CreateMapFromProps(props, map[string]string{
					"Test": "Test",
				})
				Convey("An empty map is returned", func() {
					So(m, ShouldResemble, map[string]interface{}{})
				})
			})
			Convey("When I create a map using a key that does exist in the map", func() {
				m := CreateMapFromProps(props, map[string]string{
					"hello": "hello",
				})
				Convey("The expected map is returned", func() {
					So(m, ShouldResemble, map[string]interface{}{
						"hello": "world",
					})
				})
			})
			Convey("When I create a map using multiple keys that do exist in the map", func() {
				m := CreateMapFromProps(props, map[string]string{
					"hello": "hello",
					"foo":   "foo",
					"int":   "int",
				})
				Convey("The expected map is returned", func() {
					So(m, ShouldResemble, map[string]interface{}{
						"hello": "world",
						"foo":   "bar",
						"int":   123,
					})
				})
			})
			Convey("When I create a map and rename a field", func() {
				m := CreateMapFromProps(props, map[string]string{
					"foo": "baz",
				})
				Convey("The expected map is returned", func() {
					So(m, ShouldResemble, map[string]interface{}{
						"foo": "baz",
					})
				})
			})
			Convey("When I create a map and cast a field", func() {
				Convey("From a string to an int", func() {
					m := CreateMapFromProps(props, map[string]string{
						"int:int": "stringint",
					})
					Convey("The expected map is returned", func() {
						So(m, ShouldResemble, map[string]interface{}{
							"int": 123,
						})
					})
				})
				Convey("From a int to a float", func() {
					m := CreateMapFromProps(props, map[string]string{
						"float:float": "int",
					})
					Convey("The expected map is returned", func() {
						So(m, ShouldResemble, map[string]interface{}{
							"float": float64(123),
						})
					})
				})
				Convey("From a float to an int", func() {
					m := CreateMapFromProps(props, map[string]string{
						"int:int": "float",
					})
					Convey("The expected map is returned", func() {
						So(m, ShouldResemble, map[string]interface{}{
							"int": 123,
						})
					})
				})
				Convey("From a int to an string", func() {
					m := CreateMapFromProps(props, map[string]string{
						"string:string": "int",
					})
					Convey("The expected map is returned", func() {
						So(m, ShouldResemble, map[string]interface{}{
							"string": "123",
						})
					})
				})
				Convey("From a int to an uint", func() {
					m := CreateMapFromProps(props, map[string]string{
						"uint:uint": "int",
					})
					Convey("The expected map is returned", func() {
						So(m, ShouldResemble, map[string]interface{}{
							"uint": uint(123),
						})
					})
				})
			})
		})
	})
}

func TestCat(t *testing.T) {
	Convey("Subject: Casting", t, func() {
		Convey("When casting from a boolean", func() {
			Convey("To a boolean", func() {
				v, err := Cast(true, reflect.TypeOf(true))

				Convey("No error should be returned", func() {
					So(err, ShouldBeNil)
				})
				Convey("The correct value should be returned", func() {
					So(v, ShouldEqual, true)
				})
			})
			Convey("To a string", func() {
				v, err := Cast(true, reflect.TypeOf(""))

				Convey("No error should be returned", func() {
					So(err, ShouldBeNil)
				})
				Convey("The correct value should be returned", func() {
					So(v, ShouldEqual, "true")
				})
			})
			Convey("To an int", func() {
				v, err := Cast(true, reflect.TypeOf(0))

				Convey("No error should be returned", func() {
					So(err, ShouldBeNil)
				})
				Convey("The correct value should be returned", func() {
					So(v, ShouldEqual, 1)
				})
			})
			Convey("To an invalid type", func() {
				var res = struct{}{}
				_, err := Cast(true, reflect.TypeOf(res))

				Convey("An error should be returned", func() {
					So(err, ShouldNotBeNil)
				})
			})
		})
		Convey("When casting from a string", func() {
			Convey("To a boolean", func() {
				Convey("Using a valid string", func() {
					v, err := Cast("True", reflect.TypeOf(true))

					Convey("No error should be returned", func() {
						So(err, ShouldBeNil)
					})
					Convey("The correct value should be returned", func() {
						So(v, ShouldEqual, true)
					})
				})
				Convey("Using an invalid string", func() {
					_, err := Cast("string", reflect.TypeOf(true))

					Convey("An error should be returned", func() {
						So(err, ShouldNotBeNil)
					})
				})
			})
			Convey("To a string", func() {
				v, err := Cast("string", reflect.TypeOf(""))

				Convey("No error should be returned", func() {
					So(err, ShouldBeNil)
				})
				Convey("The correct value should be returned", func() {
					So(v, ShouldEqual, "string")
				})
			})
			Convey("To an int", func() {
				Convey("Using a valid string", func() {
					v, err := Cast("123", reflect.TypeOf(0))

					Convey("No error should be returned", func() {
						So(err, ShouldBeNil)
					})
					Convey("The correct value should be returned", func() {
						So(v, ShouldEqual, 123)
					})
				})
				Convey("Using an invalid string", func() {
					_, err := Cast("string", reflect.TypeOf(0))

					Convey("An error should be returned", func() {
						So(err, ShouldNotBeNil)
					})
				})
			})
			Convey("To an uint", func() {
				Convey("Using a valid string", func() {
					v, err := Cast("123", reflect.TypeOf(uint(0)))

					Convey("No error should be returned", func() {
						So(err, ShouldBeNil)
					})
					Convey("The correct value should be returned", func() {
						So(v, ShouldEqual, uint(123))
					})
				})
				Convey("Using an invalid string", func() {
					_, err := Cast("string", reflect.TypeOf(uint(0)))

					Convey("An error should be returned", func() {
						So(err, ShouldNotBeNil)
					})
				})
			})
			Convey("To a float", func() {
				Convey("Using a valid string", func() {
					v, err := Cast("123.1", reflect.TypeOf(float64(0)))

					Convey("No error should be returned", func() {
						So(err, ShouldBeNil)
					})
					Convey("The correct value should be returned", func() {
						So(v, ShouldEqual, 123.1)
					})
				})
				Convey("Using an invalid string", func() {
					_, err := Cast("string", reflect.TypeOf(float64(0)))

					Convey("An error should be returned", func() {
						So(err, ShouldNotBeNil)
					})
				})
			})
			Convey("To an invalid type", func() {
				var res = struct{}{}
				_, err := Cast("string", reflect.TypeOf(res))

				Convey("An error should be returned", func() {
					So(err, ShouldNotBeNil)
				})
			})
		})
		Convey("When casting from an int", func() {
			Convey("To a boolean", func() {
				v, err := Cast(1, reflect.TypeOf(true))

				Convey("No error should be returned", func() {
					So(err, ShouldBeNil)
				})
				Convey("The correct value should be returned", func() {
					So(v, ShouldEqual, true)
				})
			})
			Convey("To a string", func() {
				v, err := Cast(123, reflect.TypeOf(""))

				Convey("The correct value should be returned", func() {
					So(v, ShouldEqual, "123")
				})
				Convey("No error should be returned", func() {
					So(err, ShouldBeNil)
				})
			})
			Convey("To an int", func() {
				v, err := Cast(123, reflect.TypeOf(0))

				Convey("No error should be returned", func() {
					So(err, ShouldBeNil)
				})
				Convey("The correct value should be returned", func() {
					So(v, ShouldEqual, 123)
				})
			})
			Convey("To an uint", func() {
				v, err := Cast(123, reflect.TypeOf(uint(0)))

				Convey("No error should be returned", func() {
					So(err, ShouldBeNil)
				})
				Convey("The correct value should be returned", func() {
					So(v, ShouldEqual, uint(123))
				})
			})
			Convey("To a float", func() {
				v, err := Cast(123, reflect.TypeOf(float64(0)))

				Convey("No error should be returned", func() {
					So(err, ShouldBeNil)
				})
				Convey("The correct value should be returned", func() {
					So(v, ShouldEqual, 123)
				})
			})
			Convey("To an invalid type", func() {
				var res = struct{}{}
				_, err := Cast(123, reflect.TypeOf(res))

				Convey("An error should be returned", func() {
					So(err, ShouldNotBeNil)
				})
			})
		})
		Convey("When casting from a uint", func() {
			Convey("To a boolean", func() {
				v, err := Cast(uint(1), reflect.TypeOf(true))

				Convey("No error should be returned", func() {
					So(err, ShouldBeNil)
				})
				Convey("The correct value should be returned", func() {
					So(v, ShouldEqual, true)
				})
			})
			Convey("To a string", func() {
				v, err := Cast(uint(123), reflect.TypeOf(""))

				Convey("No error should be returned", func() {
					So(err, ShouldBeNil)
				})
				Convey("The correct value should be returned", func() {
					So(v, ShouldEqual, "123")
				})
			})
			Convey("To an int", func() {
				v, err := Cast(uint(123), reflect.TypeOf(0))

				Convey("No error should be returned", func() {
					So(err, ShouldBeNil)
				})
				Convey("The correct value should be returned", func() {
					So(v, ShouldEqual, 123)
				})
			})
			Convey("To an uint", func() {
				v, err := Cast(uint(123), reflect.TypeOf(uint(0)))

				Convey("No error should be returned", func() {
					So(err, ShouldBeNil)
				})
				Convey("The correct value should be returned", func() {
					So(v, ShouldEqual, uint(123))
				})
			})
			Convey("To a float", func() {
				v, err := Cast(uint(123), reflect.TypeOf(float64(0)))

				Convey("No error should be returned", func() {
					So(err, ShouldBeNil)
				})
				Convey("The correct value should be returned", func() {
					So(v, ShouldEqual, 123)
				})
			})
			Convey("To an invalid type", func() {
				var res = struct{}{}
				_, err := Cast(uint(123), reflect.TypeOf(res))

				Convey("An error should be returned", func() {
					So(err, ShouldNotBeNil)
				})
			})
		})
		Convey("When casting from a float", func() {
			Convey("To a boolean", func() {
				v, err := Cast(float64(1), reflect.TypeOf(true))

				Convey("No error should be returned", func() {
					So(err, ShouldBeNil)
				})
				Convey("The correct value should be returned", func() {
					So(v, ShouldEqual, true)
				})
			})
			Convey("To a string", func() {
				v, err := Cast(float64(123), reflect.TypeOf(""))

				Convey("No error should be returned", func() {
					So(err, ShouldBeNil)
				})
				Convey("The correct value should be returned", func() {
					So(v, ShouldEqual, "123")
				})
			})
			Convey("To an int", func() {
				v, err := Cast(float64(123), reflect.TypeOf(0))

				Convey("No error should be returned", func() {
					So(err, ShouldBeNil)
				})
				Convey("The correct value should be returned", func() {
					So(v, ShouldEqual, 123)
				})
			})
			Convey("To an uint", func() {
				v, err := Cast(float64(123), reflect.TypeOf(uint(0)))

				Convey("No error should be returned", func() {
					So(err, ShouldBeNil)
				})
				Convey("The correct value should be returned", func() {
					So(v, ShouldEqual, uint(123))
				})
			})
			Convey("To a float", func() {
				v, err := Cast(float64(123), reflect.TypeOf(float64(0)))

				Convey("No error should be returned", func() {
					So(err, ShouldBeNil)
				})
				Convey("The correct value should be returned", func() {
					So(v, ShouldEqual, 123)
				})
			})
			Convey("To an invalid type", func() {
				var res = struct{}{}
				_, err := Cast(float64(123), reflect.TypeOf(res))

				Convey("An error should be returned", func() {
					So(err, ShouldNotBeNil)
				})
			})
		})
		Convey("When casting from an unsupported type", func() {
			_, err := Cast(struct{}{}, reflect.TypeOf(0))

			Convey("An error should be returned", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}
