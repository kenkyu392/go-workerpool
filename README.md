# go-workerpool

[![go.dev reference](https://img.shields.io/badge/go.dev-reference-00ADD8?logo=go)](https://pkg.go.dev/github.com/kenkyu392/go-workerpool)
[![go report card](https://goreportcard.com/badge/github.com/kenkyu392/go-workerpool)](https://goreportcard.com/report/github.com/kenkyu392/go-workerpool)
[![license](https://img.shields.io/github/license/kenkyu392/go-workerpool.svg)](LICENSE)

Simple and versatile Go worker pool.

## Installation

```
go get -u github.com/kenkyu392/go-workerpool
```

## Usage

```go
package main

import (
	"log"
	"math/rand"
	"runtime"
	"time"

	"github.com/kenkyu392/go-workerpool"
)

func init() {
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	animals := []string{"Lion", "Tiger", "Cheetah", "Jaguar", "Leopard", "Cat", "Cougar"}

	// Create a worker pool by specifying the maximum concurrent number and
	// start goroutines.
	wp := workerpool.New(2)

	// If you try to add more jobs than the number of concurrent executions,
	// processing will be blocked.
	for n, v := range animals {
		n := n
		v := v
		wp.AddJobFunc(func() error {
			d := time.Millisecond * time.Duration(rand.Intn(400)+400)
			log.Printf("number:%d goroutines:%d sleep:%v animal:%s", n, runtime.NumGoroutine(), d, v)
			time.Sleep(d)
			return nil
		})
	}

	for _, err := range wp.Wait() {
		log.Println(err)
	}
}
```

## License

[MIT](LICENSE)
