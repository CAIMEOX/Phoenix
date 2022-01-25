package function

import (
	"fmt"
	"math"
)

// Octree An octree is a data structure that allows fast retrieval of elements based
// values in three dimensions.
type Octree struct {
	root *Node
}

// CreateOctree Makes a new octree with the given min and max.
func CreateOctree(min, max Vector3f) *Octree {
	mn := min.Min(&max)
	mx := min.Max(&max)
	o := Octree{}
	o.root = &Node{box: Box{min: mn, max: mx}}
	return &o
}

// Clear Removes all the data from the Octree while
// retaining its bounding box. Returns true if octree is ready for use
// (because it has previously been initialized).
func (o *Octree) Clear() bool {
	if o.root != nil {
		// if octree has been initializes, use the same box,
		// but create a new root, freeing the other memory
		// (except where outside references have been retained).
		o.root = &Node{box: o.root.box}
		return true
	}

	return false
}

// Add Inserts the element in the tree at the specified point.
// If you may need to remove the element later, retain the
// returned node for fast removal.
func (o *Octree) Add(element interface{}, point Vector3f) *Node {
	return o.root.tryAdd([]interface{}{element}, &point)
}

// ElementsAt Retrieves a slice of elements that exist at
// the specified point in the tree.
func (o *Octree) ElementsAt(point Vector3f) []interface{} {
	return o.root.elementsAt(&point)
}

// ElementsIn Retrieves a slice of element that exist
// within the specified box.
func (o *Octree) ElementsIn(box Box) []interface{} {
	return o.root.elementsIn(&box)
}

// Remove Removes the specified element from the tree.
// Generally, RemoveUsing should used as it is faster under
// most circumstances.
func (o *Octree) Remove(element interface{}) bool {
	return o.root.remove(element)
}

// RemoveUsing Removes the specified element from the tree; node constrains the search
// for the element and should usually be the node returned when this element
// was placed in the tree using Add()
func (o *Octree) RemoveUsing(element interface{}, node *Node) bool {
	if node != nil {
		return node.remove(element)
	}
	return false
}

// ToString Get a human readable representation of the state of
// this octree and its contents.
func (o *Octree) ToString() string {
	str := "nil"
	if o.root != nil {
		str = o.root.recursiveToString("  ", "  ")
	}

	return fmt.Sprintf("Octree{\n  root: %v\n}", str)
}

// Node An element within the tree that can either act as a leaf, that can directly hold a point
// and its corresponding elements or act as a branch and hold references to child nodes.
type Node struct {
	box         Box
	point       *Vector3f
	elements    []interface{}
	hasChildren bool
	children    []*Node
}

func (n *Node) tryAdd(elements []interface{}, point *Vector3f) *Node {
	// attempt to add the elements in this node (or a descendant)
	// at the specified point.

	if !n.box.ContainsPoint(point) {
		return nil
	}

	if n.hasChildren {
		return n.addToChildren(elements, point)
	}

	if n.point != nil {
		// leaf already has assigned point
		if *n.point == *point {
			// points are equal
			n.elements = append(n.elements, elements...)
			return n
		}

		// subdivide because points are different
		return n.subdivide(elements, point)
	}

	// set own elements and point
	n.elements = elements
	n.point = point

	return n
}

func (n *Node) addToChildren(elements []interface{}, point *Vector3f) *Node {
	for _, child := range n.children {
		// try adding to child
		leaf := child.tryAdd(elements, point)

		if leaf != nil {
			// succeeded adding
			return leaf
		}
	}

	// Error: box.contains evaluated to true, but none of the children added the point
	return nil
}

func (n *Node) subdivide(addElements []interface{}, atPoint *Vector3f) *Node {
	// create child nodes for what is currently a leaf,
	// moving its current contents to one of those leafs.

	// setup this node's children
	n.hasChildren = true
	subBoxes := n.box.makeSubBoxes()

	for i := 0; i < 8; i++ {
		n.children = append(n.children, &Node{box: subBoxes[i]})
	}

	// add node's elements and point to a child
	n.addToChildren(n.elements, n.point)

	// clear elements and point from self
	n.elements = nil
	n.point = nil

	// add the new element to a child
	return n.addToChildren(addElements, atPoint)
}

