// CSS 490
// Magdalena Grzmiel
// Assignments #4
// Copyright 2015 Magdalena Grzmiel
// authserver is responsible for managing the data about the logged users

package main

import (
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/mgrzmiel/css490/assignment5/lib/config"
	"github.com/mgrzmiel/css490/assignment5/lib/sessionManager"
	"net/http"
	"strconv"
	"time"
)

// declare the map for uuid and user's names
var sessions *sessionManager.Sessions

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
		return
	}
	cookie := req.Form.Get("cookie")
	name := req.Form.Get("name")
	if name == "" || cookie == "" {
		res.WriteHeader(http.StatusBadRequest)
	} else {
		res.WriteHeader(http.StatusOK)
		sessions.SetSession(name, cookie)
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
		return
	}
	cookie := req.Form.Get("cookie")
	if cookie == "" {
		res.WriteHeader(http.StatusBadRequest)
	} else {
		name, ok := sessions.GetSession(cookie)
		res.WriteHeader(http.StatusOK)
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
	http.HandleFunc("/", unknownRequest)

	err = http.ListenAndServe(":"+portNr, Log(http.DefaultServeMux))
	if err != nil {
		log.Errorf("ListenAndServe: %s\n", err)
	}
}
