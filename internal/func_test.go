package internal

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	src := []int{1, 2, 3}
	dst := Map[int, string](src, func(item int) string { return fmt.Sprintf("[%d]", item) })
	assert.Equal(t, []string{"[1]", "[2]", "[3]"}, dst)
}

func TestFilter(t *testing.T) {
	src := []int{1, 2, 3, 4, 5, 6}
	dst := Filter[int](src, func(item int) bool { return (item % 2) == 0 })
	assert.Equal(t, []int{2, 4, 6}, dst)
}

func TestFind(t *testing.T) {
	src := []int{0, 1, 2, 3, 4, 5, 6}

	index := Find[int](src, func(item int) bool { return item > 3 })
	assert.Equal(t, 4, index)

	index = Find[int](src, func(item int) bool { return item > 100 })
	assert.Equal(t, -1, index)
}
