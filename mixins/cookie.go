package mixins

import (
	"net/http"
	"time"
)

func CreateCookie(name, token, path string) *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = token
	cookie.Expires = time.Now().Add(86400 * time.Second)
	cookie.HttpOnly = true
	cookie.Secure = true
	cookie.Path = path
	return cookie
}

func DeleteCookie(name string) *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = ""
	cookie.Expires = time.Unix(0, 0)
	return cookie
}
