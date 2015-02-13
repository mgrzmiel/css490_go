// CSS 490
// Magdalena Grzmiel
// Assignments #4
// Copyright 2015 Magdalena Grzmiel

package authClient

import (
	log "github.com/cihub/seelog"
	"io/ioutil"
	"net/http"
)

func SetRequest(portName string, port string, cookie string, name string) {
	// requestURL := "http://localhost:8070" + "/set?cookie=" + cookie + "&name=" + name
	requestURL := "http://" + portName + ":" + port + "/set?cookie=" + cookie + "&name=" + name
	res, err := http.Get(requestURL)
	if err != nil {
		log.Errorf("Set request error %s", err)
	}
	res.Body.Close()
}

func GetRequest(portName string, port string, cookie string) string {
	// requestURL := "http://localhost:8070" + "/get?cookie=" + cookie
	requestURL := "http://" + portName + ":" + port + "/get?cookie=" + cookie

	res, err := http.Get(requestURL)
	if err != nil {
		log.Errorf("Ger request error %s", err)
	}
	data, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Errorf("error %s", err)
	}
	name := string(data[:])
	return name
}
