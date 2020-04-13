package rrcf

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"time"

	"github.com/a-h/round"

	"github.com/a-h/linear/tolerance"
)

const precision = 9 // decimal places
var toleranceLevel = math.Pow(10, -(precision + 1.5))

type RCTree struct {
	rng rand.Source
	// leaves: dict containing pointers to all leaves in tree.
	leaves map[int]TreeNode
	// root: Branch or Leaf instance. Pointer to root of tree.
	root TreeNode
	// ndim is the dimension of points in the tree
	ndim         *int
	index_labels []int
	u            TreeNode
}

func (t *RCTree) GetU() TreeNode {
	return t.u
}

func (t *RCTree) SetU(u TreeNode) {
	t.u = u
}

func NewFloatVector(from []float64) FloatVector {
	fv := make(FloatVector, len(from))
	for i, v := range from {
		fv[i] = v
	}
	return fv
}

type FloatVector []float64

type FloatPredicate func(v float64) bool

func IsLessThanOrEqual(v float64) FloatPredicate {
	return func(input float64) bool {
		return input <= v
	}
}

func (fv FloatVector) Where(f FloatPredicate) BoolVector {
	op := make(BoolVector, len(fv))
	for i, v := range fv {
		op[i] = f(v)
	}
	return op
}

func (fv FloatVector) Subtract(other FloatVector) (op FloatVector) {
	op = make(FloatVector, len(fv))
	for i, v := range fv {
		op[i] = v - other[i]
	}
	return op
}

func (fv FloatVector) Divide(other FloatVector) (op FloatVector) {
	op = make(FloatVector, len(fv))
	for i, v := range fv {
		op[i] = v / other[i]
	}
	return op
}

func (fv FloatVector) DivideBy(scalar float64) (op FloatVector) {
	op = make(FloatVector, len(fv))
	for i, v := range fv {
		op[i] = v / scalar
	}
	return op
}

func (fv FloatVector) Sum() (op float64) {
	for _, v := range fv {
		op += v
	}
	return
}

type BoolVector []bool

type BoolPredicate func(b bool) bool

func IsTrue(b bool) bool {
	return b
}

// IndicesWhere is equivalent to np.flatnonzero.
func (bv BoolVector) IndicesWhere(f BoolPredicate) (indices []int) {
	for i, v := range bv {
		if f(v) {
			indices = append(indices, i)
		}
	}
	return
}

func (bv BoolVector) Where(f BoolPredicate) (op BoolVector) {
	for _, v := range bv {
		if f(v) {
			op = append(op, v)
		}
	}
	return
}

func (bv BoolVector) And(other BoolVector) (op BoolVector) {
	op = make(BoolVector, len(bv))
	for i, v := range bv {
		op[i] = v && other[i]
	}
	return op
}

func (bv BoolVector) Not() (op BoolVector) {
	op = make(BoolVector, len(bv))
	for i, v := range bv {
		op[i] = !v
	}
	return op
}

func (bv BoolVector) Sum() (sum int) {
	for _, v := range bv {
		if v {
			sum++
		}
	}
	return
}

func NewMatrix(cols int) *Matrix {
	return &Matrix{
		Data: make([][]float64, 0),
		cols: cols,
	}
}

type Matrix struct {
	Data [][]float64
	cols int
}

func (a *Matrix) Rows() int {
	return len(a.Data)
}

func (a *Matrix) Cols() int {
	return a.cols
}

func (a *Matrix) Column(index int) (op FloatVector) {
	op = make(FloatVector, len(a.Data))
	for i, row := range a.Data {
		op[i] = row[index]
	}
	return op
}

func (a *Matrix) Round(decimals int) (output *Matrix) {
	output = NewMatrix(a.Cols())
	output.Data = make([][]float64, a.Rows())
	for i, row := range a.Data {
		output.Data[i] = make([]float64, a.Cols())
		for j, v := range row {
			output.Data[i][j] = round.ToEven(v, precision)
		}
	}
	return output
}

func (a *Matrix) IndexOf(r []float64) int {
	for i := 0; i < a.Rows(); i++ {
		if a.RowEqual(i, r) {
			return i
		}
	}
	return -1
}

