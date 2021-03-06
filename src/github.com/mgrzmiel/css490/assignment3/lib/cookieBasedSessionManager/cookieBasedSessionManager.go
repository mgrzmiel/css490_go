// CSS 490
// Magdalena Grzmiel
// Assignments #3
// Copyright 2015 Magdalena Grzmiel
// cookieBasedSessionManager is resposnisible for managing
// cookies and session.

package cookieBasedSessionManager

import (
	log "github.com/cihub/seelog"
	"github.com/mgrzmiel/css490/assignment3/lib/cookiesManager"
	"github.com/mgrzmiel/css490/assignment3/lib/sessionManager"
	"html"
	"net/http"
)

const (
	CookieName    = "uuid"
	NameParameter = "name"
)

type CookieBasedSessions struct {
	SessionManager *sessionManager.Sessions
}

// CookieBasedSessions
// Return new CookieBasedSession structure
func New() *CookieBasedSessions {
	return &CookieBasedSessions{SessionManager: sessionManager.New()}
}

// Logout
// It logouts by remopving the session from map and invalidating the cookie.
func (cbs *CookieBasedSessions) Logout(res http.ResponseWriter, req *http.Request) {
	uuid, ok := cookiesManager.GetCookieValue(req, CookieName)
	if ok {
		cbs.SessionManager.RemoveSession(uuid)
		cookiesManager.RemoveCookie(res, CookieName)
	} else {
		log.Trace("Logging out without the cookie")
	}
}

// Login
// It logins by extracting the name from request, creating the session for it and store
// session key in the cookie
func (cbs *CookieBasedSessions) Login(res http.ResponseWriter, req *http.Request) bool {
	name := req.FormValue(NameParameter)
	name = html.EscapeString(name)
	log.Debugf("Log in user. Name: %s", name)
	if name != "" {
		// generate session
		uuid := cbs.SessionManager.CreateSession(name)
		cookiesManager.SetCookieValue(res, CookieName, uuid)

		// successfully loged in
		return true
	}

	return false
}

// GetSession
// It checks if there is an user correctly log in - the cookie is set up and the name of the user
// exists in sessions. It returns the right bool value
func (cbs *CookieBasedSessions) GetSession(res http.ResponseWriter, req *http.Request) (string, bool) {
	var name string
	var founded bool
	correctlyLogIn := false
	uuid, ok := cookiesManager.GetCookieValue(req, CookieName)
	if ok {
		log.Debugf("Found cookie. Value: %s", uuid)
		name, founded = cbs.SessionManager.GetSession(uuid)
		if founded {
			log.Debugf("Found session for kye: %s with value: %s", uuid, name)
			correctlyLogIn = true
		} else {
			log.Debugf("No session found for key; %s", uuid)
			cookiesManager.RemoveCookie(res, CookieName)
		}
	}

	log.Debug("Cookie is missing")
	return name, correctlyLogIn
}
