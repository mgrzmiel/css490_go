// CSS 490
// Magdalena Grzmiel
// Assignments #6
// Copyright 2015 Magdalena Grzmiel
// monitorconfig is responsible for reading and parsing all of the flags for monitor.

package monitorConfig

import (
	"flag"
)

const (
	DEFAULT_LOG_PATH = "etc/monitorlog.xml"
)

var Targets string
var SampleIntervalSec int
var RuntimeSec int
var LogPath string

func init() {
	// parse the flags
	flag.StringVar(&Targets, "targets", "", "rate")
	flag.IntVar(&SampleIntervalSec, "sample-interval-sec", 10, "sample-interval-sec")
	flag.IntVar(&RuntimeSec, "runtime-sec", 30, "runtime-sec")
	flag.StringVar(&LogPath, "log", DEFAULT_LOG_PATH, "name of log config file")
	flag.Parse()
}
