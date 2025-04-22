package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Node struct {
	key   int
	value int
	left  *Node
	right *Node
}

type OrderedMap struct {
	root  *Node
	index map[int]*Node
}

func NewOrderedMap() OrderedMap {
	return OrderedMap{
		index: make(map[int]*Node),
	}
}

func (m *OrderedMap) Insert(key, value int) {
	n, exists := m.index[key]
	if exists {
		insert(n, key, value)
		return
	}

	if m.root == nil {
		m.root = &Node{key: key, value: value}
		m.index[key] = m.root
		return
	}

	newNode := insert(m.root, key, value)
	m.index[key] = newNode
}

func insert(root *Node, key int, value int) *Node {
	if root == nil {
		return &Node{key: key, value: value}
	}

	if key < root.key {
		root.left = insert(root.left, key, value)
	} else if key > root.key {
		root.right = insert(root.right, key, value)
	}

	root.value = value
	return root
}

func findMin(node *Node) *Node {
	current := node
	for current.left != nil {
		current = current.left
	}
	return current
}

func (m *OrderedMap) Erase(key int) {
	_, exists := m.index[key]
	if !exists {
		return
	}

	m.root = m.remove(m.root, key)
	delete(m.index, key)
}

func (m *OrderedMap) remove(root *Node, key int) *Node {
	if root == nil {
		return nil
	}

	if key < root.key {
		root.left = m.remove(root.left, key)
	} else if key > root.key {
		root.right = m.remove(root.right, key)
	} else {
		if root.left == nil {
			return root.right
		} else if root.right == nil {
			return root.left
		} else {
			minRight := findMin(root.right)
			root.key = minRight.key
			root.value = minRight.value
			m.index[minRight.key] = root
			root.right = m.remove(root.right, minRight.key)
		}
	}

	return root
}

func (m *OrderedMap) Contains(key int) bool {
	_, contains := m.index[key]
	return contains
}

func (m *OrderedMap) Size() int {
	return len(m.index)
}

func (m *OrderedMap) ForEach(action func(int, int)) {
	inOrder(m.root, action)
}

func inOrder(n *Node, action func(int, int)) {
	if n == nil {
		return
	}

	inOrder(n.left, action)
	action(n.key, n.value)
	inOrder(n.right, action)
}

func TestCircularQueue(t *testing.T) {
	data := NewOrderedMap()
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
