// CSS 490
// Magdalena Grzmiel
// Assignments #6
// Copyright 2015 Magdalena Grzmiel
// This program is an example of personlized http server
// which using templates, authserver and log messages.

package main

import (
	"encoding/json"
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/mgrzmiel/css490/assignment6/lib/config"
	"github.com/mgrzmiel/css490/assignment6/lib/cookieBasedSessionManager"
	"github.com/mgrzmiel/css490/assignment6/lib/counter"
	"html/template"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// structure which keeps data which is passed to be displayed in template
type Context struct {
	Name        string
	CurrentTime string
	UTCTime     string
}

// declare variable which is a path for getting templates
var portNr string
var count int
var lock *sync.RWMutex

var (
	timeCounter = counter.New()
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

// loadTemplates
// Loads the templets by the provided template name with using provided data
func loadTemplate(res http.ResponseWriter, fileName string, data *Context) {
	tmpl := template.New("main")
	log.Debugf("Template to be rendered: %s", fileName)
	templatePath := config.TemplatePath + fileName + ".tmpl"
	mainPath := config.TemplatePath + "main.tmpl"
	menuPath := config.TemplatePath + "menu.tmpl"

	tmpl, err := tmpl.ParseFiles(mainPath, menuPath, templatePath)
	if err != nil {
		log.Errorf("Error during parsing template %s: %s", fileName, err)
		return
	}

	err = tmpl.ExecuteTemplate(res, "main", data)
	if err != nil {
		log.Errorf("Error during executing template %s: %s", fileName, err)
		return
	}
}

// logIn function
// Tryies to log in the user. If successful the page is then redirect to
// index.html endpoint, otherwise displays simple informationtion message.
func logIn(res http.ResponseWriter, req *http.Request) {
	// check if the user is login
	login, failureMessage := cookieBasedSessionManager.Login(res, req)

	if login {
		timeCounter.Incr("login", 1)
		log.Trace("User correctlly logged in. Redirecting to /index.html")

		// redirect to /index.html endpoint
		http.Redirect(res, req, "/index.html", http.StatusFound)
	} else {
		log.Trace("Could not log in user.")
		timeCounter.Incr("200s", 1)

		// if the provided input - name is empty, display this message
		res.Header().Set("Content-Type", "text/html")
		loadTemplate(res, failureMessage, nil)
	}
}

// loginForm function
// If user is not logged in, it displays login form
// Otherwise display the greeting message
func loginForm(res http.ResponseWriter, req *http.Request) {
	//check if the user is login
	name, correctlyLogIn := cookieBasedSessionManager.GetSession(res, req)

	res.Header().Set("Content-Type", "text/html")
	if correctlyLogIn {
		log.Tracef("Logged in user. Name: %s", name)

		greetingContex := Context{Name: name}
		loadTemplate(res, "greetingMessage", &greetingContex)
	} else {
		log.Trace("Unknown user.")

		loadTemplate(res, "loginForm", nil)
	}
	timeCounter.Incr("200s", 1)
}

// logout
// Log outs user in cookieBasedSessionManager
func logOut(res http.ResponseWriter, req *http.Request) {
	log.Trace("Logging out user")

	cookieBasedSessionManager.Logout(res, req)
	loadTemplate(res, "logout", nil)
	timeCounter.Incr("200s", 1)
}

// aboutUs
// Display simple information about this program
func aboutUs(res http.ResponseWriter, req *http.Request) {
	loadTemplate(res, "aboutUs", nil)
	timeCounter.Incr("200s", 1)
}

// limit
// Limit the number of concurrent request which can be handled by server
func limit(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		if config.MaxInflight <= 0 {
			handler(res, req)
		} else {
			lock.Lock()
			if config.MaxInflight > count {
				count++
				lock.Unlock()
				handler(res, req)
				lock.Lock()
				count--
				lock.Unlock()
			} else {
				lock.Unlock()
				res.WriteHeader(http.StatusInternalServerError)
				timeCounter.Incr("500s", 1)
			}
		}
	}
}

// getTime
// Displayes the time on the webside
func getTime(res http.ResponseWriter, req *http.Request) {
	now := time.Now()
	nowLoc := now.Format("3:04:05 PM")
	nowUTC := now.UTC().Format("15:04:05")
	displayName := ""
	name, correctlyLogIn := cookieBasedSessionManager.GetSession(res, req)

	if correctlyLogIn {
		timeCounter.Incr("time-user", 1)
		displayName = `, ` + name
	} else {
		timeCounter.Incr("time-anon", 1)
	}

	res.Header().Set("Content-Type", "text/html")
	timeContex := Context{
		Name:        displayName,
		CurrentTime: nowLoc,
		UTCTime:     nowUTC,
	}

	delay := rand.NormFloat64()*config.DeviationMs + config.AvgResponseMs
	var delayTime time.Duration = time.Duration(delay)
	if delayTime > 0 {
		time.Sleep(delayTime * time.Millisecond)
	}

	loadTemplate(res, "time", &timeContex)
	timeCounter.Incr("200s", 1)
}

// unknownRoute
// Handels unknown routes. Displays following message:
// "These are not the URLs you're looking for"
// It also sets the status code to 404
func unknownRoute(res http.ResponseWriter, req *http.Request) {
	log.Trace("Uknown route.")

	res.WriteHeader(http.StatusNotFound)
	timeCounter.Incr("404s", 1)
	res.Header().Set("Content-Type", "text/html")
	loadTemplate(res, "unknownRoute", nil)
}

// monitor
// It displays the json object which presents the statistics
func monitor(res http.ResponseWriter, req *http.Request) {

	// write to temp dictionery
	monitorMap := make(map[string]int)
	monitorMap["login"] = timeCounter.Get("login")
	monitorMap["time-user"] = timeCounter.Get("time-user")
	monitorMap["time-anon"] = timeCounter.Get("time-anon")
	monitorMap["200s"] = timeCounter.Get("200s")
	monitorMap["404s"] = timeCounter.Get("404s")
	monitorMap["500s"] = timeCounter.Get("500s")

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

// main function
// This function is responsible for the flow of whole program
func main() {
	// if user type -V, the V flag is set up to true
	if config.Version {
		// display the information about the version
		fmt.Println("version 6.0")
	} else {
		logger, err := log.LoggerFromConfigAsFile(config.LogPath)
		if err != nil {
			log.Errorf("Cannot open config file %s\n", err)
			return
		}

		lock = new(sync.RWMutex)
		count = 0

		log.ReplaceLogger(logger)

		log.Info("Starging server")
		log.Debugf("Port: %s", config.Port)
		log.Debugf("Template path: %s", config.TemplatePath)
		log.Debugf("logPath: %s", config.LogPath)

		// adding the styles
		http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("styles"))))

		//run the server
		portNr = strconv.Itoa(config.Port)
		http.HandleFunc("/time", limit(getTime))
		http.HandleFunc("/", unknownRoute)
		http.HandleFunc("/index.html", loginForm)
		http.HandleFunc("/login", logIn)
		http.HandleFunc("/logout", logOut)
		http.HandleFunc("/aboutUs", aboutUs)
		http.HandleFunc("/monitor", monitor)

		err = http.ListenAndServe(":"+portNr, Log(http.DefaultServeMux))
		if err != nil {
			log.Errorf("ListenAndServe: %s\n", err)
		}
	}
}
