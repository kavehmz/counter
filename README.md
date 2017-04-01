# counter
[![GoDoc](https://godoc.org/github.com/kavehmz/counter?status.svg)](https://godoc.org/github.com/kavehmz/counter)
[![Build Status](https://travis-ci.org/kavehmz/counter.svg?branch=master)](https://travis-ci.org/kavehmz/counter)
[![Coverage Status](https://coveralls.io/repos/kavehmz/counter/badge.svg?branch=master&service=github)](https://coveralls.io/github/kavehmz/counter?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/kavehmz/counter)](https://goreportcard.com/report/github.com/kavehmz/counter)

This is a [Go](http://golang.org) library to increment a counter for a window of time (say last 60 second). Its resolution is 1 second.

Results will be periodically saved in a local file.

## Installation

```bash
$ go get github.com/kavehmz/counter
```

# Usage
```bash
# to test it as a http server you can do:
$ go run example/main.go
```

```go
package main

import (
	"fmt"
	"github.com/kavehmz/counter"
)

func main() {
	c := counter.Init("/tmp/counter", time.Second, 60)
	ch := make(chan int)
	c.Inc(ch)
    fmt.Println(<-ch)
}
```

# Algorithm
This lib is using and travesing a linked list  to keep track of counts.
