package hndl

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/rob05c/sauropoda/dino"
	"github.com/rob05c/sauropoda/dinogen"
	"github.com/rob05c/sauropoda/login"
)

func hndlCatch(d RouteData, w http.ResponseWriter, r *http.Request) {
	fmt.Print(time.Now().Format(time.RFC3339) + " INFO: " + r.RequestURI + " hndlCatch\n")
	defer r.Body.Close()
	// TODO put in "refreshCookie" helper, for all handlers
	// TODO log errors

	handleErr := func(code int, msg string) {
		w.WriteHeader(code)
		w.Write([]byte(http.StatusText(code)))
		fmt.Print(time.Now().Format(time.RFC3339) + " Error: " + r.RequestURI + " hndlCatch: " + msg + "\n")
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

	idStrs, ok := r.URL.Query()["id"]
	if !ok {
		handleErr(http.StatusBadRequest, "no ID")
		return
	}
	if len(idStrs) != 1 {
		handleErr(http.StatusBadRequest, "multiple IDs")
		return
	}
	idStr := idStrs[0]

	id, err := strconv.ParseInt(idStr, 10, 64)
	if !ok {
		handleErr(http.StatusBadRequest, "ID not an int: "+err.Error())
		return
	}

	dino, ok := d.QT.GetByID(id)
	if !ok {
		handleErr(http.StatusNotFound, "ID has no dino")
		return
	}

	ownedDino := dinogen.PositionedToOwned(*dino)
	if err := insertOwnedDino(d.DB, user, ownedDino); err != nil {
		handleErr(http.StatusBadRequest, "failed to insert dino: "+err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func insertOwnedDino(db *sql.DB, player string, d dino.OwnedDinosaur) error {
	if _, err := db.Exec("insert into dinosaur (id, player, positioned_id, latitude, longitude, catch_time, name, power, health) values (?, ?, ?, ?, ?, ?, ?, ?, ?);", d.ID, player, d.PositionedDinosaur.ID, d.Latitude, d.Longitude, d.Expiration, d.Name, d.Power, d.Health); err != nil {
		// TODO return constant for already owned dino
		return errors.New("error inserting dinosaur: " + err.Error())
	}
	return nil
}
