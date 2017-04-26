# Logzio golang client
Send logs to Logzio

[![Build Status](https://travis-ci.org/dougEfresh/logzio-go.svg?branch=master)](https://travis-ci.org/dougEfresh/logzio-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/dougEfresh/logzio-go)](https://goreportcard.com/report/github.com/dougEfresh/logzio-go)
[![GoDoc](https://godoc.org/github.com/dougEfresh/logzio-go?status.svg)](https://godoc.org/github.com/dougEfresh/logzio-go)
[![license](http://img.shields.io/badge/license-apache-red.svg?style=flat)](https://raw.githubusercontent.com/dougEfresh/logzio-go/master/LICENSE)

## Getting Started

### Get Logzio token
1. Go to Logzio website
2. Sign in with your Logzio account
3. Click the top menu gear icon (Account)
4. The Logzio token is given in the account page

### Initialize Logger
```go
package main

import (
	"fmt"
	"github.com/dougEfresh/logzio-go"
	"os"
	"time"
)

func main() {
	l, err := logzio.New(os.Args[1]) // Token is required
	if err != nil {
		panic(err)
	}
	msg := fmt.Sprintf("{ \"%s\": \"%s\"}", "message", time.Now().UnixNano())
	err = l.Send([]byte(msg))
	if err != nil {
		panic(err)
	}
	l.Stop() //logs are buffered on disk. Stop will drain the buffer
}

```
