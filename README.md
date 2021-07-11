# Logzio Golang API client

Sends logs to [logz.io](https://logz.io) over HTTP. It is a low level lib that can to be integrated with other logging libs.

[![GoDoc][doc-img]][doc] [![Build Status][ci-img]][ci] [![Coverage Status][cov-img]][cov] [![Go Report][report-img]][report]

## Prerequisites
go 1.x

## Installation
```shell
$ go get -u github.com/logzio/logzio-go
```

## Quick Start

### Disk queue
```go
package main

import (
  "fmt"
  "github.com/logzio/logzio-go"
  "os"
  "time"
)

func main() {
  l, err := logzio.New(
  		"fake-token",
  		SetDebug(os.Stderr),
  		SetUrl("http://localhost:12345"),
  		SetDrainDuration(time.Minute*10),
        SetSetTempDirectory("myQueue"),
        SetDrainDiskThreshold(99)
  	) // token is required
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

### In memory queue
```go
package main

import (
  "fmt"
  "github.com/logzio/logzio-go"
  "os"
  "time"
)

func main() {
  l, err := logzio.New(
  		"fake-token",
  		SetDebug(os.Stderr),
  		SetUrl("http://localhost:12345"),
	    SetInMemoryQueue(true),
	    SetinMemoryCapacity(24000000),
	    SetlogCountLimit(6000000),
  	) // token is required
  if err != nil {
    panic(err)
  }
  msg := fmt.Sprintf("{ \"%s\": \"%s\"}", "message", time.Now().UnixNano())

  err = l.Send([]byte(msg))
  if err != nil {
     panic(err)
  }

  l.Stop() 
}
```

## Usage

- Set url mode:
    `logzio.New(token, SetUrl(ts.URL))`

- Set drain duration (flush logs on disk):
    `logzio.New(token, SetDrainDuration(time.Hour))`

- Set debug mode:
    `logzio.New(token, SetDebug(os.Stderr))`

- Set queue dir:
    `logzio.New(token, SetSetTempDirectory(os.Stderr))`

- Set the sender to check if it crosses the maximum allowed disk usage:
    `logzio.New(token, SetCheckDiskSpace(true))`

- Set disk queue threshold, once the threshold is crossed the sender will not enqueue the received logs:
    `logzio.New(token, SetDrainDiskThreshold(99))`

- Set the sender to Use in memory queue:
  `logzio.New(token, SetInMemoryQueue(true))`

- Set the sender to Use in memory queue with log count limit and capacity:
  `logzio.New(token,
  SetInMemoryQueue(true),
  SetinMemoryCapacity(500),
  SetlogCountLimit(6000000),
  )`

## Disk queue
Logzio go client uses [goleveldb](https://github.com/syndtr/goleveldb) and [goqueue](github.com/beeker1121/goque) as a persistent storage.
Every 5 seconds logs are sent to logz.io (if any are available)

## In memory queue
You can see the logzio go client queue implementation in `inMemoryQueue.go` file

## Tests

```shell
$ go test -v

```


See [travis.yaml](.travis.yml) for running benchmark tests


## Contributing
 All PRs are welcome

## Authors

* **Douglas Chimento**  - [dougEfresh][me]
* **Ido Halevi**  - [idohalevi](https://github.com/idohalevi)


## License

This project is licensed under the Apache License - see the [LICENSE](LICENSE) file for details

## Acknowledgments

* [logzio-java-sender](https://github.com/logzio/logzio-java-sender)
