package main

import (
	"bytes"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Person struct {
	Name    string `properties:"name"`
	Address string `properties:"address,omitempty"`
	Age     int    `properties:"age"`
	Married bool   `properties:"married"`
}

func Serialize(person any) string {
	v := reflect.ValueOf(person)

	var b bytes.Buffer

	for i := 0; i < v.NumField(); i++ {
		tags := getTags(v, i)
		if len(tags) == 0 {
			continue
		}
		if isOmitempty(tags, v, i) {
			continue
		}
		b.WriteString(tags[0] + "=" + getValue(v.Field(i)))
		if i < v.NumField()-1 {
			b.WriteString("\n")
		}
	}

	return b.String()
}

func getTags(v reflect.Value, i int) []string {
	tags := v.Type().Field(i).Tag.Get("properties")
	parts := strings.Split(tags, ",")
	return parts
}

func isOmitempty(parts []string, v reflect.Value, i int) bool {
	return len(parts) > 1 && parts[1] == "omitempty" && IsEmpty(v.Field(i))
}

func IsEmpty(rv reflect.Value) bool {
	switch rv.Kind() {
	case reflect.String, reflect.Array:
		return rv.Len() == 0
	case reflect.Map, reflect.Slice:
		return rv.IsNil() || rv.Len() == 0
	case reflect.Bool:
		return !rv.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return rv.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return rv.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return rv.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return rv.IsNil()
	case reflect.Struct:
		return IsEmpty(rv.Elem())
	default:
		return true
	}
}

func getValue(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.Itoa(int(v.Int()))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', -1, len(v.Bytes())*8)
	case reflect.String:
		return v.String()
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.Interface, reflect.Ptr:
		return v.Interface().(string)
	case reflect.Array, reflect.Slice:
		var str strings.Builder
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				str.WriteString(",")
			}
			str.WriteString(getValue(v.Index(i)))
		}
		return str.String()
	case reflect.Map:
		var str strings.Builder
		for _, key := range v.MapKeys() {
			str.WriteString("{")
			str.WriteString(getValue(v.MapIndex(key)))
			str.WriteString("}")
		}
		return str.String()
	case reflect.Struct:
		var str strings.Builder
		str.WriteString("{")
		for i := 0; i < v.NumField(); i++ {
			if i > 0 {
				str.WriteString(",")
			}
			str.WriteString(getValue(v.Field(i)))
		}
		str.WriteString("}")
		return str.String()
	default:
		return v.String()
	}
}

func TestSerialization(t *testing.T) {
	tests := map[string]struct {
		person Person
		result string
	}{
		"test case with empty fields": {
			result: "name=\nage=0\nmarried=false",
		},
		"test case with fields": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
			},
			result: "name=John Doe\nage=30\nmarried=true",
		},
		"test case with omitempty field": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
				Address: "Paris",
			},
			result: "name=John Doe\naddress=Paris\nage=30\nmarried=true",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := Serialize(test.person)
			assert.Equal(t, test.result, result)
		})
	}
}
