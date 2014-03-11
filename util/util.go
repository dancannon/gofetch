package util

import (
	"bytes"
	"code.google.com/p/go.net/html"
	"github.com/PuerkitoBio/goquery"
	htmlutil "html"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func SelectionToString(selection *goquery.Selection) string {
	var buf bytes.Buffer

	// Slightly optimized vs calling Each(): no single selection object created
	for i, n := range selection.Nodes {
		buf.WriteString(GetNodeText(n))

		if i < len(selection.Nodes)-1 {
			buf.WriteString("\n")
		}
	}
	return buf.String()
}

func GetNodeText(node *html.Node) string {
	if node.Type == html.TextNode {
		var re *regexp.Regexp

		t := node.Data
		re = regexp.MustCompile("[\t\r\n]+")
		t = re.ReplaceAllString(t, "")
		re = regexp.MustCompile(" {2,}")
		t = re.ReplaceAllString(t, " ")
		t = htmlutil.EscapeString(t)

		return t
	} else if node.FirstChild != nil {
		var buf bytes.Buffer
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			buf.WriteString(GetNodeText(c))
			if c.Type == html.ElementNode {
				switch c.Data {
				case "strike", "u", "b", "i", "em", "strong", "span", "sup",
					"code", "tt", "sub", "var", "font":
				default:
					buf.WriteString("\n")
				}
			}
		}
		return buf.String()
	}

	return ""
}

func CreateMapFromProps(props map[string]interface{}, keys map[string]string) map[string]interface{} {
	m := make(map[string]interface{})

	for mapKey, propKey := range keys {
		if val, ok := props[propKey]; ok {
			kp := strings.Split(mapKey, ":")
			// If the map key consists of more that 1 part then cast the value
			// to the type of the second value. For example width:int should
			// cast the value of width to an int.
			if len(kp) > 1 {
				mapKey = kp[0]

				switch kp[1] {
				case "int":
					v, err := Cast(val, reflect.TypeOf(int(0)))
					if err == nil {
						val = v
					}
				case "uint":
					v, err := Cast(val, reflect.TypeOf(uint(0)))
					if err == nil {
						val = v
					}
				case "string":
					v, err := Cast(val, reflect.TypeOf(string("")))
					if err == nil {
						val = v
					}
				case "float":
					v, err := Cast(val, reflect.TypeOf(float64(0)))
					if err == nil {
						val = v
					}
				}
			}

			m[mapKey] = val
		}
	}

	return m
}

// An CastTypeError describes a value that was
// not appropriate for a value of a specific Go type.
type CastTypeError struct {
	Value string       // description of value - "bool", "array", "number -5"
	Type  reflect.Type // type of Go value it could not be assigned to
}

func (e *CastTypeError) Error() string {
	return "cannot decode " + e.Value + " into Go value of type " + e.Type.String()
}

// Cast a value to another type.
func Cast(sv interface{}, typ reflect.Type) (interface{}, error) {
	rv, err := RCast(reflect.ValueOf(sv), typ)
	if err != nil {
		return nil, err
	}
	return rv.Interface(), nil
}

