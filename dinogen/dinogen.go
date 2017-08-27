package dinogen

import (
	"github.com/rob05c/sauropoda/db"
	"github.com/rob05c/sauropoda/dinosaur"
	"github.com/rob05c/sauropoda/quadtree"
	"math"
	"math/rand"
	"sync/atomic"
	"time"
)

const radiusMetres = 500
const dinosaursPerRadius = 5

const metresToLatitude = 110574
const metresToLongitudeTimesCosLat = 111320

func DegreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

func MetresToLatitude(metres float64) float64 {
	return metres / metresToLatitude
}

func MetresToLongitude(metres float64, latitude float64) float64 {
	return metres / (metresToLongitudeTimesCosLat * math.Cos(DegreesToRadians(latitude)))
}

// returns a slice of the existing dinosaurs, and a slice of new generated dinosaurs
func GenerateInRadius(qt quadtree.Quadtree, species map[string]db.Species, lat float64, lon float64) (existing []quadtree.PositionedDinosaur, generated []quadtree.PositionedDinosaur) {
	latRadius := MetresToLatitude(radiusMetres) // TODO precompute?
	lonRadius := MetresToLongitude(radiusMetres, lat)
	top := lat + latRadius
	bottom := lat - latRadius
	left := lon - lonRadius
	right := lon + lonRadius

	existingDinosaurs := qt.Get(top, left, bottom, right)

	var generatedDinosaurs []quadtree.PositionedDinosaur
	generateCount := dinosaursPerRadius - len(existingDinosaurs)
	for i := 0; i < generateCount; i++ {
		generatedDinosaurs = append(generatedDinosaurs, specieToPositioned(dinosaur.Generate(species), top, left, bottom, right))
	}
	return existingDinosaurs, generatedDinosaurs
	//	newSpecie := dinosaur.Generate(species)
}

// TODO change to get max id from database
var id uint64

// specieToIndividual takes the species to create an individual of,
// and the range to randomly create it in.
// TODO move to within dinosaur.Generate()?
func specieToPositioned(specie db.Species, top, left, bottom, right float64) quadtree.PositionedDinosaur {
	return quadtree.PositionedDinosaur{
		Dinosaur: dinosaur.Dinosaur{
			Name:   specie.Name,
			Power:  int64(rand.Intn(100)), // TODO implement
			Health: 100,                   // TODO implement
		},
		Latitude:   rand.Float64()*(top-bottom) + bottom,
		Longitude:  rand.Float64()*(right-left) + left,
		Expiration: time.Now().Add(time.Second * time.Duration(rand.Intn(int(specie.Popularity*3)+60))), // TODO formalise popularity/expiration ratio
		ID:         nextDinosaurID(),
	}
}

// TODO rename (removed 'owned' - also used by positioned
// TODO change to get max from database on startup
var nextOwnedDinosaurID uint64

func nextDinosaurID() int64 {
	return int64(atomic.AddUint64(&nextOwnedDinosaurID, 1))
}

func positionedToOwned(p quadtree.PositionedDinosaur) dinosaur.OwnedDinosaur {
	return dinosaur.OwnedDinosaur{Dinosaur: p.Dinosaur, ID: nextDinosaurID()}
}

func Query(qt quadtree.Quadtree, species map[string]db.Species, lat float64, lon float64) []quadtree.PositionedDinosaur {
	existing, generated := GenerateInRadius(qt, species, lat, lon)
	for _, d := range generated {
		qt.Insert(d)
		existing = append(existing, d)
	}
	return existing
}
