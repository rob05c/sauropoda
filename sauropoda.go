package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/rob05c/sauropoda/api"
	"github.com/rob05c/sauropoda/quadtree"
	"github.com/rob05c/sauropoda/sdb"
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
	fmt.Println("Serving :80")
	Serve(rd)
}

func Serve(rd api.RouteData) {
	if err := api.RegisterHandlers(rd); err != nil {
		fmt.Printf("Error registering handlers: %v\n", err)
		return
	}

	if err := http.ListenAndServe(":80", nil); err != nil {
		fmt.Printf("Error serving: %v\n", err)
		return
	}
}