func (n *Node) elementsAt(point *Vector3f) []interface{} {
	// get any alements in this node (or a descendant)
	// at the specified point

	if n.hasChildren {
		for _, child := range n.children {
			if child.box.ContainsPoint(point) {
				return child.elementsAt(point)
			}
		}
	} else {
		// when a leaf
		if n.point != nil && *point == *n.point {
			return n.elements
		}
	}

	return nil
}

func (n *Node) elementsIn(box *Box) []interface{} {
	// get any alements in this node (or a descendant)
	// within the specified box

	if n.hasChildren {
		elements := []interface{}{}

		for _, child := range n.children {
			if child.box.IsContainedIn(box) {
				// fully contained
				elements = append(elements, child.elementsIn(&child.box)...)
			} else if child.box.Contains(box) || child.box.Intersects(box) {
				// partially contained
				elements = append(elements, child.elementsIn(box)...)
			}
		}

		return elements
	}

	// when a leaf
	if n.point != nil && box.ContainsPoint(n.point) {
		return n.elements
	}

	return nil
}

func (n *Node) remove(element interface{}) bool {
	// remove the first instance of the specified element
	// in this node (or in a descendant)

	if n.hasChildren {
		for _, child := range n.children {
			if child.remove(element) {
				return true
			}
		}
		return false
	}

	for idx, val := range n.elements {
		if val == element {
			// remove element from the slice
			n.elements = append(n.elements[:idx], n.elements[idx+1:]...)
			return true
		}
	}
	return false
}

// ToString Get a human readable representation of the state of
// this node and its contents.
func (n *Node) ToString() string {
	return n.recursiveToString("", "  ")
}

func (n *Node) recursiveToString(curIndent, stepIndent string) string {
	singleIndent := curIndent + stepIndent

	// default values
	childStr := "nil"
	pointStr := "nil"
	elementStr := "nil"

	if n.hasChildren {
		doubleIndent := singleIndent + stepIndent

		// accumulate child strings
		childStr = ""
		for i, child := range n.children {
			childStr = childStr + fmt.Sprintf("%v%d: %v,\n", doubleIndent, i, child.recursiveToString(doubleIndent, stepIndent))
		}

		childStr = fmt.Sprintf("[\n%v%v]", childStr, singleIndent)
	}

	if n.point != nil {
		pointStr = n.point.ToString()
	}

	if n.elements != nil {
		// not stringifying elements since their type is unknown
		elementStr = fmt.Sprintf("[%d]", len(n.elements))
	}

	return fmt.Sprintf("Node{\n%vchildren: %v,\n%vbox: %v,\n%vpoint: %v\n%velements: %v,\n%v}", singleIndent, childStr, singleIndent, n.box.ToString(), singleIndent, pointStr, singleIndent, elementStr, curIndent)
}

// Box Defines an axis aligned rectangular solid.
type Box struct {
	min Vector3f
	max Vector3f
}

// Size Returns the dimensions of the Box.
func (b *Box) Size() Vector3f {
	return b.max.Minus(&b.min)
}

// ContainsPoint Returns whether the specified point is contained in this box.
func (b *Box) ContainsPoint(v *Vector3f) bool {
	return b.min[0] <= v[0] &&
			b.max[0] >= v[0] &&
			b.min[1] <= v[1] &&
			b.max[1] >= v[1] &&
			b.min[2] <= v[2] &&
			b.max[2] >= v[2]
}

// Contains Returns whether the specified box is contained in this box.
func (b *Box) Contains(o *Box) bool {
	return b.min[0] <= o.min[0] &&
		b.max[0] >= o.max[0] &&
		b.min[1] <= o.min[1] &&
		b.max[1] >= o.max[1] &&
		b.min[2] <= o.min[2] &&
		b.max[2] >= o.max[2]
}

