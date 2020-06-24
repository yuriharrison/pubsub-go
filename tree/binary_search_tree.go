package binarysearchtree

// IntEvaluable interface
type IntEvaluable interface {
	Value() uint32
}

// BinarySearchTree implementation
type BinarySearchTree struct {
	root *node
	sz   int
}

// Size return tree size
func (bst *BinarySearchTree) Size() int {
	return bst.sz
}

// IsEmpty return if the binary tree is empty
func (bst *BinarySearchTree) IsEmpty() bool {
	return bst.sz < 1
}

// Add an new node to the tree
func (bst *BinarySearchTree) Add(elem IntEvaluable) {
	newNode := &node{data: elem}
	if bst.IsEmpty() {
		bst.root = newNode
	} else {
		var rAdd func(n *node, new *node)
		rAdd = func(n *node, new *node) {
			switch {
			case n.ge(new.data.Value()) && n.emptyLeft():
				n.left = new
			case n.lt(new.data.Value()) && n.emptyRight():
				n.right = new
			case n.ge(new.data.Value()):
				rAdd(n.left, new)
			default:
				rAdd(n.right, new)
			}
			return
		}
		rAdd(bst.root, newNode)
	}
	bst.sz++
}

// Remove a node from the tree by value
func (bst *BinarySearchTree) Remove(value uint32) {
	var rRemove func(n *node, v uint32) *node
	rRemove = func(n *node, v uint32) *node {
		switch {
		case n.gt(v):
			n.left = rRemove(n.left, v)
		case n.lt(v):
			n.right = rRemove(n.right, v)
		default:
			switch {
			case n.emptyLeft() && n.emptyRight():
				return nil
			case n.emptyLeft():
				return n.right
			case n.emptyRight():
				return n.left
			default:
				tmp := bst.smallest(n.right)
				n.data = tmp.data
				n.right = rRemove(n.right, tmp.data.Value())
			}
		}
		return n
	}
	bst.root = rRemove(bst.root, value)
	bst.sz--
}

// Exists return if value exists in the tree
func (bst *BinarySearchTree) Exists(value uint32) bool {
	return bst.Find(value) != nil
}

// Find return the data for the given value in the tree
func (bst *BinarySearchTree) Find(value uint32) IntEvaluable {
	var rFind func(n *node) IntEvaluable
	rFind = func(n *node) IntEvaluable {
		switch {
		case n == nil:
			return nil
		case n.eq(value):
			return n.data
		case n.gt(value):
			return rFind(n.left)
		default:
			return rFind(n.right)
		}
	}
	return rFind(bst.root)
}

// Smallest returns the data for the smallest value in the tree
func (bst *BinarySearchTree) Smallest() IntEvaluable {
	return bst.smallest(bst.root).data
}

func (bst *BinarySearchTree) smallest(n *node) *node {
	for !n.emptyLeft() {
		n = n.left
	}
	return n
}

// Largest returns the data for the largest value in the tree
func (bst *BinarySearchTree) Largest() IntEvaluable {
	return bst.largest(bst.root).data
}

func (bst *BinarySearchTree) largest(n *node) *node {
	for !n.emptyRight() {
		n = n.right
	}
	return n
}

// Iterable go through values sent by the channel
type Iterable struct {
	Current IntEvaluable
	channel chan IntEvaluable
}

// Next return false when the channel is closed
func (iter *Iterable) Next() bool {
	if v, ok := <-iter.channel; ok {
		iter.Current = v
		return ok
	}
	return false
}

// Traverse tree values
func (bst *BinarySearchTree) Traverse() *Iterable {
	ch := make(chan IntEvaluable)
	var traverse func(n *node, main bool)
	traverse = func(node *node, main bool) {
		if node == nil {
			close(ch)
		}
		if node.left != nil {
			traverse(node.left, false)
		}
		if node.right != nil {
			traverse(node.right, false)
		}
		ch <- node.data
		if main {
			close(ch)
		}
	}
	go traverse(bst.root, true)
	return &Iterable{channel: ch}
}

type node struct {
	left  *node
	right *node
	data  IntEvaluable
}

func (n *node) isEmpty() bool {
	return n.data == nil
}

func (n *node) emptyLeft() bool {
	return n.left == nil
}

func (n *node) compare(value uint32) bool {
	return n.data.Value() >= value
}

func (n *node) emptyRight() bool {
	return n.right == nil
}

func (n *node) eq(value uint32) bool {
	return n.data.Value() == value
}

func (n *node) lt(value uint32) bool {
	return n.data.Value() < value
}

func (n *node) le(value uint32) bool {
	return n.data.Value() <= value
}

func (n *node) gt(value uint32) bool {
	return n.data.Value() > value
}

func (n *node) ge(value uint32) bool {
	return n.data.Value() >= value
}
