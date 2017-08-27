package main

import (
	"fmt"
	//	"io/ioutil"
	"strconv"
	//	"html"
	"database/sql"
	"encoding/json"
	"math/rand"
	"net/http"
	"strings"
	"time"

	sdb "github.com/rob05c/sauropoda/db"
	"github.com/rob05c/sauropoda/dinogen"
	"github.com/rob05c/sauropoda/quadtree"
	"github.com/rob05c/sauropoda/webui"
)

func main() {
	rand.Seed(time.Now().Unix())

	db, err := sdb.Create()
	if err != nil {
		fmt.Printf("Error creating database: %v\n", err)
		return
	}

	species, err := sdb.LoadSpecies(db)
	if err != nil {
		fmt.Printf("Error getting species from database: %v\n", err)
		return
	}

	qt := quadtree.Create()

	//	fmt.Printf("Species: %v\n", species)
	fmt.Println("Serving :47777")
	serve(db, species, qt)
}

type DataHandlerFunc func(db *sql.DB, species map[string]sdb.Species, qt quadtree.Quadtree, w http.ResponseWriter, r *http.Request)

func wrapHandler(db *sql.DB, species map[string]sdb.Species, qt quadtree.Quadtree, f DataHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f(db, species, qt, w, r)
	}
}

// TODO change to return proper error codes
func handleQuery(db *sql.DB, species map[string]sdb.Species, qt quadtree.Quadtree, w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%v %v X-Real-IP %v X-Forwarded-For %v requested %v\n", time.Now(), r.RemoteAddr, r.Header.Get("X-Real-IP"), r.Header.Get("X-Forwarded-For"), r.URL)
	urlParts := strings.Split(r.URL.Path, "/")
	if len(urlParts) < 4 {
		fmt.Fprintf(w, "Error: Not enough parts")
		return
	}
	latStr := urlParts[2]
	lonStr := urlParts[3]
	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		fmt.Fprintf(w, "Error: latitude not a number")
		return
	}
	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		fmt.Fprintf(w, "Error: longitude not a number")
		return
	}

	if lat > 90.0 || lat < -90.0 {
		fmt.Fprintf(w, "Error: latitude not between 90 and -90")
		return
	}

	if lon > 180.0 || lat < -180.0 {
		fmt.Fprintf(w, "Error: longitude not between 180 and -180")
		return
	}

	dinosaurs := dinogen.Query(qt, species, lat, lon) // []quadtree.PositionedDinosaur

	dinosaursJson, err := json.Marshal(dinosaurs)
	if err != nil {
		fmt.Printf("Error marshalling dinosaurs: %v", err)
		fmt.Fprintf(w, "Internal Server Error")
		return
	}

	fmt.Fprintf(w, "%s", string(dinosaursJson))
}

// handleNow handles a request for the current server time
func handleNow(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", time.Now())
}

// TODO create api.RegisterHandlers
func registerHandlers(db *sql.DB, species map[string]sdb.Species, qt quadtree.Quadtree) error {
	uiPathPrefix := ""
	if err := webui.RegisterHandlers(http.DefaultServeMux, uiPathPrefix, species); err != nil {
		return err
	}
	http.HandleFunc("/query/", wrapHandler(db, species, qt, handleQuery))
	http.HandleFunc("/now", handleNow)
	return nil
}

func serve(db *sql.DB, species map[string]sdb.Species, qt quadtree.Quadtree) {
	if err := registerHandlers(db, species, qt); err != nil {
		fmt.Printf("Error registering handlers: %v\n", err)
		return
	}

	if err := http.ListenAndServe(":47777", nil); err != nil {
		fmt.Printf("Error serving: %v\n", err)
		return
	}
}
