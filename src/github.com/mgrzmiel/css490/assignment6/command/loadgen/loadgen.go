// CSS 490
// Magdalena Grzmiel
// Assignments #6
// Copyright 2015 Magdalena Grzmiel
// loadgen is responsible for genrating requests to the server

package main

import (
	"flag"
	//"fmt"
	"github.com/mgrzmiel/css490/assignment6/lib/counter"
	"net/http"
	"time"
)

var rate int
var burst int
var timeoutMS int
var url string
var runtime int
var runtimeInMS time.Duration

// create counter
var (
	c = counter.New()
)

var convert = map[int]string{
	1: "100s",
	2: "200s",
	3: "300s",
	4: "400s",
	5: "500s",
}

// request
// It is responsole for sending one request to the server
func request() {
	// increse the number of total request
	c.Incr("total", 1)

	// set the max time to wait for response in milliseconds
	client := http.Client{
		Timeout: time.Duration(timeoutMS) * time.Millisecond,
	}

	// send the request
	response, err := client.Get(url)
	if err != nil {
		// if error increse the number of total errors
		c.Incr("errors", 1)
		return
	}

	// get the response status code and convert it
	key, ok := convert[response.StatusCode/100]
	if !ok {
		key = "errors"
	}

	// increase the number
	c.Incr(key, 1)
}

// load
// load is responsoble for generating the required number of requests
func load() {
	timeout := time.Tick(runtimeInMS)
	interval := time.Duration((1000000*burst)/rate) * time.Microsecond
	period := time.Tick(interval)
	for {
		//fire off burst
		for i := 0; i < burst; i++ {
			go request()
		}
		// wait for next tick
		<-period

		// if the tiemout already passed, return
		// otherwise continue
		select {
		case <-timeout:
			return
		default:
		}
	}

}

func main() {
	flag.IntVar(&rate, "rate", 200, "rate")
	flag.IntVar(&burst, "burst", 20, "burst")
	flag.IntVar(&timeoutMS, "timeout-ms", 400, "timeoutMS")
	flag.StringVar(&url, "url", "http://localhost:8080/time", "url")
	flag.IntVar(&runtime, "runtime", 20, "runtime")

	flag.Parse()

	// set the runtime to time.Duration type
	runtimeInMS = time.Duration(runtime) * time.Second

	load()

	// sleep to make sure all of the request are served
	time.Sleep(time.Duration(2*timeoutMS) * time.Millisecond)

	// print results
	// fmt.Printf("total: \t%d\n", c.Get("total"))
	// fmt.Printf("100s: \t%d\n", c.Get("100s"))
	// fmt.Printf("200s: \t%d\n", c.Get("200s"))
	// fmt.Printf("300s: \t%d\n", c.Get("300s"))
	// fmt.Printf("400s: \t%d\n", c.Get("400s"))
	// fmt.Printf("500s: \t%d\n", c.Get("500s"))
	// fmt.Printf("errors: \t%d\n", c.Get("errors"))
}
