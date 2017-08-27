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
	login "github.com/rob05c/sauropoda/login"
	"github.com/rob05c/sauropoda/quadtree"
	"github.com/rob05c/sauropoda/webui"
)

type RouteData struct {
	db       *sql.DB
	species  map[string]sdb.Species
	qt       quadtree.Quadtree
	tokenKey []byte
}

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

	cfg, err := sdb.LoadConfig(db)
	if err != nil {
		fmt.Printf("Error getting config from database: %v\n", err)
		return
	}

	rd := RouteData{
		db:       db,
		species:  species,
		qt:       quadtree.Create(),
		tokenKey: []byte(cfg.TokenKey),
	}

	//	fmt.Printf("Species: %v\n", species)
	fmt.Println("Serving :47777")
	serve(rd)
}

type DataHandlerFunc func(rd RouteData, w http.ResponseWriter, r *http.Request)

func wrapHandler(d RouteData, f DataHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f(d, w, r)
	}
}

// TODO change to return proper error codes
func handleQuery(d RouteData, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

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

	dinosaurs := dinogen.Query(d.qt, d.species, lat, lon) // []quadtree.PositionedDinosaur

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
	defer r.Body.Close()
	fmt.Fprintf(w, "%s", time.Now())
}

type LoginReq struct {
	User string `json:"user"`
	Pass string `json:"pass"`
}

func handleLogin(d RouteData, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var loginReq LoginReq
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		fmt.Printf("Login failed to decode JSON: %v\n", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	token, err := login.Login(d.db, d.tokenKey, loginReq.User, loginReq.Pass)
	if err != nil {
		fmt.Printf("Login error: %v\n", err)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	http.SetCookie(w, login.CreateCookie(token))
	w.WriteHeader(http.StatusOK) // TODO 204 No Content?
}

func handleCreateUser(d RouteData, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var loginReq LoginReq
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		fmt.Printf("Login failed to decode JSON: %v\n", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err := login.CreateUser(d.db, loginReq.User, loginReq.Pass)
	if err != nil {
		fmt.Printf("Error creating user: %v\n", err)
		// TODO confirm user exists (CreateUser could err for other reasons)
		http.Error(w, "user already exists", http.StatusBadRequest)
		return
	}
	token, err := login.Login(d.db, d.tokenKey, loginReq.User, loginReq.Pass)
	if err != nil {
		fmt.Printf("Login error: %v\n", err)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	http.SetCookie(w, login.CreateCookie(token))
	w.WriteHeader(http.StatusOK) // TODO 201 Created?
}

func handlePing(d RouteData, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	// TODO put in "refreshCookie" helper, for all handlers
	// TODO log errors
	if cookie, err := r.Cookie(login.CookieName); err == nil {
		if user, err := login.TokenValid(cookie.Value, d.db, d.tokenKey); err == nil {
			if newToken, err := login.CreateToken(user, d.tokenKey); err == nil {
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

// TODO create api.RegisterHandlers
func registerHandlers(rd RouteData) error {
	uiPathPrefix := ""
	if err := webui.RegisterHandlers(http.DefaultServeMux, uiPathPrefix, rd.species); err != nil {
		return err
	}
	http.HandleFunc("/query/", wrapHandler(rd, handleQuery))
	http.HandleFunc("/now", handleNow)
	http.HandleFunc("/login", wrapHandler(rd, handleLogin))
	http.HandleFunc("/createuser", wrapHandler(rd, handleCreateUser))
	http.HandleFunc("/ping", wrapHandler(rd, handlePing))
	return nil
}

func serve(rd RouteData) {
	if err := registerHandlers(rd); err != nil {
		fmt.Printf("Error registering handlers: %v\n", err)
		return
	}

	if err := http.ListenAndServe(":47777", nil); err != nil {
		fmt.Printf("Error serving: %v\n", err)
		return
	}
}
