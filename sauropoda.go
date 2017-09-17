package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/rob05c/sauropoda/hndl"
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

	rd := hndl.RouteData{
		DB:       db,
		Species:  species,
		QT:       quadtree.Create(),
		TokenKey: []byte(cfg.TokenKey),
	}

	//	fmt.Printf("Species: %v\n", species)
	fmt.Println("Serving :47777")
	Serve(rd)
}

func Serve(rd hndl.RouteData) {
	if err := hndl.RegisterHandlers(rd); err != nil {
		fmt.Printf("Error registering handlers: %v\n", err)
		return
	}

	if err := http.ListenAndServe(":47777", nil); err != nil {
		fmt.Printf("Error serving: %v\n", err)
		return
	}
}
