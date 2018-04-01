[![Build Status](https://travis-ci.org/rbarilani/eskip-match.svg?branch=master)](https://travis-ci.org/rbarilani/eskip-match) [![codecov](https://codecov.io/gh/rbarilani/eskip-match/branch/master/graph/badge.svg)](https://codecov.io/gh/rbarilani/eskip-match) [![Go Report Card](https://goreportcard.com/badge/github.com/rbarilani/eskip-match)](https://goreportcard.com/report/github.com/rbarilani/eskip-match) [![GoDoc](https://godoc.org/github.com/rbarilani/eskip-match?status.svg)](https://godoc.org/github.com/rbarilani/eskip-match) [![git tag](https://img.shields.io/github/tag/rbarilani/eskip-match.svg)](https://img.shields.io/github/tag/rbarilani/eskip-match.svg)

# eskip-match 

A package that helps you test [skipper](https://github.com/zalando/skipper) [`.eskip`](https://zalando.github.io/skipper/dataclients/eskip-file/) files routing matching logic.



* [Install](#install)
* [Usage](#usage)
* [CLI](#cli)
* [License](#license)

## Install

```
go get github.com/rbarilani/eskip-match/...
```

## Usage

Given an `.eskip` file:

*routes.eskip*
```
foo: Path("/foo") -> http://foo.com
foo_post: Path("/foo") && Method("POST") -> http://foopost.com
bar: Path("/bar") -> http://bar.com
```

You can write a `go test` able to check if the matching logic is what you expect.
A simple example: 

*main_test.go*
```go
package main

import (
	"path/filepath"
	"testing"

	"github.com/rbarilani/eskip-match/matcher"
)

func TestReadmeExample(t *testing.T) {
	m, err := createMatcher()

	if err != nil {
		t.Fatal(err)
		return
	}

	res := m.Test(&matcher.RequestAttributes{
		Path:   "/foo",
		Method: "POST",
	})

	route := res.Route()

	if route == nil {
		t.Error("Expect matching but no match")
		return
	}
	if route.Id != "foo_post" {
		t.Errorf("Expect matching route: %s but got %s", "foo_post", route.Id)
	}
}

func createMatcher() (matcher.Matcher, error) {
	routesFile, err := filepath.Abs("./routes.eskip")

	if err != nil {
		return nil, err
	}

	m, err := matcher.New(&matcher.Options{
		RoutesFile: routesFile,
	})

	if err != nil {
		return nil, err
	}

	return m, nil
}


```

## CLI

The package provide a binary cli tool: `eskip-match`

```
NAME:
   eskip-match - A command line tool that helps you test .eskip files routing matching logic

USAGE:
   eskip-match [global options] command [command options] [arguments...]

COMMANDS:
     test, t  Given a routes file and request attributes, checks a route matches
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config FILE, -c FILE  Load configuration from FILE
   --help, -h              show help
   --version, -v           print the version

```

### Commands

### Test

With `eskip-match test` command you can check if a route matches given specific request attributes.

#### Examples

Test if a request to path `/foo` matches a route:

```bash
eskip-match test routes.eskip -p /foo
```

Test if a request to path `/foo` using `GET` method matches a route:

```bash
eskip-match test routes.eskip -p /foo -m GET
```

Specifying headers:

```bash
eskip-match test routes.eskip -p /foo -H Accept=application/json -H Authorization="Bearer XXX"
```

Using **verbose output** might help when something doesn't seem to work as expected:

```bash
eskip-match test routes.eskip -v -p /foo
```

If your routes are using **custom filters** the tool must be informed via a **configuration file** named `.eskip-match.yml`, eg:

*.eskip-match.yml*
```yaml
customfilters:
  - myCustomFilter1
  - myCustomFilter2
```

```bash
eskip-match test routes.eskip -p /foo
```

> By default the tool will try to load `.eskip-match.yml` in the current working directory, but you can provide a custom location with `-c` global option, eg:
```bash
eskip-match -c config.yml test routes.eskip -p /foo
```



## License

Copyright 2018 Ruben Barilani

MIT License

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
