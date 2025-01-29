package main

import (
	"fmt"
	"strconv"
	"strings"
	"text/template"

	"github.com/go-sprout/sprout/sprigin"
)

func FuncMap() template.FuncMap {
	// sprout
	f := sprigin.TxtFuncMap()

	// marshaling
	f["toBool"] = toBool

	// strings
	f["countRune"] = func(s string) int {
		return len([]rune(s))
	}

	// Fix sprout regex functions
	oRegexReplaceAll := f["regexReplaceAll"].(func(...interface{}) (interface{}, error))
	oRegexReplaceAllLiteral := f["regexReplaceAllLiteral"].(func(...interface{}) (interface{}, error))
	oRegexSplit := f["regexSplit"].(func(...interface{}) (interface{}, error))

	f["reReplaceAll"] = func(regex string, replacement string, input string) string {
		out, err := oRegexReplaceAll(regex, input, replacement)
		if err != nil {
			return input // Return original string on error
		}
		return out.(string)
	}

	f["reReplaceAllLiteral"] = func(regex string, replacement string, input string) string {
		out, err := oRegexReplaceAllLiteral(regex, input, replacement)
		if err != nil {
			return input // Return original string on error
		}
		return out.(string)
	}

	f["reSplit"] = func(regex string, n int, input string) []string {
		out, err := oRegexSplit(regex, input, n)
		if err != nil {
			return []string{input} // Return single-element slice on error
		}
		return out.([]string)
	}

	return f
}

// toBool takes a string and converts it to a bool.
// On marshal error will panic if in strict mode, otherwise returns false.
// It accepts 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False.
//
// This is designed to be called from a template.
func toBool(value interface{}) bool {
	switch val := value.(type) {
	case int:
		return val != 0
	case float32:
	case float64:
		return val != 0.0
	case string:
		v := strings.Replace(val, "\"", "", -1)
		v = strings.ReplaceAll(v, "'", "")
		result, err := strconv.ParseBool(v)
		if err != nil {
			return len(v) > 0
		}
		return result
	case bool:
		return val
	case nil:
		return false
	default:
		panic(fmt.Sprintf("unsupported value type %s", val))
	}

	return false // appease the linter
}
