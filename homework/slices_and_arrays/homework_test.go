package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Int interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type CircularQueue[T Int] struct {
	values []T
	push   int
	pop    int
	cap    int
}

func NewCircularQueue[T Int](size int) CircularQueue[T] {
	data := make([]T, size)

	return CircularQueue[T]{
		values: data,
		cap:    size,
		push:   0,
		pop:    0,
	}
}

func (q *CircularQueue[T]) Push(value T) bool {
	if q.Full() {
		return false
	}

	q.values[q.push] = value
	q.push++

	q.cap--

	if q.push == len(q.values) {
		q.push = 0
	}

	return true
}

func (q *CircularQueue[T]) Pop() bool {
	if q.Empty() {
		return false
	}

	q.values[q.pop] = 0
	q.pop++

	q.cap++

	if q.pop == len(q.values) {
		q.pop = 0
	}

	return true
}

func (q *CircularQueue[T]) Front() T {
	if q.Empty() {
		return -1
	}

	return q.values[q.push]
}

func (q *CircularQueue[T]) Back() T {
	if q.Empty() {
		return -1
	}

	if q.pop == 0 {
		return q.values[len(q.values)-1]
	}

	return q.values[q.pop-1]
}

func (q *CircularQueue[T]) Empty() bool {
	return len(q.values) == q.cap
}

func (q *CircularQueue[T]) Full() bool {
	return q.cap == 0
}

func TestCircularQueue(t *testing.T) {
	const queueSize = 3
	queue := NewCircularQueue[int](queueSize)

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())

	assert.Equal(t, -1, queue.Front())
	assert.Equal(t, -1, queue.Back())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Push(1))
	assert.True(t, queue.Push(2))
	assert.True(t, queue.Push(3))
	assert.False(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int{1, 2, 3}, queue.values))

	assert.False(t, queue.Empty())
	assert.True(t, queue.Full())

	assert.Equal(t, 1, queue.Front())
	assert.Equal(t, 3, queue.Back())

	assert.True(t, queue.Pop())
	assert.False(t, queue.Empty())
	assert.False(t, queue.Full())
	assert.True(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int{4, 2, 3}, queue.values))

	assert.Equal(t, 2, queue.Front())
	assert.Equal(t, 4, queue.Back())

	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())
}