func (a *Matrix) RowEqual(i int, other []float64) bool {
	r := a.Data[i]
	if len(r) != len(other) {
		return false
	}
	for i := 0; i < len(r); i++ {
		if !tolerance.IsWithin(r[i], other[i], toleranceLevel) {
			return false
		}
	}
	return true
}

func (a *Matrix) Unique() (unique *Matrix, inverse []int, counts []int) {
	unique = NewMatrix(a.Cols())
	for _, row := range a.Data {
		//TODO: Hash the row instead, for performance.
		index := unique.IndexOf(row)
		if index == -1 {
			unique.Data = append(unique.Data, row)
			index = unique.Rows() - 1
		} else {
			counts[index]++
		}
		inverse = append(inverse, index)
	}
	return
}

func (a *Matrix) Skip(from int) (output *Matrix) {
	output = NewMatrix(a.Cols())
	output.Data = a.Data[from:]
	return output
}

func (a *Matrix) Include(includeRows []bool) (output *Matrix) {
	output = NewMatrix(a.Cols())
	for i, r := range a.Data {
		include := includeRows[i]
		if include {
			output.Data = append(output.Data, r)
		}
	}
	return
}

// Maximum value of each column across all rows.
func (a *Matrix) Max() (max FloatVector) {
	max = NewFloatVector(a.Data[0])
	for _, r := range a.Data[1:] {
		for col, v := range r {
			if v > max[col] {
				max[col] = v
			}
		}
	}
	return max
}

// Minimum value of each column across all rows.
func (a *Matrix) Min() (min FloatVector) {
	min = NewFloatVector(a.Data[0])
	for _, r := range a.Data[1:] {
		for col, v := range r {
			if v < min[col] {
				min[col] = v
			}
		}
	}
	return min
}

// Sum of each column value across all rows.
func (a *Matrix) Sum() (sum float64) {
	for _, r := range a.Data {
		for _, v := range r {
			sum += v
		}
	}
	return
}

func (a *Matrix) Subtract(b *Matrix) (output *Matrix) {
	output = NewMatrix(a.Cols())
	output.Data = make([][]float64, a.Rows())
	for i, row := range a.Data {
		output.Data[i] = make([]float64, a.Cols())
		for j, v := range row {
			output.Data[i][j] = v - b.Data[i][j]
		}
	}
	return output
}

func (a *Matrix) Divide(by float64) (output *Matrix) {
	output = NewMatrix(a.Cols())
	output.Data = make([][]float64, a.Rows())
	for i, row := range a.Data {
		output.Data[i] = make([]float64, a.Cols())
		for j, v := range row {
			output.Data[i][j] = v / by
		}
	}
	return output
}

func maxIntOrZero(in []int) (max int) {
	if len(in) == 0 {
		return
	}
	max = in[0]
	if len(in) == 1 {
		return
	}
	for _, v := range in[1:] {
		if v > max {
			max = v
		}
	}
	return
}

func npOnes(count int) []int {
	op := make([]int, count)
	for i := 0; i < count; i++ {
		op[i] = 1
	}
	return op
}

func npBools(count int) []bool {
	op := make([]bool, count)
	for i := 0; i < count; i++ {
		op[i] = true
	}
	return op
}

// TreeNode is either a Branch or Leaf
type TreeNode interface {
	GetU() TreeNode
	SetU(n TreeNode)
}

type TreeBranchLeaf interface {
	TreeNode
	GetN() int
	SetN(int)
	// GetB gets the bounding box.
	GetB() BoundingBox
	// SetB sets the bounding box.
	SetB(BoundingBox)
}

