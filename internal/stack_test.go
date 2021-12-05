package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPush(t *testing.T) {
	var s []int
	Push(&s, 1)
	assert.Equal(t, s, []int{1})
}
