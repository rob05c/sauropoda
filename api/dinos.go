package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/rob05c/sauropoda/dino"
	"github.com/rob05c/sauropoda/login"
)

func hndlDinos(d RouteData, w http.ResponseWriter, r *http.Request) {
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

	dinos, err := queryOwnedDino(d.DB, user)
	if err != nil {
		handleErr(http.StatusInternalServerError, "getting dinos: "+err.Error())
		return
	}

	dinosJson, err := json.Marshal(dinos)
	if err != nil {
		handleErr(http.StatusInternalServerError, "marshalling dinos: "+err.Error())
		return
	}
	w.Write(dinosJson)
}

func queryOwnedDino(db *sql.DB, player string) ([]dino.OwnedDinosaur, error) {
	rows, err := db.Query("select id, positioned_id, latitude, longitude, catch_time, name, power, health from dinosaur where player = ?;", player)
	if err != nil {
		return nil, errors.New("error querying dinosaurs: " + err.Error())
	}
	defer rows.Close()
	dinos := []dino.OwnedDinosaur{}
	for rows.Next() {
		d := dino.OwnedDinosaur{}
		expiration := int64(0)
		if err := rows.Scan(&d.ID, &d.PositionedID, &d.Latitude, &d.Longitude, &expiration, &d.Name, &d.Power, &d.Health); err != nil {
			return nil, errors.New("error scanning dinosaurs: " + err.Error())
		}
		d.Expiration = time.Unix(0, expiration)
		dinos = append(dinos, d)
	}
	return dinos, nil
}
