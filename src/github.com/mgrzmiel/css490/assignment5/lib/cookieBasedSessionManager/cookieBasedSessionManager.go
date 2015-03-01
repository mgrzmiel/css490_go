// CSS 490
// Magdalena Grzmiel
// Assignments #5
// Copyright 2015 Magdalena Grzmiel
// cookieBasedSessionManager is resposnisible for managing
// cookies and session.

package cookieBasedSessionManager

import (
	"bytes"
	log "github.com/cihub/seelog"
	"github.com/mgrzmiel/css490/assignment5/lib/authClient"
	"github.com/mgrzmiel/css490/assignment5/lib/cookiesManager"
	"html"
	"net/http"
	"os/exec"
	"strings"
)

const (
	CookieName    = "uuid"
	NameParameter = "name"
)

// Logout
// It logouts by removing the session from map and invalidating the cookie.
func Logout(res http.ResponseWriter, req *http.Request) {
	_, ok := cookiesManager.GetCookieValue(req, CookieName)
	if ok {
		// cbs.SessionManager.RemoveSession(uuid)
		cookiesManager.RemoveCookie(res, CookieName)
	} else {
		log.Trace("Logging out without the cookie")
	}
}

// Login
// It logins by extracting the name from request, creating the session for it and store
// session key in the cookie
func Login(res http.ResponseWriter, req *http.Request) (bool, string) {
	name := req.FormValue(NameParameter)
	name = html.EscapeString(name)
	log.Debugf("Log in user. Name: %s", name)
	if name != "" {
		uuid := generateRandomUUID()
		success := authClient.SetRequest(uuid, name)
		if success {
			cookiesManager.SetCookieValue(res, CookieName, uuid)
		}
		// successfully loged in
		if success {
			return success, ""
		} else {
			return success, "authServerFail"
		}

	}

	return false, "noName"
}

// generateRandomUUID
// This functions uses linux uuidgen to generate random UUID for
// session's needs
func generateRandomUUID() string {
	//generate uuidgen
	cmd := exec.Command("/usr/bin/uuidgen")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Errorf("Not able to generate uuid", err)
	}
	uuid := out.String()
	return strings.Replace(uuid, "\n", "", 1)
}

// GetSession
// It checks if there is an user correctly log in - the cookie is set up and the name of the user
// exists in sessions. It returns the right bool value
func GetSession(res http.ResponseWriter, req *http.Request) (string, bool) {
	var name string
	var success bool
	correctlyLogIn := false
	uuid, ok := cookiesManager.GetCookieValue(req, CookieName)
	if ok {
		log.Debugf("Found cookie. Value: %s", uuid)
		name, success = authClient.GetRequest(uuid)
		if name != "" && success {
			log.Debugf("Found session for key: %s with value: %s", uuid, name)
			correctlyLogIn = true
		} else {
			log.Debugf("No session found for key; %s", uuid)
			cookiesManager.RemoveCookie(res, CookieName)
		}
	}

	log.Debug("Cookie is missing")
	return name, correctlyLogIn
}
