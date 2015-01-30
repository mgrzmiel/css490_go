// CSS 490
// Magdalena Grzmiel
// Assignments #2
// Copyright 2015 Magdalena Grzmiel
// This program is an example of personlized http server
// which prints a more personalized message for logged-in users.

package main

import (
	"flag"
	"fmt"
	"github.com/mgrzmiel/css490/assignment3/lib/cookiesManager"
	"github.com/mgrzmiel/css490/assignment3/lib/sessionManager"
	"html"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Context struct {
	Name        string
	CurrentTime string
}

// declare the map for uuid and user's names
var sessions *sessionManager.Sessions
var templatePath string

// Log function
// Wrapper around DefaultServeMutex for printing each request
// before it's being handled by a handle function
func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		fmt.Println(req)
		handler.ServeHTTP(res, req)
	})
}

func loadTemplate(res http.ResponseWriter, fileName string) {
	tmpl := template.New("main")
	templatPath := templatePath + fileName + ".tmpl"
	mainPath := templatePath + "main.tmpl"
	menuPath := templatePath + "menu.tmpl"

	tmpl, err := tmpl.ParseFiles(mainPath, menuPath, templatPath)
	if err != nil {
		fmt.Printf("parsing template: %s\n", err)
		return
	}

	err = tmpl.ExecuteTemplate(res, "main", "")
	if err != nil {
		fmt.Printf("executing template: %s\n", err)
		return
	}
	fmt.Println()
}

func loadTemplateWitData(res http.ResponseWriter, fileName string, data Context) {
	tmpl := template.New("main")
	templatPath := templatePath + fileName + ".tmpl"
	mainPath := templatePath + "main.tmpl"
	menuPath := templatePath + "menu.tmpl"

	tmpl, err := tmpl.ParseFiles(mainPath, menuPath, templatPath)
	if err != nil {
		fmt.Printf("parsing template: %s\n", err)
		return
	}

	err = tmpl.ExecuteTemplate(res, "main", data)
	if err != nil {
		fmt.Printf("executing template: %s\n", err)
		return
	}
	fmt.Println()
}

// logIn function
// If the user name is provided in the request, the login function generates
// uuid and adds it together with name to map.
// Then it redirects page to index.html endpoint.
// If the name is an empty string, it just displays simple message.
func logIn(res http.ResponseWriter, req *http.Request) {
	// retrive the name form URL
	name := req.FormValue("name")
	name = html.EscapeString(name)
	if name != "" {

		uuid := sessions.CreateSession(name)
		cookiesManager.SetCookieValue(res, "uuid", uuid)
		// redirect to /index.html endpoint
		http.Redirect(res, req, "/index.html", http.StatusFound)
	} else {
		// if the provided input - name is empty, display this message
		res.Header().Set("Content-Type", "text/html")
		loadTemplate(res, "noName")
	}
}

// loginForm function
// If user is not login, it displays login form
// Otherwise display the greeting message
func loginForm(res http.ResponseWriter, req *http.Request) {
	//check if the user is login
	cookie, ok := cookiesManager.GetCookieValue(req, "uuid")
	var correctlyLogIn bool
	if ok {
		name, founded := sessions.GetSession(cookie)
		if founded {
			correctlyLogIn = true
			greetingContex := Context{
				Name: name,
			}
			loadTemplateWitData(res, "greetingMessage", greetingContex)
		}
	}

	res.Header().Set("Content-Type", "text/html")
	if !correctlyLogIn {
		loadTemplate(res, "loginForm")
	}
}

// logout
// It invalidates the cookie since user is no longer login
// and displays good bye message
func logOut(res http.ResponseWriter, req *http.Request) {
	cookiesManager.RemoveCookie(res, "uuid")
	loadTemplate(res, "logout")
}

func aboutUs(res http.ResponseWriter, req *http.Request) {
	loadTemplate(res, "aboutUs")
}

// getNameAndCookie
// It checks if the cookie is set up and if the name for that cookie exists in map.
// Based on that, it sets up the correctlyLogIn variable.
// func getNameAndCookie(res http.ResponseWriter, req *http.Request) (string, bool) {
// 	var name string
// 	var ok bool
// 	var cookie, err = req.Cookie("uuid")

// 	//correctlyLogIn - means that both cookie and name exists
// 	correctlyLogIn := false

// 	// if the cookie is set up
// 	if err == nil {

// 		// retrive the name, before the access to map, lock it
// 		sessionsSyncLoc.RLock()
// 		name, ok = sessions[cookie.Value]
// 		sessionsSyncLoc.RUnlock()

// 		if ok {
// 			// if the name exists, set correctllyLogIn to true
// 			correctlyLogIn = true
// 		} else {
// 			// no name so invalidate cookie
// 			invalidateCookie(res)
// 		}
// 	}

// 	return name, correctlyLogIn
// }

// getTime
// It is called when the /time endpoint is used
// It displayes the time on the webside
func getTime(res http.ResponseWriter, req *http.Request) {
	now := time.Now().Format("3:04:05 PM")
	displayName := ""
	correctlyLogIn := false
	var name string
	var founded bool
	cookie, ok := cookiesManager.GetCookieValue(req, "uuid")
	if ok {
		name, founded = sessions.GetSession(cookie)
		if founded {
			correctlyLogIn = true
		} else {
			cookiesManager.RemoveCookie(res, cookie)
		}
	}

	if correctlyLogIn {
		displayName = `, ` + name
	}

	res.Header().Set("Content-Type", "text/html")

	timeContex := Context{
		Name:        displayName,
		CurrentTime: now,
	}
	loadTemplateWitData(res, "time", timeContex)
}

// unknownRoute
// If the endpint is unknown, this method is called.
// It displays following message:
// "These are not the URLs you're looking for"
// It also sets the status code to 404
func unknownRoute(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusNotFound)
	res.Header().Set("Content-Type", "html")
	loadTemplate(res, "unknownRoute")
}

// main function
// This function is responsible for the flow of whole program
func main() {

	var port int
	var version bool
	sessions = sessionManager.New()

	// parse the flags
	flag.IntVar(&port, "port", 8080, "used port")
	flag.BoolVar(&version, "V", false, "version of the program")
	flag.StringVar(&templatePath, "templates", "", "path to the templates")
	flag.Parse()

	// if user type -V, the V flag is set up to true
	if version {
		// display the information about the version
		fmt.Println("version 2.0")
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
			log.Fatal("ListenAndServe: ", err)
		}
	}
}
