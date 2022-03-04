package main

import (
	"strconv"
	"text/template"

	sprig "github.com/Masterminds/sprig/v3"
)

func FuncMap(templateName string) template.FuncMap {
	f := sprig.TxtFuncMap()
	// marshaling
	f["toBool"] = toBool

	// strings
	f["countRune"] = func(s string) int {
		return len([]rune(s))
	}

	// Fix sprig regex functions
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
func toBool(value string) bool {
	result, err := strconv.ParseBool(value)
	if err != nil {
		if Strict {
			panic(err.Error())
		}
		return false
	}
	return result
}
