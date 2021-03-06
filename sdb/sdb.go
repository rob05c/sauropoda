package sdb

import (
	"database/sql"
	"fmt"

	"github.com/rob05c/sauropoda/dino"

	_ "github.com/mattn/go-sqlite3"
)

func Create(path string) (*sql.DB, error) {
	return sql.Open("sqlite3", path)
}

func LoadSpecies(db *sql.DB) (map[string]dino.Species, error) {
	rows, err := db.Query("select name, height_m, length_m, weight_kg, popularity from species;")
	if err != nil {
		return nil, err
	}

	species := map[string]dino.Species{}
	for rows.Next() {
		s := dino.Species{}
		err := rows.Scan(&s.Name, &s.HeightMetres, &s.LengthMetres, &s.WeightKg, &s.Popularity)
		if err != nil {
			return nil, err
		}
		species[s.Name] = s
	}
	return species, nil
}

type DBConfig struct {
	TokenKey string
}

func LoadConfig(db *sql.DB) (*DBConfig, error) {
	cfg := &DBConfig{}
	if err := db.QueryRow("select token_key from config").Scan(&cfg.TokenKey); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("database missing config row")
		}
		return nil, fmt.Errorf("querying: %v", err)
	}
	return cfg, nil
}
