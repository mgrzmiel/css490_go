// CSS 490
// Magdalena Grzmiel
// Assignments #4
// Copyright 2015 Magdalena Grzmiel
// Authclient send request to the server.

package authClient

import (
	log "github.com/cihub/seelog"
	"github.com/mgrzmiel/css490/assignment5/lib/config"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// SetRequest
// It sends set request to the server to login the user
func SetRequest(cookie string, name string) bool {
	// create the request
	requestURL := "http://" + config.Authhost + ":" + strconv.Itoa(config.Authport) + "/set?cookie=" + cookie + "&name=" + name

	timeout := time.Duration(time.Duration(config.AuthtimeoutMs) * time.Millisecond)
	client := http.Client{
		Timeout: timeout,
	}

	res, err := client.Get(requestURL)
	if err != nil {
		log.Errorf("Set request error %s", err)
		return false
	}
	res.Body.Close()
	return true
}

// GetRequest
// It sends get request to the srever to get the name for the given cookie-uuid
func GetRequest(cookie string) (string, bool) {
	// create the rerquest
	requestURL := "http://" + config.Authhost + ":" + strconv.Itoa(config.Authport) + "/get?cookie=" + cookie

	// send the request to server
	timeout := time.Duration(time.Duration(config.AuthtimeoutMs) * time.Millisecond)
	client := http.Client{
		Timeout: timeout,
	}
	res, err := client.Get(requestURL)
	if err != nil {
		log.Errorf("Get request error %s", err)
		return "", false
	}

	data, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Errorf("Error while reading the request: %s", err)
		return "", false
	}
	name := string(data[:])
	return name, true
}
