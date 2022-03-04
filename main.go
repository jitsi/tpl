package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/subchen/go-cli"
)

// version
var (
	BuildVersion   string
	BuildGitBranch string
	BuildGitRev    string
	BuildGitCommit string
	BuildDate      string
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
	app := cli.NewApp()
	app.Name = "tol"
	app.Usage = "Generate file using template"
	app.UsageText = "[options] input-file[:output-file] ..."
	app.Authors = "Guoqiang Chen <subchen@gmail.com>"

	app.Examples = strings.TrimSpace(`
frep nginx.conf.in -e webroot=/usr/share/nginx/html -e port=8080
frep nginx.conf.in:/etc/nginx.conf -e webroot=/usr/share/nginx/html -e port=8080
frep nginx.conf.in --json '{"webroot": "/usr/share/nginx/html", "port": 8080}'
frep nginx.conf.in --load config.json --overwrite
echo "{{ .Env.PATH }}"  | frep -
`)

	app.Version = BuildVersion
	app.BuildInfo = &cli.BuildInfo{
		GitBranch:   BuildGitBranch,
		GitCommit:   BuildGitCommit,
		GitRevCount: BuildGitRev,
		Timestamp:   BuildDate,
	}

	app.Action = func(c *cli.Context) {
		if c.NArg() == 0 {
			c.ShowHelp()
			return
		}

		defer func() {
			if err := recover(); err != nil {
				os.Stderr.WriteString(fmt.Sprintf("fatal: %v\n", err))
				os.Exit(1)
			}
		}()

		ctx = newTemplateVariables()
		for _, file := range c.Args() {
			filePair := strings.SplitN(file, ":", 2)
			srcFile := filePair[0]

			t := template.New(srcFile)
			t.Option("missingkey=error")
			t.Funcs(FuncMap(file))

			templateExecute(t, file)
		}
	}

	app.Run(os.Args)
}
