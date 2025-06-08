package main

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

func Trace(stacks [][]uintptr) []uintptr {
	type el struct {
		ptr uintptr
		pos int
	}

	var result []uintptr
	visited := make(map[uintptr]struct{})
	var queue []el

	// Первый проход — собираем все уникальные указатели из стеков
	for _, stack := range stacks {
		for _, ptr := range stack {
			if ptr == 0 {
				continue
			}
			if _, ok := visited[ptr]; ok {
				continue
			}
			visited[ptr] = struct{}{}
			pos := len(result)
			result = append(result, ptr)

			next := *(*uintptr)(unsafe.Pointer(ptr))
			if next != 0 {
				queue = append(queue, el{next, pos + 1})
			}
		}
	}

	// BFS
	for len(queue) > 0 {
		item := queue[0]
		queue = queue[1:]

		ptr := item.ptr
		if ptr == 0 {
			continue
		}
		if _, ok := visited[ptr]; ok {
			continue
		}
		visited[ptr] = struct{}{}

		// Вставка в нужную позицию
		if item.pos >= len(result) {
			result = append(result, ptr)
		} else {
			result = append(result[:item.pos], append([]uintptr{ptr}, result[item.pos:]...)...)
		}

		next := *(*uintptr)(unsafe.Pointer(ptr))
		if next != 0 {
			queue = append(queue, el{next, item.pos + 1})
		}
	}

	return result
}

func TestTrace(t *testing.T) {
	var heapObjects = []int{
		0x00, 0x00, 0x00, 0x00, 0x00,
	}

	var heapPointer1 *int = &heapObjects[1]
	var heapPointer2 *int = &heapObjects[2]
	var heapPointer3 *int = nil
	var heapPointer4 **int = &heapPointer3

	var stacks = [][]uintptr{
		{
			uintptr(unsafe.Pointer(&heapPointer1)), 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, uintptr(unsafe.Pointer(&heapObjects[0])),
			0x00, 0x00, 0x00, 0x00,
		},
		{
			uintptr(unsafe.Pointer(&heapPointer2)), 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, uintptr(unsafe.Pointer(&heapObjects[1])),
			0x00, 0x00, 0x00, uintptr(unsafe.Pointer(&heapObjects[2])),
			uintptr(unsafe.Pointer(&heapPointer4)), 0x00, 0x00, 0x00,
		},
		{
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, uintptr(unsafe.Pointer(&heapObjects[3])),
		},
	}

	pointers := Trace(stacks)
	expectedPointers := []uintptr{
		uintptr(unsafe.Pointer(&heapPointer1)),
		uintptr(unsafe.Pointer(&heapObjects[0])),
		uintptr(unsafe.Pointer(&heapPointer2)),
		uintptr(unsafe.Pointer(&heapObjects[1])),
		uintptr(unsafe.Pointer(&heapObjects[2])),
		uintptr(unsafe.Pointer(&heapPointer4)),
		uintptr(unsafe.Pointer(&heapPointer3)),
		uintptr(unsafe.Pointer(&heapObjects[3])),
	}

	assert.True(t, reflect.DeepEqual(expectedPointers, pointers))
}
