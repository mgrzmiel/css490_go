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
		log.Warn("Logging out without the cookie")
	}
}

// Login
// It logins by extrcating the name from request, generating the value for cookie and
// setting the cookie
func (cbs *CookieBasedSessions) Login(res http.ResponseWriter, req *http.Request) bool {
	name := req.FormValue(NameParameter)
	name = html.EscapeString(name)
	log.Info("got the name from request " + name)
	if name != "" {
		// generate uuid
		uuid := cbs.SessionManager.CreateSession(name)
		cookiesManager.SetCookieValue(res, CookieName, uuid)

		// successfully loged in
		return true
	}
	return false
}

// GetSession
// It checks if there is an user correctly log in - the cookie is set up and the name of the user
// exists in map. It returns the right bool value
func (cbs *CookieBasedSessions) GetSession(res http.ResponseWriter, req *http.Request) (string, bool) {
	var name string
	var founded bool
	correctlyLogIn := false
	uuid, ok := cookiesManager.GetCookieValue(req, CookieName)
	if ok {
		log.Debug("Cookie is set up and it is follwoing" + uuid)
		name, founded = cbs.SessionManager.GetSession(uuid)
		if founded {
			log.Debug("the name is set up and is follwoing" + name)
			correctlyLogIn = true
		} else {
			log.Warn("but no name set up")
			cookiesManager.RemoveCookie(res, CookieName)
		}
	}
	log.Debug("No cookie sets up")
	return name, correctlyLogIn
}
