package range_tree

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

type FakeTree struct {
	array *[]int64
}

func NewFakeTree(count int64) *FakeTree {
	arr := make([]int64, count)
	for i := int64(0); i < count; i++ {
		arr[i] = 0
	}
	return &FakeTree{&arr}
}

func (ft FakeTree) Update(value int64, left int64, right int64) {
	arr := *ft.array
	for i := left; i <= right; i++ {
		arr[i] += value
	}
}

func (ft *FakeTree) Get(left int64, right int64) int64 {
	arr := *ft.array
	minElem := arr[left]
	for i := left + 1; i <= right; i++ {
		minElem = min(minElem, arr[i])
	}
	return minElem
}

func TestTree_GetUpdate(t *testing.T) {
	const size = 20000
	tree := NewTree(size, 0)
	fakeTree := NewFakeTree(size)
	for i := 0; i < 20000; i++ {
		left := rand.Int63() % (size)
		right := rand.Int63() % (size)
		updateWithRandom(t, left, right, tree, fakeTree)
	}
}

func TestTree_GetUpdateWithGlobal(t *testing.T) {
	const globalModification = 4
	const size = 2000
	tree := NewTree(size, globalModification)
	fakeTree := NewFakeTree(size)
	fakeTree.Update(globalModification, 0, size-1)
	for i := 0; i < 100000; i++ {
		left := rand.Int63() % (size)
		right := rand.Int63() % (size)
		updateWithRandom(t, left, right, tree, fakeTree)
	}
}

func TestTree_GetUpdateWithSize(t *testing.T) {
	size := int64(8)
	finSize := size * 128
	tree := NewTree(size, 0)
	fakeTree := NewFakeTree(finSize)
	for ; size <= finSize; size *= 2 {
		for i := 0; i < 10; i++ {
			left := rand.Int63() % (size)
			right := rand.Int63() % (size)
			updateWithRandom(t, left, right, tree, fakeTree)
		}
		if size*4 == finSize {
			size *= 2
		}
	}
}

func TestTree_GetUpdateWithSizeAndGlobal(t *testing.T) {
	const globalModification = 0
	size := int64(8)
	finSize := size * 64
	tree := NewTree(size, globalModification)
	fakeTree := NewFakeTree(finSize)
	fakeTree.Update(globalModification, 0, finSize-1)
	for ; size <= finSize; size *= 2 {
		for i := 0; i < 100; i++ {
			left := rand.Int63() % (size)
			right := rand.Int63() % (size)
			updateWithRandom(t, left, right, tree, fakeTree)
		}
	}
}

func TestTree_GetUpdateWithEdgeSize(t *testing.T) {
	size := int64(8)
	finSize := size * 256
	tree := NewTree(size, 0)
	fakeTree := NewFakeTree(finSize)
	modificator := int64(2)
	for ; size <= finSize; size *= modificator {
		left := size
		right := size
		updateWithRandom(t, left, right, tree, fakeTree)
		modificator *= 2
	}
}

func updateWithRandom(t *testing.T, randomLeft int64, randomRight int64, tree *Tree, fakeTree *FakeTree) {
	if randomLeft > randomRight {
		tmp := randomLeft
		randomLeft = randomRight
		randomRight = tmp
	}
	tree.Update(1-((randomLeft+randomRight)%2)*2, randomLeft, randomRight)
	fakeTree.Update(1-((randomLeft+randomRight)%2)*2, randomLeft, randomRight)
	assert.Equal(t, tree.Get(randomLeft, randomRight), fakeTree.Get(randomLeft, randomRight))
}

func Test_nextClosestPowerOf2(t *testing.T) {
	tests := []struct {
		x    int64
		want int64
	}{
		{1024, 1024},
		{1025, 2048},
		{1023, 1024},
		{0, 0},
		{1, 1},
		{2, 2},
		{3, 4},
	}
	for _, tt := range tests {
		t.Run("test", func(t *testing.T) {
			if got := nextClosestPowerOf2(tt.x); got != tt.want {
				t.Errorf("nextClosestPowerOf2() = %v, want %v", got, tt.want)
			}
		})
	}
}