// Cast a value to another type using reflection.
func RCast(sv reflect.Value, typ reflect.Type) (reflect.Value, error) {
	dv := reflect.New(typ)
	dv = indirect(dv, true)

	// Special case for if sv is nil:
	if sv.Kind() == reflect.Invalid {
		return dv, nil
	}

	// Attempt to convert the value from the source type to the destination type
	switch sv.Kind() {
	case reflect.Bool:
		switch dv.Kind() {
		default:
			return reflect.Value{}, &CastTypeError{"bool", dv.Type()}
		case reflect.Bool:
			dv.Set(sv)
		case reflect.String:
			dv.SetString(strconv.FormatBool(sv.Bool()))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if sv.Bool() {
				dv.SetInt(1)
			} else {
				dv.SetInt(0)
			}
		}

	case reflect.String:
		switch dv.Kind() {
		default:
			return reflect.Value{}, &CastTypeError{sv.Type().String(), dv.Type()}
		case reflect.String:
			dv.Set(sv)
		case reflect.Bool:
			b, err := strconv.ParseBool(sv.String())
			if err != nil {
				return reflect.Value{}, &CastTypeError{sv.Type().String(), dv.Type()}
			}
			dv.SetBool(b)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			n, err := strconv.ParseInt(sv.String(), 10, 64)
			if err != nil || dv.OverflowInt(n) {
				return reflect.Value{}, &CastTypeError{sv.Type().String(), dv.Type()}
			}
			dv.SetInt(n)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			n, err := strconv.ParseUint(sv.String(), 10, 64)
			if err != nil || dv.OverflowUint(n) {
				return reflect.Value{}, &CastTypeError{sv.Type().String(), dv.Type()}
			}
			dv.SetUint(n)
		case reflect.Float32, reflect.Float64:
			n, err := strconv.ParseFloat(sv.String(), 64)
			if err != nil || dv.OverflowFloat(n) {
				return reflect.Value{}, &CastTypeError{sv.Type().String(), dv.Type()}
			}
			dv.SetFloat(n)
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch dv.Kind() {
		default:
			return reflect.Value{}, &CastTypeError{sv.Type().String(), dv.Type()}
		case reflect.Bool:
			dv.SetBool(sv.Int() == 1)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			dv.SetInt(int64(sv.Int()))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			dv.SetUint(uint64(sv.Int()))
		case reflect.Float32, reflect.Float64:
			dv.SetFloat(float64(sv.Int()))
		case reflect.String:
			dv.SetString(strconv.FormatInt(int64(sv.Int()), 10))
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		switch dv.Kind() {
		default:
			return reflect.Value{}, &CastTypeError{sv.Type().String(), dv.Type()}
		case reflect.Bool:
			dv.SetBool(sv.Uint() == 1)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			dv.SetInt(int64(sv.Uint()))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			dv.SetUint(uint64(sv.Uint()))
		case reflect.Float32, reflect.Float64:
			dv.SetFloat(float64(sv.Uint()))
		case reflect.String:
			dv.SetString(strconv.FormatUint(uint64(sv.Uint()), 10))
		}
	case reflect.Float32, reflect.Float64:
		switch dv.Kind() {
		default:
			return reflect.Value{}, &CastTypeError{sv.Type().String(), dv.Type()}
		case reflect.Bool:
			dv.SetBool(sv.Float() == 1)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			dv.SetInt(int64(sv.Float()))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			dv.SetUint(uint64(sv.Float()))
		case reflect.Float32, reflect.Float64:
			dv.SetFloat(float64(sv.Float()))
		case reflect.String:
			dv.SetString(strconv.FormatFloat(float64(sv.Float()), 'g', -1, 64))
		}
	default:
		return reflect.Value{}, &CastTypeError{sv.Type().String(), dv.Type()}
	}

	return dv, nil
}

// indirect walks down v allocating pointers as needed,
// until it gets to a non-pointer.
func indirect(v reflect.Value, decodeNull bool) reflect.Value {
	// If v is a named type and is addressable,
	// start with its address, so that if the type has pointer methods,
	// we find them.
	if v.Kind() != reflect.Ptr && v.Type().Name() != "" && v.CanAddr() {
		v = v.Addr()
	}
	for {
		// Load value from interface, but only if the result will be
		// usefully addressable.
		if v.Kind() == reflect.Interface && !v.IsNil() {
			e := v.Elem()
			if e.Kind() == reflect.Ptr && !e.IsNil() && (!decodeNull || e.Elem().Kind() == reflect.Ptr) {
				v = e
				continue
			}
		}

		if v.Kind() != reflect.Ptr {
			break
		}

		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		v = v.Elem()
	}
	return v
}
