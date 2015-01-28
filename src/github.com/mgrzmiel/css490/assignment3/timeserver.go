// CSS 490
// Magdalena Grzmiel
// Assignments #2
// Copyright 2015 Magdalena Grzmiel
// This program is an example of personlized http server
// which prints a more personalized message for logged-in users.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"html"
	"html/template"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Context struct {
	Name        string
	CurrentTime string
}

// sync lock for cuncurrent accessing sessions object
var sessionsSyncLoc *sync.RWMutex

// declare the map for uuid and user's names
var sessions map[string]string

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
	templatPath := "templates/" + fileName + ".tmpl"
	menuPath := "templates/main.tmpl"
	tmpl, err := tmpl.ParseFiles(menuPath, templatPath)
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
	templatPath := "templates/" + fileName + ".tmpl"
	menuPath := "templates/main.tmpl"
	tmpl, err := tmpl.ParseFiles(menuPath, templatPath)
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
		uuid := generateUniqueId()
		sessionsSyncLoc.Lock()
		sessions[uuid] = name
		sessionsSyncLoc.Unlock()

		// save uuid in the cookie
		cookie := http.Cookie{Name: "uuid", Value: uuid, Path: "/"}
		http.SetCookie(res, &cookie)

		// redirect to /index.html endpoint
		http.Redirect(res, req, "/index.html", http.StatusFound)
	} else {
		// if the provided input - name is empty, display this message
		res.Header().Set("Content-Type", "text/html")
		loadTemplate(res, "noName")
	}
}

// generateUniqueId
// This function generates univerally unique identifier for cookie
func generateUniqueId() string {
	cmd := exec.Command("/usr/bin/uuidgen")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	uuid := out.String()
	uuid = strings.Replace(uuid, "\n", "", 1)
	return uuid
}

// loginForm function
// If user is not login, it displays login form
// Otherwise display the greeting message
func loginForm(res http.ResponseWriter, req *http.Request) {
	//check if the user is login
	name, correctlyLogIn := getNameAndCookie(res, req)
	res.Header().Set("Content-Type", "text/html")
	if !correctlyLogIn {
		loadTemplate(res, "loginForm")
	} else {
		greetingContex := Context{
			Name: name,
		}
		loadTemplateWitData(res, "greetingMessage", greetingContex)
	}
}

// invalidate cookie
// It invalidates cookies since no name exists for that uuid in map
func invalidateCookie(res http.ResponseWriter) {
	// set the experiation date to last year
	expire := time.Now().AddDate(-1, 0, 0)
	cookie := http.Cookie{Name: "uuid", Path: "/", Expires: expire}
	http.SetCookie(res, &cookie)
}

// logout
// It invalidates the cookie since user is no longer login
// and displays good bye message
func logOut(res http.ResponseWriter, req *http.Request) {
	invalidateCookie(res)
	loadTemplate(res, "logout")
}

// getNameAndCookie
// It checks if the cookie is set up and if the name for that cookie exists in map.
// Based on that, it sets up the correctlyLogIn variable.
func getNameAndCookie(res http.ResponseWriter, req *http.Request) (string, bool) {
	var name string
	var ok bool
	var cookie, err = req.Cookie("uuid")

	//correctlyLogIn - means that both cookie and name exists
	correctlyLogIn := false

	// if the cookie is set up
	if err == nil {

		// retrive the name, before the access to map, lock it
		sessionsSyncLoc.RLock()
		name, ok = sessions[cookie.Value]
		sessionsSyncLoc.RUnlock()

		if ok {
			// if the name exists, set correctllyLogIn to true
			correctlyLogIn = true
		} else {
			// no name so invalidate cookie
			invalidateCookie(res)
		}
	}

	return name, correctlyLogIn
}

// getTime
// It is called when the /time endpoint is used
// It displayes the time on the webside
func getTime(res http.ResponseWriter, req *http.Request) {
	now := time.Now().Format("3:04:05 PM")
	displayName := ""
	name, correctlyLogIn := getNameAndCookie(res, req)
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
	sessions = make(map[string]string)
	sessionsSyncLoc = new(sync.RWMutex)

	// parse the flags
	flag.IntVar(&port, "port", 8080, "used port")
	flag.BoolVar(&version, "V", false, "version of the program")
	flag.Parse()

	// if user type -V, the V flag is set up to true
	if version {
		// display the information about the version
		fmt.Println("version 2.0")
	} else {
		// adding the styles
		http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("styles"))))
		// otherwise run the server
		portNr := strconv.Itoa(port)
		http.HandleFunc("/time", getTime)
		http.HandleFunc("/", unknownRoute)
		http.HandleFunc("/index.html", loginForm)
		http.HandleFunc("/login", logIn)
		http.HandleFunc("/logout", logOut)

		err := http.ListenAndServe(":"+portNr, Log(http.DefaultServeMux))
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}
}
