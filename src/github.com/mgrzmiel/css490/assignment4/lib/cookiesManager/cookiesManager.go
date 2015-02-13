// CSS 490
// Magdalena Grzmiel
// Assignments #4
// Copyright 2015 Magdalena Grzmiel
// cookiesManager is resposnisible for managing the cookies.

package cookiesManager

import (
	log "github.com/cihub/seelog"
	"net/http"
	"time"
)

// GetCookieValue
// Returns cookie's value with the specified name
func GetCookieValue(req *http.Request, cookieName string) (string, bool) {
	var cookie, err = req.Cookie(cookieName)
	if err != nil {
		log.Debugf("No cookie with name %s found", cookieName)
		return "", false
	}

	log.Debugf("Found cookie with name %s and value %s", cookieName, cookie.Value)
	return cookie.Value, true
}

// SetCookieValue
// Sets the cookie with the provided name and value
func SetCookieValue(res http.ResponseWriter, name string, value string) {
	// save uuid in the cookie
	cookie := http.Cookie{Name: name, Value: value, Path: "/"}
	log.Debugf("Setting cookie with name %s and value %s", name, value)
	http.SetCookie(res, &cookie)
}

// RemoveCookie
// Reset the cookie with setting the expirention date to one day back
func RemoveCookie(res http.ResponseWriter, name string) {
	// set the experiation date to last year
	expire := time.Now().AddDate(-1, 0, 0)
	log.Debugf("Expirering cookie with name %s", name)
	cookie := http.Cookie{Name: name, Path: "/", Expires: expire}
	http.SetCookie(res, &cookie)
}