// IsContainedIn Returns whether the specified box contains this box.
func (b *Box) IsContainedIn(o *Box) bool {
	return o.Contains(b)
}

// Intersects Returns whether any portion of this box intersects with
// the specified box.
func (b *Box) Intersects(o *Box) bool {
	return !(b.max[0] < o.min[0] ||
		o.max[0] < b.min[0] ||
		b.max[1] < o.min[1] ||
		o.max[1] < b.min[1] ||
		b.max[2] < o.min[2] ||
		o.max[2] < b.min[2])
}

// ToString Get a human readable representation of the state of
// this box.
func (b *Box) ToString() string {
	return fmt.Sprintf("Box{min: %v, max: %v}", b.min.ToString(), b.max.ToString())
}

func (b *Box) makeSubBoxes() [8]Box {
	// gets the child boxes (octants) of the box.
	center := b.min.Lerp(&b.max, 0.5)

	return [8]Box{
		Box{Vector3f{b.min[1], b.min[1], b.min[2]}, Vector3f{center[0], center[1], center[2]}},
		Box{Vector3f{center[1], b.min[1], b.min[2]}, Vector3f{b.max[0], center[1], center[2]}},
		Box{Vector3f{b.min[1], center[1], b.min[2]}, Vector3f{center[0], b.max[1], center[2]}},
		Box{Vector3f{center[1], center[1], b.min[2]}, Vector3f{b.max[0], b.max[1], center[2]}},
		Box{Vector3f{b.min[1], b.min[1], center[2]}, Vector3f{center[0], center[1], b.max[2]}},
		Box{Vector3f{center[1], b.min[1], center[2]}, Vector3f{b.max[0], center[1], b.max[2]}},
		Box{Vector3f{b.min[1], center[1], center[2]}, Vector3f{center[0], b.max[1], b.max[2]}},
		Box{Vector3f{center[1], center[1], center[2]}, Vector3f{b.max[0], b.max[1], b.max[2]}},
	}
}

// Vector3f ...
type Vector3f [3]float64

// Minus Subtracts another Vector3f from this Vector3f and returns the result.
func (v *Vector3f) Minus(other *Vector3f) Vector3f {
	return Vector3f{v[0] - other[0], v[1] - other[1], v[2] - other[2]}
}

// Plus Returns the addition of the Vector3f(s).
func (v *Vector3f) Plus(other *Vector3f) Vector3f {
	return Vector3f{v[0] + other[0], v[1] + other[1], v[2] + other[2]}
}

// Scale Returns the multiplication of the Vector3f by a number.
func (v *Vector3f) Scale(f float64) Vector3f {
	return Vector3f{v[0] * f, v[1] * f, v[2] * f}
}

// Min Returns the a vector where each component is the lesser of the
// corresponding component in this and the specified vector.
func (v *Vector3f) Min(other *Vector3f) Vector3f {
	return Vector3f{
		math.Min(v[0], other[0]),
		math.Min(v[1], other[1]),
		math.Min(v[2], other[2]),
	}
}

// Max Returns the a vector where each component is the greater of the
// corresponding component in this and the specified vector.
func (v *Vector3f) Max(other *Vector3f) Vector3f {
	return Vector3f{
		math.Max(v[0], other[0]),
		math.Max(v[1], other[1]),
		math.Max(v[2], other[2]),
	}
}

// Lerp Returns the linear interpolation between two Vector3f(s).
func (v *Vector3f) Lerp(other *Vector3f, f float64) Vector3f {
	return Vector3f{
		(other[0]-v[0])*f + v[0],
		(other[1]-v[1])*f + v[1],
		(other[2]-v[2])*f + v[2],
	}
}

// ToString Get a human readable representation of the state of
// this vector.
func (v *Vector3f) ToString() string {
	return fmt.Sprintf("Vector3f{%f, %f, %f}", v[0], v[1], v[2])
}
