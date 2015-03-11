// CSS 490
// Magdalena Grzmiel
// Assignments #6
// Copyright 2015 Magdalena Grzmiel
// loadgenconfig is responsible for reading and parsing all of the flags for load generator.

package loadgenConfig

import (
	"flag"
)

const (
	DEFAULT_LOG_PATH = "etc/loadgenlog.xml"
)

var Rate int
var Burst int
var TimeoutMS int
var Url string
var Runtime int
var LogPath string

func init() {
	// parse the flags
	flag.IntVar(&Rate, "rate", 200, "rate")
	flag.IntVar(&Burst, "burst", 20, "burst")
	flag.IntVar(&TimeoutMS, "timeout-ms", 400, "timeoutMS")
	flag.StringVar(&Url, "url", "http://localhost:8080/time", "url")
	flag.IntVar(&Runtime, "runtime", 20, "runtime")
	flag.StringVar(&LogPath, "log", DEFAULT_LOG_PATH, "name of log config file")

	flag.Parse()
}
