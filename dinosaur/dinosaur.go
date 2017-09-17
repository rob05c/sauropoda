package dinosaur

import (
	"math/rand"
	"time"
)

type Species struct {
	Name         string
	HeightMetres float64
	LengthMetres float64
	WeightKg     float64
	Popularity   int64
}

type Dinosaur struct {
	Name   string
	Power  int64
	Health int64
}

// TODO add json lowercase tags
type PositionedDinosaur struct {
	Dinosaur
	Latitude   float64
	Longitude  float64
	Expiration time.Time // TODO rename, as this is used for "caught" time too
	ID         int64
}

type OwnedDinosaur struct {
	PositionedDinosaur
	ID int64
}

// Generator takes the species map, and returns an array of names, which when randomly indexed will result in the appropriate distribution
// TODO add location info, e.g. sauropterygia should generate near water.
func Generator(species map[string]Species) []string {
	// TODO make this more efficient
	names := []string{}
	for name, specie := range species {
		for i := 0; i < int(specie.Popularity); i++ {
			names = append(names, name)
		}
	}
	return names
}

// TODO change to take a generator, so it isn't created every time
// TODO change to keep a *Rand, so it isn't created every time
func Generate(species map[string]Species) Species {
	generator := Generator(species)
	i := rand.Intn(len(generator))
	return species[generator[i]]
}
