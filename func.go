package main

import (
	"fmt"
	"strconv"
	"strings"
	"text/template"

	"github.com/go-sprout/sprout"
)

func FuncMap() template.FuncMap {
	// sprout
	f := sprout.TxtFuncMap()

	// marshaling
	f["toBool"] = toBool

	// strings
	f["countRune"] = func(s string) int {
		return len([]rune(s))
	}

	// Fix sprout regex functions
	oRegexReplaceAll := f["regexReplaceAll"].(func(regex string, s string, repl string) string)
	oRegexReplaceAllLiteral := f["regexReplaceAllLiteral"].(func(regex string, s string, repl string) string)
	oRegexSplit := f["regexSplit"].(func(regex string, s string, n int) []string)
	f["reReplaceAll"] = func(regex string, replacement string, input string) string {
		return oRegexReplaceAll(regex, input, replacement)
	}
	f["reReplaceAllLiteral"] = func(regex string, replacement string, input string) string {
		return oRegexReplaceAllLiteral(regex, input, replacement)
	}
	f["reSplit"] = func(regex string, n int, input string) []string {
		return oRegexSplit(regex, input, n)
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
