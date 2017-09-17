package dinogen

import (
	"github.com/rob05c/sauropoda/db"
	"github.com/rob05c/sauropoda/dino"
	"github.com/rob05c/sauropoda/quadtree"
	"math/rand"
	"testing"
	"time"
)

// TODO deduplicate from quadtree_test.go?
func randStr(maxLen int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ !@#$%^&*()_+1234567890-=|}{:\">?<,./;'[]\\")
	n := rand.Intn(maxLen) + 1
	s := make([]rune, n)
	for i := 0; i < n; i++ {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

// TODO deduplicate from quadtree_test.go?
func randDino(lat, lon float64) quadtree.PositionedDinosaur {
	latRadius := MetresToLatitude(radiusMetres)
	lonRadius := MetresToLongitude(radiusMetres, lat)

	return quadtree.PositionedDinosaur{
		Dinosaur: dino.Dinosaur{
			Name:   randStr(100),
			Power:  int64(rand.Intn(500)),
			Health: int64(rand.Intn(100)),
		},
		Latitude:   lat - (latRadius / 2) + rand.Float64()*latRadius,
		Longitude:  lon - (lonRadius / 2) + rand.Float64()*lonRadius,
		Expiration: time.Now().Add(time.Second * time.Duration(rand.Intn(500))),
	}
}

func testSpecies() map[string]db.Species {
	return map[string]db.Species{
		"Brontosaurus": db.Species{
			Name:         "Brontosaurus",
			HeightMetres: 42,
			LengthMetres: 99,
			WeightKg:     5001,
			Popularity:   100,
		},
		"Apatasaurus": db.Species{
			Name:         "Apatasaurus",
			HeightMetres: 14,
			LengthMetres: 19,
			WeightKg:     5002,
			Popularity:   99,
		},
		"Doyouthinkeesaurus": db.Species{
			Name:         "Doyouthinkeesaurus",
			HeightMetres: 4,
			LengthMetres: 8,
			WeightKg:     16,
			Popularity:   32,
		},
	}

}

func TestGenerateInRadius(t *testing.T) {
	lat := 50.0
	lon := 50.0
	qt := quadtree.Create()
	for i := 0; i < dinosaursPerRadius-1; i++ {
		qt.Insert(randDino(lat, lon))
	}

	species := testSpecies()

	existing, generated := GenerateInRadius(qt, species, lat, lon)

	if len(existing) != dinosaursPerRadius-1 {
		t.Errorf("expected existing len: %v actual: %v", dinosaursPerRadius-1, len(existing))
	}

	if len(generated) != 1 {
		t.Errorf("expected generated len: 1 actual: %v", len(generated))
	}
}
