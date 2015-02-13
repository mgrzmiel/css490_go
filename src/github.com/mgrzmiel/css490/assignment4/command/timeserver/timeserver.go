// CSS 490
// Magdalena Grzmiel
// Assignments #4
// Copyright 2015 Magdalena Grzmiel
// This program is an example of personlized http server
// which using templates and log messages.

package main

import (
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/mgrzmiel/css490/assignment4/lib/config"
	"github.com/mgrzmiel/css490/assignment4/lib/cookieBasedSessionManager"
	"html/template"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// const (
// 	DEFAULT_LOG_PATH      = "etc/seelog.xml"
// 	DEFAULT_TEMPLATE_PATH = "templates/"
// )

// structure which keeps data which is passed to be displayed in template
type Context struct {
	Name        string
	CurrentTime string
	UTCTime     string
}

// declare variable which is a path for getting templates
// var templatePath string
var portNr string
var count int
var lock *sync.RWMutex
var number int

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
	login := cookieBasedSessionManager.Login(res, req)

	if login {
		log.Trace("User correctlly logged in. Redirecting to /index.html")

		// redirect to /index.html endpoint
		http.Redirect(res, req, "/index.html", http.StatusFound)
	} else {
		log.Trace("Could not log in user.")

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
}

// logout
// Log outs user in cookieBasedSessionManager
func logOut(res http.ResponseWriter, req *http.Request) {
	log.Trace("Logging out user")

	cookieBasedSessionManager.Logout(res, req)
	loadTemplate(res, "logout", nil)
}

// aboutUs
// Display simple information about this program
func aboutUs(res http.ResponseWriter, req *http.Request) {
	loadTemplate(res, "aboutUs", nil)
}

//limit
func limit(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	fmt.Printf("limit is called \n")
	fmt.Printf("max %s\n", config.MaxInflight)
	fmt.Printf("count %s\n", count)
	number++
	fmt.Printf("function is called times: %s\n", number)

	return func(res http.ResponseWriter, req *http.Request) {
		if config.MaxInflight == -1 {
			handler(res, req)
		}

		lock.Lock()
		if config.MaxInflight > count {
			count++
			fmt.Printf("Count after increment %s\n", count)
			lock.Unlock()
			handler(res, req)
			lock.Lock()
			count--
			fmt.Printf("Count after decrement %s\n", count)
			lock.Unlock()
		} else {
			lock.Unlock()
			fmt.Print("to many requests \n")
			res.WriteHeader(http.StatusInternalServerError)
		}
	}
}

// 	lock.Lock()
// 	if config.MaxInflight > count {

// 		count++
// 		fmt.Printf("Count after increment %s\n", count)
// 		lock.Unlock()
// 		return func(res http.ResponseWriter, req *http.Request) {
// 			handler(res, req)

// 			lock.Lock()
// 			count--
// 			fmt.Printf("Count after decrement %s\n", count)
// 			lock.Unlock()
// 		}
// 	} else {
// 		lock.Unlock()
// 		return func(res http.ResponseWriter, req *http.Request) {
// 			fmt.Print("to many requests \n")
// 			res.WriteHeader(http.StatusInternalServerError)
// 		}
// 	}
// }

// getTime
// Displayes the time on the webside
func getTime(res http.ResponseWriter, req *http.Request) {
	fmt.Print("time func is called \n")
	now := time.Now()
	nowLoc := now.Format("3:04:05 PM")
	nowUTC := now.UTC().Format("15:04:05")
	displayName := ""
	name, correctlyLogIn := cookieBasedSessionManager.GetSession(res, req)

	if correctlyLogIn {
		displayName = `, ` + name
	}

	res.Header().Set("Content-Type", "text/html")
	timeContex := Context{
		Name:        displayName,
		CurrentTime: nowLoc,
		UTCTime:     nowUTC,
	}

	delay := rand.NormFloat64()*config.DeviationMs + config.AvgResponseMs
	var delayTime time.Duration = time.Duration(delay)
	if delayTime < 0 {
		return
	}

	time.Sleep(delayTime * time.Millisecond)
	loadTemplate(res, "time", &timeContex)
}

// unknownRoute
// Handels unknown routes. Displays following message:
// "These are not the URLs you're looking for"
// It also sets the status code to 404
func unknownRoute(res http.ResponseWriter, req *http.Request) {
	log.Trace("Uknown route.")

	res.WriteHeader(http.StatusNotFound)
	res.Header().Set("Content-Type", "text/html")
	loadTemplate(res, "unknownRoute", nil)
}

// main function
// This function is responsible for the flow of whole program
func main() {
	// if user type -V, the V flag is set up to true
	if config.Version {
		// display the information about the version
		fmt.Println("version 4.0")
	} else {
		logger, err := log.LoggerFromConfigAsFile(config.LogPath)
		if err != nil {
			log.Errorf("Cannot open config file %s\n", err)
			return
		}

		lock = new(sync.RWMutex)
		count = 0
		number = 0

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

		err = http.ListenAndServe(":"+portNr, Log(http.DefaultServeMux))
		if err != nil {
			log.Errorf("ListenAndServe: %s\n", err)
		}
	}
}
