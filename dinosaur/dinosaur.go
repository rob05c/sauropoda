package dinosaur

import (
	"github.com/rob05c/sauropoda/db"
	"math/rand"
)

type Dinosaur struct {
	Name   string
	Power  int64
	Health int64
}

type OwnedDinosaur struct {
	Dinosaur
	ID int64
}

// Generator takes the species map, and returns an array of names, which when randomly indexed will result in the appropriate distribution
// TODO add location info, e.g. sauropterygia should generate near water.
func Generator(species map[string]db.Species) []string {
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
func Generate(species map[string]db.Species) db.Species {
	generator := Generator(species)
	i := rand.Intn(len(generator))
	return species[generator[i]]
}
