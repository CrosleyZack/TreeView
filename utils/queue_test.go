package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStack(t *testing.T) {
	s := NewQueue[int]()
	assert.Len(t, s, 0)
	s = s.Push(1)
	assert.Len(t, s, 1)
	s = s.Push(2)
	assert.Len(t, s, 2)
	s, v := s.Pop()
	assert.Len(t, s, 1)
	assert.Equal(t, 1, v)
	s, v = s.Pop()
	assert.Len(t, s, 0)
	assert.Equal(t, 2, v)
	s, v = s.Pop()
	assert.Len(t, s, 0)
	assert.Equal(t, 0, v) // default value for int
}