func New(X *Matrix) *RCTree {
	self := &RCTree{
		rng:    rand.NewSource(time.Now().Unix()),
		leaves: map[int]TreeNode{},
		root:   nil,
		ndim:   nil,
	}
	if X != nil {
		// Round data to avoid sorting errors
		X = X.Round(precision)
		// Initialize index labels, if they exist
		self.index_labels = make([]int, X.Rows())
		for i := 0; i < X.Rows(); i++ {
			self.index_labels[i] = i
		}
		// Check for duplicates.
		U, I, N := X.Unique()
		// If duplicates exist, take unique elements
		var n, d int
		if maxIntOrZero(N) > 1 {
			n, d = U.Rows(), U.Cols()
			X = U
		} else {
			n, d = X.Rows(), X.Cols()
			N = npOnes(n)
			I = nil
		}
		// Store dimension of dataset
		self.ndim = &d
		// Set node above to None in case of bottom-up search
		self.u = nil
		// Create RRC Tree
		S := npBools(n)
		self._mktree(X, S, N, I, self, "root", 0)
		// Remove parent of root
		self.root.SetU(nil)
		// Count all leaves under each branch
		self._count_all_top_down(self.root)
		// Set bboxes of all branches
		self._get_bbox_top_down(self.root)
	}
	return self
}

func print_push(depth, treestr string, char rune) (string, string) {
	depth += fmt.Sprintf(" %c  ", char)
	return depth, treestr
}

func print_pop(depth, treestr string) (string, string) {
	depth = depth[:len(depth)-4]
	return depth, treestr
}

func print_tree(depth, treestr string, n TreeNode) (string, string) {
	if l, ok := n.(*Leaf); ok {
		treestr += fmt.Sprintf("(%d)\n", l.i)
	}
	if b, ok := n.(*Branch); ok {
		treestr += fmt.Sprintf("%c%v\n", rune(9472), '+')
		treestr += fmt.Sprintf("%s %c%c%c", depth, rune(9500), rune(9472), rune(9472))
		depth, treestr = print_push(depth, treestr, rune(9474))
		depth, treestr = print_tree(depth, treestr, b.l)
		depth, treestr = print_pop(depth, treestr)
		treestr += fmt.Sprintf("%s %c%c%c", depth, rune(9492), rune(9472), rune(9472))
		depth, treestr = print_push(depth, treestr, ' ')
		depth, treestr = print_tree(depth, treestr, b.r)
		depth, treestr = print_pop(depth, treestr)
	}
	return depth, treestr
}

func (t *RCTree) String() string {
	// return fmt.Sprintf("%v\n", reflect.TypeOf(t.root))
	_, treestr := print_tree("", "", t.root)
	return treestr
}

func choose(rng rand.Source, probabilities []float64) (index int) {
	rnd := rand.New(rng)
	var max float64
	for i, p := range probabilities {
		if v := rnd.Float64() * p; v > max {
			index = i
			max = v
		}
	}
	return
}

// Compute bbox of node based on bboxes of node's children.
func (self *RCTree) _lr_branch_bbox(node *Branch) BoundingBox {
	if node.r == nil {
		panic("nil node")
	}
	lbb, rbb := node.l.GetB(), node.r.GetB()
	bb := BoundingBox{
		Ax: lbb.Ax,
		Ay: lbb.Ay,
		Dx: lbb.Dy,
		Dy: lbb.Dy,
	}
	if rbb.Ax < bb.Ax {
		bb.Ax = rbb.Ax
	}
	if rbb.Ay > bb.Ay {
		bb.Ay = rbb.Ay
	}
	if rbb.Dx < bb.Dx {
		bb.Dx = rbb.Dx
	}
	if rbb.Dy > bb.Dy {
		bb.Dy = rbb.Dy
	}
	return bb
}

// Recursively compute bboxes of all branches from root to leaves.
func (self *RCTree) _get_bbox_top_down(node TreeNode) {
	if n, ok := node.(*Branch); ok {
		fmt.Println("it's a branch", reflect.TypeOf(n))
		if n.l != nil {
			self._get_bbox_top_down(n.l)
		}
		if n.r != nil {
			self._get_bbox_top_down(n.r)
		}
		n.b = self._lr_branch_bbox(n)
	} else {
		fmt.Println("it ain't no branch", reflect.TypeOf(n), n)
	}
}

