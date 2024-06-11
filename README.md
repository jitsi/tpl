# tpl

Simple templating from environment variables.

## Usage

Given the following sample file (sample.txt):

```
Hello {{ .Env.USER }}
```

```bash
$ tpl sample.txt
Hello saghul
```

The output is always stdout. This utility is means to be used as follows:

```bash
$ tpl nginx.conf.tpl > nginx.conf
```

## Template context

Templates use Golang [text/template](http://golang.org/pkg/text/template/).

You can access environment variables within a template in the `.Env` object.

There are some built-in functions as well: sprout
- github: https://github.com/go-sprout/sprout
- doc: https://docs.atom.codes/sprout

More functions:
- toBool
- countRune
- pipeline compatible regex functions from sprout
    - reReplaceAll
    - reReplaceAllLiteral
    - reSplit

## Thanks

This project is a fork of [frep](https://github.com/subchen/frep) with a more
limited scope. Thank you Guoqiang Chen for creating frep!
