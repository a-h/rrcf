package rrcf

import "fmt"

// Branch of RCTree containing two children and at most one parent.
type Branch struct {
	// q: Dimension of cut
	q int
	// p: Value of cut
	p float64
	// l: Pointer to left child
	l TreeBranchLeaf
	// r: Pointer to right child
	r TreeBranchLeaf
	// u: Pointer to parent
	u TreeNode
	// n: Number of leaves under branch
	n int
	// b: Bounding box of points under branch (2 x d)
	b BoundingBox
}

func (b *Branch) GetU() TreeNode {
	return b.u
}

func (b *Branch) SetU(u TreeNode) {
	b.u = u
}

func (b *Branch) GetN() int {
	return b.n
}

func (b *Branch) SetN(n int) {
	b.n = n
}

func (b *Branch) GetB() BoundingBox {
	return b.b
}

func (b *Branch) SetB(bb BoundingBox) {
	b.b = bb
}

func (b *Branch) String() string {
	return fmt.Sprintf("Branch(q=%v, p={:%2f})", b.q, b.p)
}
