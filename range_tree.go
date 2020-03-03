package range_tree

import (
	"math"
)

type node struct {
	left, right  *node
	min          int64
	modification int64
}

type Tree struct {
	head               *node
	size               int64
	globalModification int64
}

func NewTree(size int64, globalModification int64) *Tree {
	if size < 0 {
		panic("size should be bigger than zero")
	}
	next := nextClosestPowerOf2(size) - 1
	return &Tree{head: &node{}, size: next, globalModification: globalModification}
}

func (tree *Tree) Update(value int64, leftEdge int64, rightEdge int64) {
	if leftEdge < 0 || rightEdge < 0 {
		panic("Wrong edges for tree update")
	}
	if tree.size < rightEdge {
		tree.ensureCapacity(rightEdge)
	}
	tree.head.update(value, 0, tree.size, leftEdge, rightEdge)
}

func (tree *Tree) ensureCapacity(rightEdge int64) {
	tree.size++
	leftNode := tree.head
	// do while
	for {
		leftNode = &node{left: leftNode}
		tree.size *= 2
		if !(tree.size-1 < rightEdge) {
			break
		}
	}
	tree.head = leftNode
	tree.size--
}

func (currentNode *node) updateMin() {
	var newMinRight, newMinLeft int64
	if currentNode.right != nil {
		newMinRight = currentNode.right.modification + currentNode.right.min
	}
	if currentNode.left != nil {
		newMinLeft = currentNode.left.modification + currentNode.left.min
	}
	currentNode.min = min(newMinLeft, newMinRight)
}

func (currentNode *node) update(value int64, leftNodeEdge int64, rightNodeEdge int64, leftEdge int64, rightEdge int64) int64 {
	if (leftNodeEdge == leftEdge) && (rightNodeEdge == rightEdge) {
		currentNode.modification += value
	} else {
		midNodeEdge := leftNodeEdge + (rightNodeEdge-leftNodeEdge)/2
		switch {
		case midNodeEdge >= rightEdge:
			if currentNode.left == nil {
				currentNode.left = &node{}
			}
			currentNode.left.update(value, leftNodeEdge, midNodeEdge, leftEdge, rightEdge)
			currentNode.updateMin()
		case midNodeEdge < leftEdge:
			if currentNode.right == nil {
				currentNode.right = &node{}
			}
			currentNode.right.update(value, midNodeEdge+1, rightNodeEdge, leftEdge, rightEdge) // +1??
			currentNode.updateMin()
		case rightEdge > midNodeEdge:
			if currentNode.left == nil {
				currentNode.left = &node{}
			}
			if currentNode.right == nil {
				currentNode.right = &node{}
			}
			currentNode.left.update(value, leftNodeEdge, midNodeEdge, leftEdge, midNodeEdge)
			currentNode.right.update(value, midNodeEdge+1, rightNodeEdge, midNodeEdge+1, rightEdge)
			currentNode.updateMin()
		}
	}
	return currentNode.min + currentNode.modification
}

func (tree *Tree) Get(leftEdge int64, rightEdge int64) int64 {
	if leftEdge < 0 || rightEdge < 0 || leftEdge > rightEdge  {
		panic("Wrong values for tree get")
	}

	effectiveRightEdge := min(tree.size, rightEdge)
	newMin := tree.head.get(0, tree.size, leftEdge, effectiveRightEdge)
	if leftEdge > tree.size {
		return tree.globalModification
	}

	if effectiveRightEdge < rightEdge {
		newMin = min(newMin, 0)
	}
	return tree.globalModification + newMin
}

func (currentNode *node) get(leftNodeEdge int64, rightNodeEdge int64, leftEdge int64, rightEdge int64) int64 {
	if !((leftNodeEdge == leftEdge) && (rightNodeEdge == rightEdge)) {
		midNodeEdge := leftNodeEdge + (rightNodeEdge-leftNodeEdge)/2
		var newMinLeft int64
		var newMinRight int64
		switch {
		case midNodeEdge >= rightEdge:
			if currentNode.left != nil {
				newMinLeft = currentNode.left.get(leftNodeEdge, midNodeEdge, leftEdge, rightEdge)
			}
			return newMinLeft + currentNode.modification
		case midNodeEdge < leftEdge:
			if currentNode.right != nil {
				newMinRight = currentNode.right.get(midNodeEdge+1, rightNodeEdge, leftEdge, rightEdge)
			}
			return newMinRight + currentNode.modification
		case rightEdge > midNodeEdge:
			if currentNode.left != nil {
				newMinLeft = currentNode.left.get(leftNodeEdge, midNodeEdge, leftEdge, midNodeEdge)
			}
			if currentNode.right != nil {
				newMinRight = currentNode.right.get(midNodeEdge+1, rightNodeEdge, midNodeEdge+1, rightEdge)
			}
			return min(newMinLeft, newMinRight) + currentNode.modification
		}
	}
	return currentNode.min + currentNode.modification
}

func min(x int64, y int64) int64 {
	if x <= y {
		return x
	}
	return y
}

func nextClosestPowerOf2(x int64) int64 {
	return int64(math.Pow(2., math.Ceil(math.Log2(float64(x)))))
}
