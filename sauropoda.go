package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/rob05c/sauropoda/api"
	"github.com/rob05c/sauropoda/quadtree"
	"github.com/rob05c/sauropoda/sdb"
)

func main() {
	rand.Seed(time.Now().Unix())

	var port = flag.Int("port", -1, "HTTP port; ignored if https-port is nonnegative")
	var crtPath = flag.String("cert-path", "", "HTTPS certificate path")
	var crtKeyPath = flag.String("cert-key-path", "", "HTTPS certificate key path")
	var dbPath = flag.String("db", "./db.sqlite", "Database file path")
	flag.Parse()
	if *port == -1 {
		if *crtPath == "" {
			*port = 80
		} else {
			*port = 443
		}
	}

	db, err := sdb.Create(*dbPath)
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

	rd := api.RouteData{
		DB:       db,
		Species:  species,
		QT:       quadtree.Create(),
		TokenKey: []byte(cfg.TokenKey),
	}

	if err := api.RegisterHandlers(rd); err != nil {
		fmt.Printf("Error registering handlers: %v\n", err)
		return
	}

	//	fmt.Printf("Species: %v\n", species)
	if *crtPath == "" {
		fmt.Println("Serving HTTP: " + strconv.Itoa(*port))
		if err := http.ListenAndServe(":"+strconv.Itoa(*port), nil); err != nil {
			fmt.Printf("Error serving: %v\n", err)
			return
		}
	} else {
		fmt.Println("Redirecting HTTP 80, Serving HTTPS: " + strconv.Itoa(*port))
		go http.ListenAndServe(":80", http.HandlerFunc(redirectHTTP))
		if err := http.ListenAndServeTLS(":"+strconv.Itoa(*port), *crtPath, *crtKeyPath, nil); err != nil {
			fmt.Printf("Error serving: %v\n", err)
			return
		}
	}
}

func redirectHTTP(w http.ResponseWriter, req *http.Request) {
	target := "https://" + req.Host + req.URL.Path
	if req.URL.RawQuery != "" {
		target += "?" + req.URL.RawQuery
	}
	http.Redirect(w, req, target, http.StatusTemporaryRedirect)
}
