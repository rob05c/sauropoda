package hndl

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rob05c/sauropoda/login"
)

type LoginReq struct {
	User string `json:"user"`
	Pass string `json:"pass"`
}

func handleCreateUser(d RouteData, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var loginReq LoginReq
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		fmt.Printf("Login failed to decode JSON: %v\n", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err := login.CreateUser(d.DB, loginReq.User, loginReq.Pass)
	if err != nil {
		fmt.Printf("Error creating user: %v\n", err)
		// TODO confirm user exists (CreateUser could err for other reasons)
		http.Error(w, "user already exists", http.StatusBadRequest)
		return
	}
	token, err := login.Login(d.DB, d.TokenKey, loginReq.User, loginReq.Pass)
	if err != nil {
		fmt.Printf("Login error: %v\n", err)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	http.SetCookie(w, login.CreateCookie(token))
	w.WriteHeader(http.StatusOK) // TODO 201 Created?
}
