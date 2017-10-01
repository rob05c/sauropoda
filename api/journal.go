package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/rob05c/sauropoda/dino"
	"github.com/rob05c/sauropoda/login"
)

func hndlJournal(d RouteData, w http.ResponseWriter, r *http.Request) {
	fmt.Print(time.Now().Format(time.RFC3339) + " INFO: " + r.RequestURI + " hdnlJournal\n")
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

	species, err := queryPlayerSpecies(d.DB, user)
	if err != nil {
		handleErr(http.StatusInternalServerError, "getting player speices: "+err.Error())
		return
	}

	missingSpecies, err := queryPlayerMissingSpecies(d.DB, user)
	if err != nil {
		handleErr(http.StatusInternalServerError, "getting player missing species: "+err.Error())
		return
	}

	// TODO replace this with a tree, when parents are implemented
	for _, s := range missingSpecies {
		species = append(species, dino.Species{Name: strings.Repeat("?", len(s))})
	}

	speciesJson, err := json.Marshal(species)
	if err != nil {
		handleErr(http.StatusInternalServerError, "marshalling player species: "+err.Error())
		return
	}
	w.Write(speciesJson)
}

func queryPlayerSpecies(db *sql.DB, player string) ([]dino.Species, error) {
	rows, err := db.Query("select name, height_m, length_m, weight_kg from species where name in (select distinct name from dinosaur where player = ?);", player)
	if err != nil {
		return nil, errors.New("error querying species: " + err.Error())
	}
	defer rows.Close()
	species := []dino.Species{}
	for rows.Next() {
		s := dino.Species{}
		if err := rows.Scan(&s.Name, &s.HeightMetres, &s.LengthMetres, &s.WeightKg); err != nil {
			return nil, errors.New("error scanning species: " + err.Error())
		}
		species = append(species, s)
	}
	return species, nil
}

func queryPlayerMissingSpecies(db *sql.DB, player string) ([]string, error) {
	// TODO add parent here, so a tree/cladogram can be constructed, when parents exist
	rows, err := db.Query("select name from species where name not in (select distinct name from dinosaur where player = ?);", player)
	if err != nil {
		return nil, errors.New("error querying species names: " + err.Error())
	}
	defer rows.Close()
	names := []string{}
	for rows.Next() {
		n := ""
		if err := rows.Scan(&n); err != nil {
			return nil, errors.New("error scanning species names: " + err.Error())
		}
		names = append(names, n)
	}
	return names, nil
}
