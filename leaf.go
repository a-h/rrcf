package rrcf

import "fmt"

// Leaf of RCTree containing no children and at most one parent.
type Leaf struct {
	//     i: Index of leaf (user-specified)
	i int
	//     d: Depth of leaf
	d int
	//     u: Pointer to parent
	u TreeNode
	//     x: Original point (1 x d)
	x []float64
	//     n: Number of points in leaf (1 if no duplicates)
	n int
	//     b: Bounding box of point (1 x d)
	b BoundingBox
}

func (self *Leaf) GetU() TreeNode {
	return self.u
}

func (self *Leaf) SetU(u TreeNode) {
	self.u = u
}

func (self *Leaf) GetN() int {
	return self.n
}

func (self *Leaf) SetN(n int) {
	self.n = n
}

func (self *Leaf) GetB() BoundingBox {
	return self.b
}

func (self *Leaf) SetB(bb BoundingBox) {
	self.b = bb
}

func (self *Leaf) String() string {
	return fmt.Sprintf("Leaf(%d)", self.i)
}
