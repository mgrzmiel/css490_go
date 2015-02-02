// CSS 490
// Magdalena Grzmiel
// Assignments #3
// Copyright 2015 Magdalena Grzmiel
// cookiesManager is resposnisible for managing the cookies.

package cookiesManager

import (
	log "github.com/cihub/seelog"
	"net/http"
	"time"
)

// GetCookieValue
// Returns cookies values return the cookies value
func GetCookieValue(req *http.Request, cookieName string) (string, bool) {
	var cookie, err = req.Cookie(cookieName)
	if err != nil {
		log.Info("no cookie set up")
		return "", false
	}

	log.Debug("cookie" + cookie.Value)
	return cookie.Value, true
}

// SetCookieValue
// Sets the cookie with the provided name and value
func SetCookieValue(res http.ResponseWriter, name string, value string) {
	// save uuid in the cookie
	cookie := http.Cookie{Name: name, Value: value, Path: "/"}
	log.Debug("set follwoing cookie" + cookie.Value)
	http.SetCookie(res, &cookie)
}

// RemoveCookie
// Reset the cookie with setting the expirention date to one day back
func RemoveCookie(res http.ResponseWriter, name string) {
	// set the experiation date to last year
	expire := time.Now().AddDate(-1, 0, 0)
	log.Debug("new expiration date is " + expire.Format("3:04:05 PM"))
	//http.SetCookie(res, &cookie)
	cookie := http.Cookie{Name: name, Path: "/", Expires: expire}
	http.SetCookie(res, &cookie)
}
