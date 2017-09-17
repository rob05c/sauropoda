package hndl

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rob05c/sauropoda/login"
)

func handleLogin(d RouteData, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var loginReq LoginReq
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		fmt.Printf("Login failed to decode JSON: %v\n", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	token, err := login.Login(d.DB, d.TokenKey, loginReq.User, loginReq.Pass)
	if err != nil {
		fmt.Printf("Login error: %v\n", err)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	http.SetCookie(w, login.CreateCookie(token))
	w.WriteHeader(http.StatusOK) // TODO 204 No Content?
}
