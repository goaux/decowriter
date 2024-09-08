# decowriter

Package decowriter provides a writer that add a prefix and a suffix to each line.

[![Go Reference](https://pkg.go.dev/badge/github.com/goaux/decowriter.svg)](https://pkg.go.dev/github.com/goaux/decowriter)
[![Go Report Card](https://goreportcard.com/badge/github.com/goaux/decowriter)](https://goreportcard.com/report/github.com/goaux/decowriter)

A line is defines as a sequence of zero or more non-'\n' bytes followed by a '\n'.

A prefix is written before one or more bytes of a line.

A suffix is written just before the trailing '\n'.

A string without a trailing '\n' is not recognized as a line, so no suffix is written.
This allows you to split a single line writing across two Write calls.

## Features

- Adds a specified prefix to the beginning of each line
- Adds a specified suffix to the end of each line
- Implements the `io.Writer` interface
- Tracks total bytes written

## Usage

```go
package main

import (
	"fmt"
	"os"

	"github.com/goaux/decowriter"
)

func main() {
	w := decowriter.New(bufio.NewWriter(os.Stdout), []byte(">>"), []byte("<<"))
	fmt.Fprintln(w, "This is a log message")
	fmt.Fprintln(w, "Another log message")
}
```

OUTPUT:

```
>>This is a log message<<
>>Another log message<<
```
