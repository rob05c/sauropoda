package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rob05c/sauropoda/login"
)

type PlayerInfo struct {
	Name string `json:"name"`
}

func hndlPlayer(d RouteData, w http.ResponseWriter, r *http.Request) {
	fmt.Print(time.Now().Format(time.RFC3339) + " INFO: " + r.RequestURI + " hndlDinos\n")
	defer r.Body.Close()
	// TODO put in "refreshCookie" helper, for all handlers
	// TODO log errors

	handleErr := func(code int, msg string) {
		w.WriteHeader(code)
		w.Write([]byte(http.StatusText(code)))
		fmt.Print(time.Now().Format(time.RFC3339) + " Error: " + r.RequestURI + " hndlDinos: " + msg + "\n")
	}

	cookie, err := r.Cookie(login.CookieName)
	if err != nil {
		handleErr(http.StatusUnauthorized, "getting cookie: "+err.Error())
		return
	}

	user, err := login.TokenValid(cookie.Value, d.DB, d.TokenKey)
	if err != nil {
		handleErr(http.StatusUnauthorized, "token invalid: "+err.Error())
		return
	}

	playerInfo := PlayerInfo{Name: user}

	playerJson, err := json.Marshal(playerInfo)
	if err != nil {
		handleErr(http.StatusInternalServerError, "marshalling playerInfo: "+err.Error())
		return
	}
	w.Write(playerJson) // TODO protect against injection, e.g. with `1;` or `"result":`
}
