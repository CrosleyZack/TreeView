package utils

// Queue is a generic queue type
type Queue[T any] []T

// NewQueue creates a new queue of type T
func NewQueue[T any]() Queue[T] {
	return make(Queue[T], 0)
}

func (s Queue[T]) Push(v T) Queue[T] {
	return append(s, v)
}

func (s Queue[T]) Pop() (Queue[T], T) {
	l := len(s)
	if l == 0 {
		return nil, *new(T)
	}
	return s[1:], s[0]
}
