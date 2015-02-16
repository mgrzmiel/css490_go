// CSS 490
// Magdalena Grzmiel
// Assignments #4
// Copyright 2015 Magdalena Grzmiel

package main

import (
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/mgrzmiel/css490/assignment4/lib/config"
	"github.com/mgrzmiel/css490/assignment4/lib/sessionManager"
	"net/http"
	"strconv"
	"time"
)

const (
	DEFAULT_LOG_PATH = "etc/seelog.xml"
)

// declare the map for uuid and user's names
var sessions *sessionManager.Sessions

func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		log.Infof("HTTP request. URL: %s", req.URL.Path)
		handler.ServeHTTP(res, req)
	})
}

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
		return
	} else {
		res.WriteHeader(http.StatusOK)
		sessions.SetSession(name, cookie)
	}
}

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
		return
	} else {
		name, ok := sessions.GetSession(cookie)
		if ok {
			res.WriteHeader(http.StatusOK)
			fmt.Fprintf(res, name)
		} else {
			res.WriteHeader(http.StatusOK)
		}
	}
}

func unknownRequest(res http.ResponseWriter, req *http.Request) {
	log.Trace("Uknown request")
	res.WriteHeader(http.StatusNotFound)
}

// main function
// This function is responsible for the flow of whole program
func main() {
	fileName := config.Dumpfile
	sessions = sessionManager.New(fileName)

	// sessions.SessionsMap = readFromFile(fileName)
	// fmt.Print("Print sessions maps")
	// for key, value := range sessions.SessionsMap {
	// 	fmt.Println("Key:", key, "Value:", value)
	// }

	if config.CheckpointInterval != -1 {
		var checkpointInterval time.Duration = time.Duration(config.CheckpointInterval)
		go sessions.WriteToFile(fileName, checkpointInterval)
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
