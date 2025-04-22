package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type KeyType interface {
	comparable
}

type ValueType interface {
	any
}

type Node[Key KeyType, Value ValueType] struct {
	key   Key
	value Value
	left  *Node[Key, Value]
	right *Node[Key, Value]
}

type OrderedMap[Key KeyType, Value ValueType] struct {
	root  *Node[Key, Value]
	index map[Key]*Node[Key, Value]
	less  func(a, b Key) bool
}

func NewOrderedMap[Key KeyType, Value ValueType](less func(a, b Key) bool) OrderedMap[Key, Value] {
	return OrderedMap[Key, Value]{
		index: make(map[Key]*Node[Key, Value]),
		less:  less,
	}
}

func (m *OrderedMap[Key, Value]) Insert(key Key, value Value) {
	n, exists := m.index[key]
	if exists {
		m.insert(n, key, value)
		return
	}

	if m.root == nil {
		m.root = &Node[Key, Value]{key: key, value: value}
		m.index[key] = m.root
		return
	}

	newNode := m.insert(m.root, key, value)
	m.index[key] = newNode
}

func (m *OrderedMap[Key, Value]) insert(node *Node[Key, Value], key Key, value Value) *Node[Key, Value] {
	if node == nil {
		return &Node[Key, Value]{key: key, value: value}
	}

	if m.less(key, node.key) {
		node.left = m.insert(node.left, key, value)
	} else if key != node.key {
		node.right = m.insert(node.right, key, value)
	}

	node.value = value
	return node
}

func (m *OrderedMap[Key, Value]) Erase(key Key) {
	_, exists := m.index[key]
	if !exists {
		return
	}

	m.root = m.remove(m.root, key)
	delete(m.index, key)
}

func (m *OrderedMap[Key, Value]) remove(node *Node[Key, Value], key Key) *Node[Key, Value] {
	if node == nil {
		return nil
	}

	if m.less(key, node.key) {
		node.left = m.remove(node.left, key)
	} else if key != node.key {
		node.right = m.remove(node.right, key)
	} else {
		if node.left == nil {
			return node.right
		} else if node.right == nil {
			return node.left
		} else {
			minRight := m.findMin(node.right)
			node.key = minRight.key
			node.value = minRight.value
			m.index[minRight.key] = node
			node.right = m.remove(node.right, minRight.key)
		}
	}

	return node
}

func (m *OrderedMap[Key, Value]) findMin(node *Node[Key, Value]) *Node[Key, Value] {
	current := node
	for current.left != nil {
		current = current.left
	}
	return current
}

func (m *OrderedMap[Key, Value]) Contains(key Key) bool {
	_, contains := m.index[key]
	return contains
}

func (m *OrderedMap[Key, Value]) Size() int {
	return len(m.index)
}

func (m *OrderedMap[Key, Value]) ForEach(action func(Key, Value)) {
	m.inOrder(m.root, action)
}

func (m *OrderedMap[Key, Value]) inOrder(node *Node[Key, Value], action func(Key, Value)) {
	if node == nil {
		return
	}

	m.inOrder(node.left, action)
	action(node.key, node.value)
	m.inOrder(node.right, action)
}

func TestCircularQueue(t *testing.T) {
	data := NewOrderedMap[int, int](func(a, b int) bool { return a < b })
	assert.Zero(t, data.Size())

	data.Insert(10, 10)
	data.Insert(5, 5)
	data.Insert(15, 15)
	data.Insert(2, 2)
	data.Insert(4, 4)
	data.Insert(12, 12)
	data.Insert(14, 14)

	assert.Equal(t, 7, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(3))
	assert.False(t, data.Contains(13))

	var keys []int
	expectedKeys := []int{2, 4, 5, 10, 12, 14, 15}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))

	data.Erase(15)
	data.Erase(14)
	data.Erase(2)

	assert.Equal(t, 4, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(2))
	assert.False(t, data.Contains(14))

	keys = nil
	expectedKeys = []int{4, 5, 10, 12}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))
}
