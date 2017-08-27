package quadtree

import (
	"github.com/rob05c/sauropoda/dinosaur"
	"math/rand"
	"testing"
	"time"
)

func randStr(maxLen int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ !@#$%^&*()_+1234567890-=|}{:\">?<,./;'[]\\")
	n := rand.Intn(maxLen) + 1
	s := make([]rune, n)
	for i := 0; i < n; i++ {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func randDino() PositionedDinosaur {
	return PositionedDinosaur{
		Dinosaur: dinosaur.Dinosaur{
			Name:   randStr(100),
			Power:  int64(rand.Intn(500)),
			Health: int64(rand.Intn(100)),
		},
		Latitude:   rand.Float64()*180 - 90,
		Longitude:  rand.Float64()*360 - 180,
		Expiration: time.Now().Add(time.Second * time.Duration(rand.Intn(500))),
	}
}

func TestCreate(t *testing.T) {
	qt := Create()
	if qt.root.left != -180.0 || qt.root.right != 180.0 || qt.root.top != 90.0 || qt.root.bottom != -90.0 {
		t.Errorf("quadtree dimensions invalid")
	}
	if qt.root.topLeft != nil || qt.root.topRight != nil || qt.root.bottomLeft != nil || qt.root.bottomRight != nil {
		t.Errorf("new quadtree children not empty")
	}
	if len(qt.root.elements) != 0 {
		t.Errorf("new quadtree elements not empty")
	}
}

func TestInsert(t *testing.T) {
	qt := Create()
	d := randDino()
	qt.Insert(d)
	if len(qt.root.elements) != 1 {
		t.Fatalf("qt.root.elements len expected: 1 actual: %d", len(qt.root.elements))
	}
	insertedD := qt.root.elements[0]
	if insertedD != d {
		t.Errorf("qt.root.element[0] expected: %v actual: %v", d, insertedD)
	}
}

func TestSplit(t *testing.T) {
	qt := Create()
	for i := 0; i < elementSize+1; i++ {
		qt.Insert(randDino())
	}

	if len(qt.root.elements) > 0 {
		t.Errorf("len qt.root.element expected: 0 actual: %v", len(qt.root.elements))
	}

	if qt.root.topLeft == nil || qt.root.topRight == nil || qt.root.bottomLeft == nil || qt.root.bottomRight == nil {
		t.Errorf("len qt.root.element children expected: not nil actual: nil")
	}

	if childrenLens := len(qt.root.topLeft.elements) + len(qt.root.topRight.elements) + len(qt.root.bottomLeft.elements) + len(qt.root.bottomRight.elements); childrenLens != elementSize+1 {
		t.Errorf("len qt.root.element children len expected: %d actual: %d with root len %v", elementSize+1, childrenLens, len(qt.root.elements))
	}
}

func TestGet(t *testing.T) {
	qt := Create()
	for i := 0; i < elementSize*50+1; i++ {
		qt.Insert(randDino())
	}

	d := PositionedDinosaur{
		Dinosaur: dinosaur.Dinosaur{
			Name:   "George",
			Power:  42,
			Health: 97,
		},
		Latitude:   19,
		Longitude:  -72,
		Expiration: time.Now().Add(time.Second * 60),
	}

	qt.Insert(d)

	dinos := qt.Get(20, -73, 18, -71)

	foundI := -1
	for i, dino := range dinos {
		if dino.Dinosaur.Name == "George" {
			foundI = i
			break
		}
	}
	if foundI == -1 {
		t.Errorf("len qt.root.element Get() expected: d actual: no d")
	}

}

func TestRemoveExpired(t *testing.T) {
	qt := Create()
	d := PositionedDinosaur{
		Dinosaur: dinosaur.Dinosaur{
			Name:   "George",
			Power:  42,
			Health: 97,
		},
		Latitude:   19,
		Longitude:  -72,
		Expiration: time.Now().Add(time.Second * -1),
	}
	qt.Insert(d)
	dinos := qt.Get(20, -73, 18, -71)
	if len(dinos) > 0 {
		t.Errorf("Get len expected: 0 actual: %v", len(dinos))
	}

	qt = Create()
	d = PositionedDinosaur{
		Dinosaur: dinosaur.Dinosaur{
			Name:   "George",
			Power:  42,
			Health: 97,
		},
		Latitude:   19,
		Longitude:  -72,
		Expiration: time.Now().Add(time.Second * 1),
	}
	qt.Insert(d)
	dinos = qt.Get(20, -73, 18, -71)
	if len(dinos) != 1 {
		t.Errorf("Get len expected: 1 actual: %v", len(dinos))
	}

}
