package hndl

import (
	"fmt"
	"github.com/rob05c/sauropoda/login"
	"net/http"
)

func handlePing(d RouteData, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	// TODO put in "refreshCookie" helper, for all handlers
	// TODO log errors
	if cookie, err := r.Cookie(login.CookieName); err == nil {
		if user, err := login.TokenValid(cookie.Value, d.DB, d.TokenKey); err == nil {
			if newToken, err := login.CreateToken(user, d.TokenKey); err == nil {
				http.SetCookie(w, login.CreateCookie(newToken))
			} else {
				fmt.Printf("Ping Error CreateToken: %v\n", err)
			}
		} else {
			fmt.Printf("Ping Error TokenValid: %v\n", err)
		}
	} else {
		fmt.Printf("Ping Error getting cookie: %v\n", err)
	}
	w.Write([]byte("pong\n"))
}
