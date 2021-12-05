package internal

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppend(t *testing.T) {
	l := []int{0, 1, 2, 3, 4}
	r := Append[int](l[2:], l[:3])
	assert.Equal(t, []int{2, 3, 4, 0, 1, 2}, r)
}

func TestRemoveIndex(t *testing.T) {
	l := []int{0, 1, 2, 3}

	assert.Equal(t, []int{1, 2, 3}, RemoveIndex[int](l, 0))
	assert.Equal(t, []int{0, 2, 3}, RemoveIndex[int](l, 1))
	assert.Equal(t, []int{0, 1, 2}, RemoveIndex[int](l, 3))
}

func TestGetOrCreate(t *testing.T) {
	m := map[string]int{"one": 1, "two": 2}

	assert.Equal(t, 2, GetOrCreate[string, int](m, "two", func(key string) int { return 0 }))
	assert.Equal(t, map[string]int{"one": 1, "two": 2}, m)

	assert.Equal(t, 3, GetOrCreate[string, int](m, "three", func(key string) int { return 3 }))
	assert.Equal(t, map[string]int{"one": 1, "two": 2, "three": 3}, m)
}

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

func TestFindValue(t *testing.T) {
	src := []int{0, 1, 2, 3, 4, 5, 6}

	assert.Equal(t, 3, FindValue[int](src, 3))
	assert.Equal(t, -1, FindValue[int](src, 100))
}

func TestReduce(t *testing.T) {
	src := []int{1, 2, 3, 4}

	sum := Reduce[int](src, 0, func(acc int, val int) int { return acc + val })
	assert.Equal(t, 10, sum)
}
