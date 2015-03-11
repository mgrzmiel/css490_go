// CSS 490
// Magdalena Grzmiel
// Assignments #6
// Copyright 2015 Magdalena Grzmiel
// authserver is responsible for managing the data about the logged users

package main

import (
	"encoding/json"
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/mgrzmiel/css490/assignment6/lib/config"
	"github.com/mgrzmiel/css490/assignment6/lib/counter"
	"github.com/mgrzmiel/css490/assignment6/lib/sessionManager"
	"net/http"
	"strconv"
	"time"
)

// declare the map for uuid and user's names
var sessions *sessionManager.Sessions
var (
	authServerCounter = counter.New()
)

// Log function
// Wrapper around DefaultServeMutex for printing each request
// before it's being handled by a handle function
func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		log.Infof("HTTP request. URL: %s", req.URL.Path)
		handler.ServeHTTP(res, req)
	})
}

// setFunc
// It is responsible for retriving the parameters from set request
// and if they are valid, set the session
func setFunc(res http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		log.Errorf("Problem with retriving data from url %s", err)
		res.WriteHeader(http.StatusBadRequest)
		authServerCounter.Incr("400s", 1)
		return
	}
	cookie := req.Form.Get("cookie")
	name := req.Form.Get("name")
	if name == "" || cookie == "" {
		res.WriteHeader(http.StatusBadRequest)
		authServerCounter.Incr("400s", 1)
	} else {
		res.WriteHeader(http.StatusOK)
		sessions.SetSession(name, cookie)
		authServerCounter.Incr("set-cookie", 1)
		authServerCounter.Incr("200s", 1)
	}
}

// getFunc
// It is responsible for retrving the cookie from url and then
// based on the cookie retriving the name and return the name.
func getFunc(res http.ResponseWriter, req *http.Request) {

	err := req.ParseForm()
	if err != nil {
		log.Errorf("Problem with retriving data from url %s", err)
		res.WriteHeader(http.StatusBadRequest)
		authServerCounter.Incr("400s", 1)
		return
	}
	cookie := req.Form.Get("cookie")
	if cookie == "" {
		res.WriteHeader(http.StatusBadRequest)
		authServerCounter.Incr("400s", 1)
		authServerCounter.Incr("no-cookie", 1)
	} else {
		name, ok := sessions.GetSession(cookie)
		authServerCounter.Incr("get-cookie", 1)
		res.WriteHeader(http.StatusOK)
		authServerCounter.Incr("200s", 1)
		if ok {
			fmt.Fprintf(res, name)
		}
	}
}

// All other request than set and get are treaded as unknown request
// and return the 404 status code
func unknownRequest(res http.ResponseWriter, req *http.Request) {
	log.Trace("Uknown request")
	res.WriteHeader(http.StatusNotFound)
	authServerCounter.Incr("404s", 1)
}

// monitorFunc
// It displays the json object which presents the statistics
func monitorFunc(res http.ResponseWriter, req *http.Request) {

	// write to temp dictionery
	monitorMap := make(map[string]int)
	monitorMap["set-cookie"] = authServerCounter.Get("set-cookie")
	monitorMap["get-cookie"] = authServerCounter.Get("get-cookie")
	monitorMap["no-cookie"] = authServerCounter.Get("no-cookie")
	monitorMap["200s"] = authServerCounter.Get("200s")
	monitorMap["400s"] = authServerCounter.Get("400s")
	monitorMap["404s"] = authServerCounter.Get("404s")

	// marshall the data
	data, err := json.Marshal(monitorMap)
	if err != nil {
		log.Errorf("Not able to marshall the data")
		return
	} else {
		res.Header().Set("Content-Type", "text/json")
		fmt.Fprintf(
			res,
			string(data),
		)
	}
}

func periodicallySaveSessions(fileName string, checkpointInterval time.Duration) {
	performWriting := time.Tick(checkpointInterval * time.Second)
	for _ = range performWriting {
		sessions.WriteToFile(fileName)
	}
}

// main function
// This function is responsible for the flow of whole program
func main() {

	// get the dumpfile name
	fileName := config.Dumpfile

	// create new session object, if the file exist the map will contain the data
	sessions = sessionManager.New(fileName)

	// the interval time was provided, start to save the data to file
	if config.CheckpointInterval >= 0 {
		var checkpointInterval time.Duration = time.Duration(config.CheckpointInterval)
		go periodicallySaveSessions(fileName, checkpointInterval)
	}

	logger, err := log.LoggerFromConfigAsFile(config.LogPath)
	if err != nil {
		log.Errorf("Cannot open config file %s\n", err)
		return
	}

	log.ReplaceLogger(logger)
	log.Info("Starging authserver")
	log.Debugf("Authport: %s", config.Authport)

	//run the server
	portNr := strconv.Itoa(config.Authport)
	http.HandleFunc("/set", setFunc)
	http.HandleFunc("/get", getFunc)
	http.HandleFunc("/monitor", monitorFunc)
	http.HandleFunc("/", unknownRequest)

	err = http.ListenAndServe(":"+portNr, Log(http.DefaultServeMux))
	if err != nil {
		log.Errorf("ListenAndServe: %s\n", err)
	}
}
