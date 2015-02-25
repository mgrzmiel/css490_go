// CSS 490
// Magdalena Grzmiel
// Assignments #4
// Copyright 2015 Magdalena Grzmiel
// config is responsible for reading and parsing all of the flags.

package config

import (
	"flag"
	"strings"
)

const (
	DEFAULT_LOG_PATH      = "etc/seelog.xml"
	DEFAULT_TEMPLATE_PATH = "templates/"
	DEFAULT_HOST_NAME     = "localhost"
)

var Port int
var Version bool
var LogPath string
var TemplatePath string
var Authport int
var Authhost string
var AuthtimeoutMs int
var AvgResponseMs float64
var DeviationMs float64
var MaxInflight int
var Dumpfile string
var CheckpointInterval int

func init() {
	// parse the flags
	flag.IntVar(&Port, "port", 8080, "used port")
	flag.BoolVar(&Version, "V", false, "version of the program")
	flag.StringVar(&TemplatePath, "templates", DEFAULT_TEMPLATE_PATH, "path to the templates")
	flag.StringVar(&LogPath, "log", DEFAULT_LOG_PATH, "name of log config file")
	flag.IntVar(&Authport, "authport", 8070, "auth_used port")
	flag.StringVar(&Authhost, "authhost", DEFAULT_HOST_NAME, "name of the host")
	flag.IntVar(&AuthtimeoutMs, "authtime-ms", 10, "time to wait for response")
	flag.Float64Var(&AvgResponseMs, "avg-response-ms", 0, "time to wait for response")
	flag.Float64Var(&DeviationMs, "response-deviation-ms", 0, "time to wait for response")
	flag.IntVar(&MaxInflight, "max-inflight", 0, "maximum number of in-flight time requests")
	flag.StringVar(&Dumpfile, "dumpfile", "", "path to the dumpfile")
	flag.IntVar(&CheckpointInterval, "checkpoint-interval", -1, "checkpoint interval")

	flag.Parse()

	if !strings.HasSuffix(TemplatePath, "/") {
		TemplatePath += "/"
	}
}
