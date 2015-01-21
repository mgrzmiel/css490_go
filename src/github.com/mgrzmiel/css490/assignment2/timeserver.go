// CSS 490
// Magdalena Grzmiel
// Assignments #2
// Copyright 2015 Magdalena Grzmiel
// This is a simple example of http server.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var sessions map[string]string

// logIn is a function which read the user name from request
func logIn(res http.ResponseWriter, req *http.Request) {
	name := req.URL.Query().Get("name")
	id := generateUniqueId()
	sessions[id] = name
	cookie := http.Cookie{Name: "uuid", Value: id, Path: "/"}
	http.SetCookie(res, &cookie)
	http.Redirect(res, req, "/index.html", http.StatusFound)
}

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
func loginForm(res http.ResponseWriter, req *http.Request) {
	name, correctlyLogIn:= getNameAndCookie(res, req)
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

func invalidateCookie(res http.ResponseWriter) {
	expire := time.Now().AddDate(-1, 0, -1)
	cookie2 := http.Cookie{Name: "uuid", Path: "/", Expires: expire}
	http.SetCookie(res, &cookie2)
}

func logOut(res http.ResponseWriter, req *http.Request) {
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

func getNameAndCookie (res http.ResponseWriter, req *http.Request) (string, bool){
	var name string
	var ok bool
	var cookie, err = req.Cookie("uuid")
	correctlyLogIn := false
	// if the cookie is set up
	if err == nil {
		// retrive the name
		name, ok = sessions[cookie.Value]
		// if the name exist, print greetings
		if ok {
			correctlyLogIn = true
		// no name so invalidate cookie
		} else {
			invalidateCookie(res)
		}
	}
	return name, correctlyLogIn
}

// getTime is a function which display the time on the webside
func getTime(res http.ResponseWriter, req *http.Request) {
	now := time.Now().Format("3:04:05 PM")
	name, correctlyLogIn := getNameAndCookie(res, req)
	res.Header().Set("Content-Type", "text/html")
	if (correctlyLogIn){
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
	}else{
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

// unknownRoute is a function which display "These are not the URLs you're looking for"
// It also sets the status code to 404
func unknownRoute(res http.ResponseWriter, req *http.Request) {
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

	// parse the flags
	flag.IntVar(&port, "port", 8080, "used port")
	flag.BoolVar(&version, "V", false, "version of the program")
	flag.Parse()

	// if user type -V, the V flag is set up to true
	if version {
		// display the information about the version
		fmt.Println("version 1.4")
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