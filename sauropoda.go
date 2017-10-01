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

	var port = flag.Int("port", 80, "HTTP port")
	var dbPath = flag.String("db", "./db.sqlite", "Database file path")
	flag.Parse()

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

	//	fmt.Printf("Species: %v\n", species)
	fmt.Println("Serving :" + strconv.Itoa(*port))
	Serve(*port, rd)
}

func Serve(port int, rd api.RouteData) {
	if err := api.RegisterHandlers(rd); err != nil {
		fmt.Printf("Error registering handlers: %v\n", err)
		return
	}

	if err := http.ListenAndServe(":"+strconv.Itoa(port), nil); err != nil {
		fmt.Printf("Error serving: %v\n", err)
		return
	}
}
