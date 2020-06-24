package binarysearchtree

import (
	"testing"
)

type Number uint32

func (v Number) Value() uint32 {
	return uint32(v)
}

func TestBinarySeachTree(t *testing.T) {
	bst := BinarySearchTree{}
	numbers := []Number{7, 3, 13, 5, 19, 12, 22, 2, 1, 5, 4}
	for _, n := range numbers {
		bst.Add(n)
	}
	if smallest := bst.Smallest(); smallest.Value() != 1 {
		t.Errorf("Expecting %v got %v", 1, smallest.Value())
	}
	if largest := bst.Largest(); largest.Value() != 22 {
		t.Errorf("Expecting %v got %v", 22, largest.Value())
	}

	iter := bst.Traverse()
	traverseCount := 0
	for iter.Next() {
		if iter.Current == nil {
			t.Errorf("iter.Current value is nil")
		}
		traverseCount++
	}
	if traverseCount != len(numbers) {
		t.Errorf(
			"Traverse don't output the same len of numbers numbers=%v traverse=%v",
			len(numbers),
			traverseCount,
		)
	}
}
