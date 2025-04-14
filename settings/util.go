package settings

import (
	"reflect"
	"strings"
	"unicode"
)

// get field from settings name / fieldMapping key
// from an alias that is case insensitive
func getMappingKeyFromAlias(alias string) (o string) {
	alias = strings.ToLower(alias)
	switch alias {
	case "amp-scalar":
		o = "AmpScalar"
	case "deviceidx", "device-idx":
		o = "DeviceIdx"
	case "-d":
		o = "--device-idx"
	case "port":
		o = "--device-idx"
	}

	return
}

func getFuncType(_key string) (
	funcType interface{},
	ok bool,
	foundKey string,
) {
	var f func(string, bool) string

	f = func(key string, isAlias bool) string {
		funcType, ok = fieldMapping[key]

		// if we didn't find the key
		// try to lookup with alias
		if !ok && !isAlias {
			// fmt.Println("alias", key)
			return f(getMappingKeyFromAlias(key), true)
		}
		return key
	}

	foundKey = f(_key, false)

	return
}

func kebabToPascal(s string) string {
	// remove all suffix/prefix -
	s = strings.Trim(s, "-")

	var (
		hadDash = false
		sb      strings.Builder
	)
	sb.WriteRune(unicode.ToUpper(rune(s[0])))
	for _, c := range s[1:] {
		if c == '-' {
			hadDash = true
			continue
		}

		if hadDash {
			sb.WriteRune(unicode.ToUpper(c))
			hadDash = false
			continue
		}

		sb.WriteRune(c)
	}

	return sb.String()
}


func getFieldPointer(s *Settings, field string) any {
	structVal := reflect.ValueOf(s)
	if !structVal.IsValid() {
		panic("val of struct not valid")
	}

	structVal = structVal.Elem()
	if !structVal.IsValid() {
		panic("struct val elem not valid")
	}

	f := structVal.FieldByName(field)
	if !f.IsValid() {
		panic("field not valid")
	}
	if !f.Addr().IsValid() {
		panic("field addr not valid")
	}

	return f.Addr().Interface()
}
