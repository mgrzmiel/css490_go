package cookiesManager

import (
	"net/http"
	"time"
)

func GetCookieValue(req *http.Request, cookieName string) (string, bool) {
	var cookie, err = req.Cookie(cookieName)
	if err != nil {
		return "", false
	}

	return cookie.Value, true
}

func SetCookieValue(res http.ResponseWriter, name string, value string) {
	// save uuid in the cookie
	cookie := http.Cookie{Name: name, Value: value, Path: "/"}
	http.SetCookie(res, &cookie)
}

// invalidate cookie
// It invalidates cookies since no name exists for that uuid in map
func RemoveCookie(res http.ResponseWriter, name string) {
	// set the experiation date to last year
	expire := time.Now().AddDate(-1, 0, 0)
	cookie := http.Cookie{Name: name, Path: "/", Expires: expire}
	http.SetCookie(res, &cookie)
}
