// CSS 490
// Magdalena Grzmiel
// Assignments #3
// Copyright 2015 Magdalena Grzmiel
// This program is an example of personlized http server
// which using templates and log messages.

package main

import (
	"flag"
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/mgrzmiel/css490/assignment3/lib/cookieBasedSessionManager"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// structure which keep data which is displayed in template
type Context struct {
	Name        string
	CurrentTime string
	UTCTime     string
}

// declare the map for uuid and user's names
//var sessions *sessionManager.Sessions
var cookieBasedSessions *cookieBasedSessionManager.CookieBasedSessions

// declare variable which is a path for getting templates
var templatePath string

// Log function
// Wrapper around DefaultServeMutex for printing each request
// before it's being handled by a handle function
func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		//fmt.Println(req)
		log.Info(req)
		handler.ServeHTTP(res, req)
	})
}

// loadTemplates
// Loads the templets by the provided template name with using provided data
func loadTemplate(res http.ResponseWriter, fileName string, data *Context) {
	tmpl := template.New("main")
	templatPath := templatePath + fileName + ".tmpl"
	log.Trace("Load the follwoing templete: " + templatePath)
	mainPath := templatePath + "main.tmpl"
	log.Info("main path: " + mainPath)
	menuPath := templatePath + "menu.tmpl"
	log.Info("menu path: " + menuPath)

	tmpl, err := tmpl.ParseFiles(mainPath, menuPath, templatPath)
	if err != nil {
		//fmt.Printf("parsing template: %s\n", err)
		log.Warnf("Parsing template %s", err)
		return
	}

	err = tmpl.ExecuteTemplate(res, "main", data)
	if err != nil {
		//fmt.Printf("executing template: %s\n", err)
		log.Warnf("executing template %s", err)
		return
	}
	fmt.Println()
}

// logIn function
// If the user name is provided in the request, uuid is generated
// and added with name to map. The page isthen redirect to index.html endpoint.
// If the name is an empty string, displays simple message.
func logIn(res http.ResponseWriter, req *http.Request) {
	// retrive the name form URL
	correctlyLogIn := cookieBasedSessions.Login(res, req)

	if correctlyLogIn {
		// redirect to /index.html endpoint
		http.Redirect(res, req, "/index.html", http.StatusFound)
	} else {
		// if the provided input - name is empty, display this message
		res.Header().Set("Content-Type", "text/html")
		loadTemplate(res, "noName", nil)
	}
}

// loginForm function
// If user is not login, it displays login form
// Otherwise display the greeting message
func loginForm(res http.ResponseWriter, req *http.Request) {
	//check if the user is login
	name, correctlyLogIn := cookieBasedSessions.GetSession(res, req)

	res.Header().Set("Content-Type", "text/html")
	if correctlyLogIn {
		greetingContex := Context{Name: name}
		loadTemplate(res, "greetingMessage", &greetingContex)
	} else {
		log.Debug("the cookie was not set up")
		loadTemplate(res, "loginForm", nil)
	}
}

// logout
// It invalidates the cookie since user is no longer login
// and displays good bye message
func logOut(res http.ResponseWriter, req *http.Request) {
	cookieBasedSessions.Logout(res, req)
	loadTemplate(res, "logout", nil)
}

// aboutUs
// It displays simple information about this program
func aboutUs(res http.ResponseWriter, req *http.Request) {
	loadTemplate(res, "aboutUs", nil)
}

// getTime
// It is called when the /time endpoint is used
// It displayes the time on the webside
func getTime(res http.ResponseWriter, req *http.Request) {
	now := time.Now()
	nowLoc := now.Format("3:04:05 PM")
	nowUTC := now.UTC().Format("15:04:05")
	displayName := ""
	name, correctlyLogIn := cookieBasedSessions.GetSession(res, req)

	if correctlyLogIn {
		displayName = `, ` + name
	}

	res.Header().Set("Content-Type", "text/html")
	timeContex := Context{
		Name:        displayName,
		CurrentTime: nowLoc,
		UTCTime:     nowUTC,
	}
	loadTemplate(res, "time", &timeContex)
}

// unknownRoute
// If the endpint is unknown, this method is called.
// It displays following message:
// "These are not the URLs you're looking for"
// It also sets the status code to 404
func unknownRoute(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusNotFound)
	res.Header().Set("Content-Type", "text/html")
	loadTemplate(res, "unknownRoute", nil)
}

// main function
// This function is responsible for the flow of whole program
func main() {

	var port int
	var version bool
	// sessions = sessionManager.New()
	cookieBasedSessions = cookieBasedSessionManager.New()
	var printlogs string

	// parse the flags
	flag.IntVar(&port, "port", 8080, "used port")
	flag.BoolVar(&version, "V", false, "version of the program")
	flag.StringVar(&templatePath, "templates", "", "path to the templates")
	flag.StringVar(&printlogs, "log", "", "name of log config file")
	flag.Parse()

	// if user type -V, the V flag is set up to true
	if version {
		// display the information about the version
		fmt.Println("version 3.0")
	} else {
		// check if the provided path ends with "/", if not add it
		if len(templatePath) > 0 && !strings.HasSuffix(templatePath, "/") {
			templatePath += "/"
		}

		// adding the styles
		http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("styles"))))
		//run the server
		portNr := strconv.Itoa(port)
		http.HandleFunc("/time", getTime)
		http.HandleFunc("/", unknownRoute)
		http.HandleFunc("/index.html", loginForm)
		http.HandleFunc("/login", logIn)
		http.HandleFunc("/logout", logOut)
		http.HandleFunc("/aboutUs", aboutUs)

		err := http.ListenAndServe(":"+portNr, Log(http.DefaultServeMux))
		if err != nil {
			log.Criticalf("ListenAndServe: %s", err)
		}
	}
}
