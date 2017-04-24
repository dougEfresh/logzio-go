# Logzio for go-kit logger
Send go-kit logs to Logzio

[![Build Status](https://travis-ci.org/dougEfresh/kitz.svg?branch=master)](https://travis-ci.org/dougEfresh/kitz)
[![Go Report Card](https://goreportcard.com/badge/github.com/dougEfresh/kitz)](https://goreportcard.com/report/github.com/dougEfresh/kitz)
[![GoDoc](https://godoc.org/github.com/dougEfresh/kitz?status.svg)](https://godoc.org/github.com/dougEfresh/kitz)
[![license](http://img.shields.io/badge/license-apache-red.svg?style=flat)](https://raw.githubusercontent.com/dougEfresh/kitz/master/LICENSE)

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
        "github.com/dougEfresh/kitz"
        "github.com/go-kit/kit/log"
)

func main() {
        klogger, err := kitz.New("123456789")
        if err != nil {
                panic(err)
        }
        // returns the go-kit logger
        logger := klogger.Build()
        // message is required
        logger.Log("message", "hello!")
}
```