// Recursively compute number of leaves below each branch from root to leaves.
func (self *RCTree) _count_all_top_down(node TreeNode) {
	if n, ok := node.(*Branch); ok {
		if n.l != nil {
			self._count_all_top_down(n.l)
		}
		if n.r != nil {
			self._count_all_top_down(n.r)
		}
		var count int
		if l, ok := n.l.(TreeBranchLeaf); ok {
			fmt.Println("TreeBranchLeafLeft", reflect.TypeOf(n))
			count += l.GetN()
		}
		if r, ok := n.r.(TreeBranchLeaf); ok {
			fmt.Println("TreeBranchLeafRight", reflect.TypeOf(n))
			count += r.GetN()
		}
		n.n = count
	}
}

// uniform selects a random float from the given range with a uniform distribution of
// probability.
func uniform(rng rand.Source, from, to float64) float64 {
	return (rand.New(rng).Float64() * (to - from)) + from
}

func (self *RCTree) _cut(X *Matrix, S []bool, parent TreeNode, side string) (S1 BoolVector, S2 BoolVector, child *Branch) {
	// Find max and min over all d dimensions
	xmax := X.Include(S).Max()
	xmin := X.Include(S).Min()
	// Compute l
	l := xmax.Subtract(xmin)
	l = l.DivideBy(l.Sum())
	// Determine dimension to cut
	q := choose(self.rng, l)
	// Determine value for split
	p := uniform(self.rng, xmin[q], xmax[q])
	// Determine subset of points to left
	S1 = X.Column(q).Where(IsLessThanOrEqual(p)).And(S)
	// Determine subset of points to right
	S2 = S1.Not().And(S)
	// Create new child node
	child = &Branch{
		q: q,
		p: p,
		u: parent,
	}
	// Link child node to parent
	if parent != nil {
		if t, ok := parent.(*RCTree); ok {
			switch side {
			case "root":
				t.root = child
				break
			}
		}
		if b, ok := parent.(*Branch); ok {
			switch side {
			case "l":
				b.l = child
				break
			case "r":
				b.r = child
				break
			}
		}
	}
	return
}

func indicesOfValue(array []int, value int) (indices []int) {
	for i, v := range array {
		if v == value {
			indices = append(indices, i)
		}
	}
	return
}

func (self *RCTree) _mktree(X *Matrix, S []bool, N []int, I []int, parent TreeNode, side string, depth int) {
	// Increment depth as we traverse down
	depth++
	// Create a cut according to definition 1
	S1, S2, branch := self._cut(X, S, parent, side)
	// If S1 does not contain an isolated point...
	if S1.Sum() > 1 {
		// Recursively construct tree on S1
		self._mktree(X, S1, N, I, branch, "l", depth)
		// Otherwise...
	} else {
		// Create a leaf node from isolated point
		// There is only one true value, or we'd have ended up in the branch above
		i := S1.IndicesWhere(IsTrue)[0]
		leaf := &Leaf{
			i: i,
			d: depth,
			u: branch,
			x: X.Data[i][:],
			n: N[i],
		}
		// Link leaf node to parent
		branch.l = leaf
		// If duplicates exist...
		if I != nil {
			// Add a key in the leaves dict pointing to leaf for all duplicate indices
			J := indicesOfValue(I, i)[0]
			// Get index label
			label := self.index_labels[J]
			self.leaves[label] = leaf
		} else {
			i = self.index_labels[i]
			self.leaves[i] = leaf
		}
		// If S2 does not contain an isolated point...
		if S2.Sum() > 1 {
			// Recursively construct tree on S2
			self._mktree(X, S2, N, I, branch, "r", depth)
		} else {
			// Create a leaf node from isolated point
			// There is only one 'true' value, or we'd have ended up in the branch above
			i := S1.IndicesWhere(IsTrue)[0]
			leaf = &Leaf{
				i: i,
				d: depth,
				u: branch,
				x: X.Data[i][:],
				n: N[i],
			}
			// Link leaf node to parent
			branch.r = leaf
			// If duplicates exist...
			if I != nil {
				// Add a key in the leaves dict pointing to leaf for all duplicate indices
				J := indicesOfValue(I, i)[0]
				// Get index label
				label := self.index_labels[J]
				self.leaves[label] = leaf
			} else {
				i = self.index_labels[i]
				self.leaves[i] = leaf
			}
		}
	}
	// Decrement depth as we traverse back up
	depth--
}
