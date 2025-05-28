package main

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type el struct {
	uptr uintptr
	idx  int
}

func Trace(stacks [][]uintptr) []uintptr {
	result := make([]uintptr, 0, len(stacks))
	visited := make(map[uintptr]struct{}, len(stacks))
	var queue []el

	for i := 0; i < len(stacks); i++ {
		for j := 0; j < len(stacks[i]); j++ {
			if stacks[i][j] == 0 {
				continue
			}

			if _, ok := visited[stacks[i][j]]; ok || stacks[i][j] == 0 {
				continue
			}

			result = append(result, stacks[i][j])
			visited[stacks[i][j]] = struct{}{}

			next := (*uintptr)(unsafe.Pointer(stacks[i][j]))
			if _, ok := visited[*next]; !ok && *next != 0 {
				queue = append(queue, el{*next, len(result)})
			}
		}
	}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		if _, ok := visited[current.uptr]; ok || current.uptr == 0 {
			continue
		}

		tmp := make([]uintptr, len(result)-current.idx)
		copy(tmp, result[current.idx:])
		result = result[:current.idx]
		result = append(result, current.uptr)
		result = append(result, tmp...)

		next := (*uintptr)(unsafe.Pointer(current.uptr))
		if _, ok := visited[*next]; !ok && *next != 0 {
			queue = append(queue, el{*next, current.idx})
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
