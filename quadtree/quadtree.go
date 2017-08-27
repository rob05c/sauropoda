package quadtree

import (
	"github.com/rob05c/sauropoda/dinosaur"
	"math/rand"
	"time"
)

const elementSize = 100

type Quadtree struct {
	root *Node
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
	elements    []PositionedDinosaur
}

func (n *Node) hasSplit() bool {
	return n.topLeft != nil // it has split, if any child nodes are non-nil
}

// TODO add json lowercase tags
type PositionedDinosaur struct {
	Dinosaur   dinosaur.Dinosaur
	Latitude   float64
	Longitude  float64
	Expiration time.Time
	ID         int64
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
	return Quadtree{root: &Node{left: l, right: r, top: t, bottom: b, elements: make([]PositionedDinosaur, 0, elementSize)}}
}

// Insert inserts the given dinosaur into the quadtree
func (q Quadtree) Insert(d PositionedDinosaur) {
	q.root.insert(d)
}

func (n *Node) insert(d PositionedDinosaur) {
	if !n.hasSplit() {
		n.elements = append(n.elements, d)
		n.splitIfNecessary()
		return
	}

	n.insertAppropriateChild(d)
}

func (n *Node) insertAppropriateChild(d PositionedDinosaur) {
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

func (q Quadtree) Get(top float64, left float64, bottom float64, right float64) []PositionedDinosaur {
	return q.root.get(top, left, bottom, right)
}

func (n *Node) get(top float64, left float64, bottom float64, right float64) []PositionedDinosaur {
	if top < n.bottom || bottom > n.top || left > n.right || right < n.left {
		return nil
	}

	if !n.hasSplit() {
		n.removeExpired()
		return elementsInRect(top, left, bottom, right, n.elements)
	}

	topLeftElements := n.topLeft.get(top, left, bottom, right)
	topRightElements := n.topRight.get(top, left, bottom, right)
	bottomLeftElements := n.bottomLeft.get(top, left, bottom, right)
	bottomRightElements := n.bottomRight.get(top, left, bottom, right)

	all := make([]PositionedDinosaur, 0, len(topLeftElements)+len(topRightElements)+len(bottomLeftElements)+len(bottomRightElements))
	all = append(all, topLeftElements...)
	all = append(all, topRightElements...)
	all = append(all, bottomLeftElements...)
	all = append(all, bottomRightElements...)
	return all
}

func inRect(rtop float64, rleft float64, rbottom float64, rright float64, latitude float64, longitude float64) bool {
	return latitude < rtop && latitude > rbottom && longitude > rleft && longitude < rright
}

func elementsInRect(rtop float64, rleft float64, rbottom float64, rright float64, elements []PositionedDinosaur) []PositionedDinosaur {
	in := []PositionedDinosaur{}
	for _, e := range elements {
		if inRect(rtop, rleft, rbottom, rright, e.Latitude, e.Longitude) {
			in = append(in, e)
		}
	}
	return in
}

// removeExpired removes expired elements.
// Note this only removes expired elements from this node, not children.
func (n *Node) removeExpired() {
	now := time.Now()
	newElements := []PositionedDinosaur{}
	for _, e := range n.elements {
		if !e.Expiration.Before(now) {
			newElements = append(newElements, e)
		}
	}
	n.elements = newElements
}
