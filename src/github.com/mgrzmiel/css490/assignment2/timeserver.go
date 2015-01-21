// CSS 490
// Magdalena Grzmiel
// Assignments #2
// Copyright 2015 Magdalena Grzmiel
// This program is an example of personlized http server.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"html"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

// declare the RWMutex structure
var sessionsSyncLoc *sync.RWMutex

// declare the map for uuid and user's names
var sessions map[string]string

// logIn function
// If the name is provided in the request, the login function generate
// uuid and add it together with name to map. Then redirect page to index.html
//endpoint.
// If the name is an empty string, it just display simple message
func logIn(res http.ResponseWriter, req *http.Request) {
	fmt.Println(req)
	//retrive the name form URL
	name := req.FormValue("name")
	name = html.EscapeString(name)
	if name != "" {
		id := generateUniqueId()                                  // generate uuid
		sessionsSyncLoc.Lock()                                    // before modifying the map, lock it
		sessions[id] = name                                       // add name with uuid to map
		sessionsSyncLoc.Unlock()                                  // unlock map
		cookie := http.Cookie{Name: "uuid", Value: id, Path: "/"} // create cookie
		http.SetCookie(res, &cookie)                              //set the cookie and redirect to the /index.htm endpoint
		http.Redirect(res, req, "/index.html", http.StatusFound)
	} else {
		// if the provided input - name is empty, display this message
		res.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(
			res,
			`<html>
			<body>
			<form action="login">
			  C'mon, I need a name.
			</form>
			</p>
			</body>
			</html>`,
		)
	}
}

// generateUniqueId
// This function generates cookie containing a univerally unique identifier
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
	fmt.Println(req)
	//check if the user is login -
	name, correctlyLogIn := getNameAndCookie(res, req)
	res.Header().Set("Content-Type", "text/html")
	if !correctlyLogIn {
		fmt.Fprintf(
			res,
			`<html>
			<body>
			<form action="login">
			  What is your name, Earthling?
			  <input type="text" name="name" size="50">
			  <input type="submit">
			</form>
			</p>
			</body>
			</html>`,
		)
	} else {
		fmt.Fprintf(
			res,
			`<html>
			<body>
			<p> Greetings, `+name+`</p>
			</body>
			</html>`,
		)
	}
}

// invalidate cookie
// It invalidates cookies since no name exists for that cookie in map
func invalidateCookie(res http.ResponseWriter) {
	// set the experiation date to yesterday
	expire := time.Now().AddDate(-1, 0, -1)
	cookie2 := http.Cookie{Name: "uuid", Path: "/", Expires: expire}
	http.SetCookie(res, &cookie2)
}

//logout
// It invalidate the cookie since user is no longer liginh
// and displat good bye message
func logOut(res http.ResponseWriter, req *http.Request) {
	fmt.Println(req)
	invalidateCookie(res)
	fmt.Fprintf(
		res,
		`<html>
		<head>
		<META http-equiv="refresh" content="10;URL=/index.html">
		<body>
		<p>Good-bye.</p>
		</body>
		</html>`,
	)
}

//getNameAndCookie
//check if the cookie is set up and if the name for that cookie exists in map
// based on that, it sets up the correctlyLogIn variable.
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
		// if the name exist, set correctllyLogIn to true
		if ok {
			correctlyLogIn = true
			// no name so invalidate cookie
		} else {
			invalidateCookie(res)
		}
	}

	return name, correctlyLogIn
}

// getTime
// It is caleed when the /time endpoint is used
// It displayes the time on the webside
func getTime(res http.ResponseWriter, req *http.Request) {
	fmt.Println(req)
	now := time.Now().Format("3:04:05 PM")
	name, correctlyLogIn := getNameAndCookie(res, req)
	res.Header().Set("Content-Type", "text/html")
	if correctlyLogIn {
		fmt.Fprintf(
			res,
			`<doctype html>
	        <html>
			<head>
			<style>
			p {font-size: xx-large}
			span.time {color: red}
			</style>
			</head>
			<body>
			<p>The time is now <span class="time">`+now+`</span>, `+name+
				`.</p>
			</body>
			</html>`,
		)
	} else {
		fmt.Fprintf(
			res,
			`<doctype html>
	        <html>
			<head>
			<style>
			p {font-size: xx-large}
			span.time {color: red}
			</style>
			</head>
			<body>
			<p>The time is now <span class="time">`+now+`</span>.</p>
			</body>
			</html>`,
		)
	}
}

// unknownRoute
// If the endpint is not known, this method is called
// It displays following message:
// "These are not the URLs you're looking for"
// It also sets the status code to 404
func unknownRoute(res http.ResponseWriter, req *http.Request) {
	fmt.Println(req)
	res.WriteHeader(http.StatusNotFound)
	res.Header().Set("Content-Type", "html")
	fmt.Fprintf(
		res,
		`<html>
		<body>
		<p>These are not the URLs you're looking for.</p>
		</body>
		</html>`,
	)
}

// main function
// This function is responsible for the floow of whole program
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
		fmt.Println("version 1.9")
		// otherwise run the server
	} else {
		portNr := strconv.Itoa(port)
		http.HandleFunc("/time", getTime)
		http.HandleFunc("/", unknownRoute)
		http.HandleFunc("/index.html", loginForm)
		http.HandleFunc("/login", logIn)
		http.HandleFunc("/logout", logOut)
		err := http.ListenAndServe(":"+portNr, nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}
}
