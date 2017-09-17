package quadtree

import (
	"github.com/rob05c/sauropoda/dino"
	"math/rand"
	"time"
)

const elementSize = 100

type Quadtree struct {
	root  *Node
	dinos map[int64]*dino.PositionedDinosaur
}

func (q *Quadtree) GetByID(id int64) (*dino.PositionedDinosaur, bool) {
	dino, ok := q.dinos[id]
	if !ok {
		return nil, false
	}
	if time.Now().After(dino.Expiration) {
		return nil, false
	}
	return dino, true
}

type Node struct {
	left        float64
	right       float64
	top         float64
	bottom      float64
	topLeft     *Node
	topRight    *Node
	bottomLeft  *Node
	bottomRight *Node
	elements    []dino.PositionedDinosaur
}

func (n *Node) hasSplit() bool {
	return n.topLeft != nil // it has split, if any child nodes are non-nil
}

func init() {
	rand.Seed(time.Now().Unix())
}

// Create Returns a Latitude-Longitude Quadtree
// Note this creates a latlon quadtree. In Lat Lon, right > left and top > bottom
func Create() Quadtree {
	return createGeneric(-180.0, 180.0, 90.0, -90.0)
}

func createGeneric(l, r, t, b float64) Quadtree {
	return Quadtree{
		root: &Node{
			left:     l,
			right:    r,
			top:      t,
			bottom:   b,
			elements: make([]dino.PositionedDinosaur, 0, elementSize),
		},
		dinos: map[int64]*dino.PositionedDinosaur{},
	}
}

// Insert inserts the given dinosaur into the quadtree
func (q *Quadtree) Insert(d dino.PositionedDinosaur) {
	q.root.insert(d)
	q.dinos[d.PositionedID] = &d
}

func (n *Node) insert(d dino.PositionedDinosaur) {
	if !n.hasSplit() {
		n.elements = append(n.elements, d)
		n.splitIfNecessary()
		return
	}

	n.insertAppropriateChild(d)
}

func (n *Node) insertAppropriateChild(d dino.PositionedDinosaur) {
	leftMid := (n.right-n.left)/2 + n.left
	topMid := (n.top-n.bottom)/2 + n.bottom
	left := d.Longitude < leftMid
	top := d.Latitude > topMid
	if left && top {
		n.topLeft.insert(d)
	} else if top {
		n.topRight.insert(d)
	} else if left {
		n.bottomLeft.insert(d)
	} else {
		n.bottomRight.insert(d)
	}
}

func (n *Node) splitIfNecessary() {
	if len(n.elements) > elementSize {
		n.split()
	}
}

func (n *Node) split() {
	leftMid := (n.right-n.left)/2 + n.left
	topMid := (n.top-n.bottom)/2 + n.bottom
	n.topLeft = &Node{left: n.left, right: leftMid, top: n.top, bottom: topMid}
	n.topRight = &Node{left: leftMid, right: n.right, top: n.top, bottom: topMid}
	n.bottomLeft = &Node{left: n.left, right: leftMid, top: topMid, bottom: n.bottom}
	n.bottomRight = &Node{left: leftMid, right: n.right, top: topMid, bottom: n.bottom}
	for _, e := range n.elements {
		n.insertAppropriateChild(e)
	}
	n.elements = nil
}

func (q *Quadtree) Get(top float64, left float64, bottom float64, right float64) []dino.PositionedDinosaur {
	dinos, expired := q.root.get(top, left, bottom, right)
	for _, dino := range expired {
		delete(q.dinos, dino.PositionedID)
	}
	return dinos
}

// get returns the dinosaurs in the given rect, and the dinosaurs which have expired and should be removed from the Quadtree's map of IDs.
func (n *Node) get(top float64, left float64, bottom float64, right float64) ([]dino.PositionedDinosaur, []dino.PositionedDinosaur) {
	if top < n.bottom || bottom > n.top || left > n.right || right < n.left {
		return nil, nil
	}

	if !n.hasSplit() {
		expired := n.removeExpired()
		return elementsInRect(top, left, bottom, right, n.elements), expired
	}

	topLeftElements, topLeftExpired := n.topLeft.get(top, left, bottom, right)
	topRightElements, topRightExpired := n.topRight.get(top, left, bottom, right)
	bottomLeftElements, bottomLeftExpired := n.bottomLeft.get(top, left, bottom, right)
	bottomRightElements, bottomRightExpired := n.bottomRight.get(top, left, bottom, right)

	allExpired := make([]dino.PositionedDinosaur, 0, len(topLeftExpired)+len(topRightExpired)+len(bottomLeftExpired)+len(bottomRightExpired))
	allExpired = append(allExpired, topLeftExpired...)
	allExpired = append(allExpired, topRightExpired...)
	allExpired = append(allExpired, bottomLeftExpired...)
	allExpired = append(allExpired, bottomRightExpired...)

	all := make([]dino.PositionedDinosaur, 0, len(topLeftElements)+len(topRightElements)+len(bottomLeftElements)+len(bottomRightElements))
	all = append(all, topLeftElements...)
	all = append(all, topRightElements...)
	all = append(all, bottomLeftElements...)
	all = append(all, bottomRightElements...)
	return all, allExpired
}

func inRect(rtop float64, rleft float64, rbottom float64, rright float64, latitude float64, longitude float64) bool {
	return latitude < rtop && latitude > rbottom && longitude > rleft && longitude < rright
}

func elementsInRect(rtop float64, rleft float64, rbottom float64, rright float64, elements []dino.PositionedDinosaur) []dino.PositionedDinosaur {
	in := []dino.PositionedDinosaur{}
	for _, e := range elements {
		if inRect(rtop, rleft, rbottom, rright, e.Latitude, e.Longitude) {
			in = append(in, e)
		}
	}
	return in
}

// removeExpired removes expired elements.
// Note this only removes expired elements from this node, not children.
// Returns expired dinosaurs, so they may be removed from the Quadtree map of IDs.
func (n *Node) removeExpired() []dino.PositionedDinosaur {
	now := time.Now()
	newElements := []dino.PositionedDinosaur{}
	expired := []dino.PositionedDinosaur{}
	for _, e := range n.elements {
		if !e.Expiration.Before(now) {
			newElements = append(newElements, e)
		} else {
			expired = append(expired, e)
		}
	}
	n.elements = newElements
	return expired
}
