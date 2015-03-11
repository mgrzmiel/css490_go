// CSS 490
// Magdalena Grzmiel
// Assignments #6
// Copyright 2015 Magdalena Grzmiel
// loadgen is responsible for genrating requests to the server

package main

import (
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/mgrzmiel/css490/assignment6/lib/counter"
	"github.com/mgrzmiel/css490/assignment6/lib/loadgenConfig"
	"net/http"
	"time"
)

var runtimeDuration time.Duration

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
		Timeout: time.Duration(loadgenConfig.TimeoutMS) * time.Millisecond,
	}

	// send the request
	response, err := client.Get(loadgenConfig.Url)
	if err != nil {
		// if error increse the number of total errors
		c.Incr("errors", 1)
		log.Errorf("Error while getting response from client: %s", err)
		return
	}

	// get the response status code and convert it
	key, ok := convert[response.StatusCode/100]
	if !ok {
		key = "errors"
		log.Error("Not able to get the response status code")
	}

	// increase the number
	c.Incr(key, 1)
}

// load
// load is responsoble for generating the required number of requests
func load() {
	timeout := time.Tick(runtimeDuration)
	interval := time.Duration((1000000*loadgenConfig.Burst)/loadgenConfig.Rate) * time.Microsecond
	period := time.Tick(interval)
	for {
		//fire off burst
		for i := 0; i < loadgenConfig.Burst; i++ {
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

	// set the runtime to time.Duration type
	runtimeDuration = time.Duration(loadgenConfig.Runtime) * time.Second

	logger, err := log.LoggerFromConfigAsFile(loadgenConfig.LogPath)
	if err != nil {
		log.Errorf("Cannot open config file %s\n", err)
		return
	}

	log.ReplaceLogger(logger)

	go load()

	// sleep to make sure all of the request are served
	time.Sleep(runtimeDuration + time.Duration(2*loadgenConfig.TimeoutMS)*time.Millisecond)

	// print results
	fmt.Printf("total: \t%d\n", c.Get("total"))
	fmt.Printf("100s: \t%d\n", c.Get("100s"))
	fmt.Printf("200s: \t%d\n", c.Get("200s"))
	fmt.Printf("300s: \t%d\n", c.Get("300s"))
	fmt.Printf("400s: \t%d\n", c.Get("400s"))
	fmt.Printf("500s: \t%d\n", c.Get("500s"))
	fmt.Printf("errors: \t%d\n", c.Get("errors"))
}
