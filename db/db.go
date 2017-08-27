package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Species struct {
	Name         string
	HeightMetres float64
	LengthMetres float64
	WeightKg     float64
	Popularity   int64
}

func Create() (*sql.DB, error) {
	return sql.Open("sqlite3", "./db.sqlite")
}

func LoadSpecies(db *sql.DB) (map[string]Species, error) {
	rows, err := db.Query("select name, height_m, length_m, weight_kg, popularity from species;")
	if err != nil {
		return nil, err
	}

	species := map[string]Species{}
	for rows.Next() {
		s := Species{}
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
