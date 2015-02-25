package main

import (
	"flag"
	"fmt"
	"github.com/mgrzmiel/css490/assignment5/lib/counter"
	"net/http"
	"time"
)

var rate int      //    = 200
var burst int     //    = 20
var timeoutMS int // = 400
var url string    //     = "http://localhost:8080/time"
var runtime int   //  = 20 * time.Second
var runtimeInMS time.Duration

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

//one request
func request() {
	c.Incr("total", 1)
	client := http.Client{
		Timeout: time.Duration(timeoutMS) * time.Millisecond,
	}
	response, err := client.Get(url)
	if err != nil {
		c.Incr("errors", 1)
		return
	}
	key, ok := convert[response.StatusCode/100]
	if !ok {
		key = "errors"
	}
	c.Incr(key, 1)

	// if response.StatusCode<200{
	// 	c.Incr('100s',1)
	// }else if response.StatusCode<300{
	// 	c.Incr('200', 1)
	//} //etc

	// c.Incr(fmt.Sprintf("%ds", (response.StatusCode /100)*100), 1)
}

func load() {
	timeout := time.Tick(runtimeInMS)
	interval := time.Duration((1000000*burst)/rate) * time.Microsecond
	fmt.Println(interval)
	fmt.Println(runtimeInMS)
	period := time.Tick(interval)
	for {
		//fire off burst
		for i := 0; i < burst; i++ {
			go request()
		}
		// wait for next tick
		<-period
		fmt.Println("p")
		select {
		case <-timeout:
			return
		default:
		}
		//val, ok:= <-timeout
		//poll or timeout

	}

}

func main() {

	flag.IntVar(&rate, "rate", 200, "rate")
	flag.IntVar(&burst, "burst", 20, "burst")
	flag.IntVar(&timeoutMS, "timeout-ms", 400, "timeoutMS")
	flag.StringVar(&url, "url", "http://localhost:8080/time", "url")
	flag.IntVar(&runtime, "runtime", 20, "runtime")

	flag.Parse()

	runtimeInMS = time.Duration(runtime) * time.Millisecond
	//go load()
	load()

	time.Sleep(time.Duration(2*timeoutMS) * time.Millisecond)

	fmt.Printf("total: \t%d\n", c.Get("total"))
	fmt.Printf("100s: \t%d\n", c.Get("100s"))
	fmt.Printf("200s: \t%d\n", c.Get("200s"))
	fmt.Printf("300s: \t%d\n", c.Get("300s"))
	fmt.Printf("400s: \t%d\n", c.Get("400s"))
	fmt.Printf("500s: \t%d\n", c.Get("500s"))
	fmt.Printf("errors: \t%d\n", c.Get("errors"))
}
