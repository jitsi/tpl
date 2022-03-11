package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

// template shared context
var (
	ctx interface{}
)

// create template context
func newTemplateVariables() map[string]interface{} {
	var vars = make(map[string]interface{})

	envs := make(map[string]interface{})
	for _, env := range os.Environ() {
		kv := strings.SplitN(env, "=", 2)
		envs[kv[0]] = kv[1] // .Env.name
	}
	vars["Env"] = envs

	return vars
}

func templateExecute(t *template.Template, srcFile string) {
	var err error
	var templateBytes []byte

	templateBytes, err = ioutil.ReadFile(srcFile)
	if err != nil {
		panic(fmt.Errorf("unable to read from %v: %v", srcFile, err))
	}

	tmpl, err := t.Parse(string(templateBytes))
	if err != nil {
		panic(fmt.Errorf("unable to parse template file: %v", err))
	}

	var output bytes.Buffer
	err = tmpl.Execute(&output, ctx)
	if err != nil {
		panic(fmt.Errorf("render template error: %v", err))
	}

	_, err = os.Stdout.Write(output.Bytes())
	if err != nil {
		panic(fmt.Errorf("error writing template: %v", err))
	}
}

func main() {
	argLength := len(os.Args[1:])
	if argLength != 1 {
		fmt.Print("Missing source file")
		os.Exit(1)
	}

	srcFile := os.Args[1]
	ctx = newTemplateVariables()
	t := template.New(srcFile)
	t.Option("missingkey=error")
	t.Funcs(FuncMap())
	templateExecute(t, srcFile)
}
