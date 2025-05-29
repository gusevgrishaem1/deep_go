package main

import (
	"container/heap"
	"testing"

	"github.com/stretchr/testify/assert"
)

type maxTaskPriorityHeap []Task

func (h maxTaskPriorityHeap) Len() int {
	return len(h)
}

func (h maxTaskPriorityHeap) Less(i, j int) bool {
	return h[i].Priority > h[j].Priority
}

func (h maxTaskPriorityHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *maxTaskPriorityHeap) Push(val interface{}) {
	*h = append(*h, val.(Task))
}

func (h *maxTaskPriorityHeap) Pop() interface{} {
	elements := *h

	size := len(elements)
	val := elements[size-1]
	*h = elements[:size-1]

	return val
}

type Task struct {
	Identifier int
	Priority   int
}

type Scheduler struct {
	h maxTaskPriorityHeap
}

func NewScheduler() Scheduler {
	return Scheduler{maxTaskPriorityHeap{}}
}

func (s *Scheduler) AddTask(task Task) {
	heap.Push(&s.h, task)
}

func (s *Scheduler) ChangeTaskPriority(taskID int, newPriority int) {
	for i, task := range s.h {
		if task.Identifier == taskID {
			s.h[i].Priority = newPriority
			heap.Fix(&s.h, i)
			return
		}
	}
}

func (s *Scheduler) GetTask() Task {
	return heap.Pop(&s.h).(Task)
}

func TestTrace(t *testing.T) {
	task1 := Task{Identifier: 1, Priority: 10}
	task2 := Task{Identifier: 2, Priority: 20}
	task3 := Task{Identifier: 3, Priority: 30}
	task4 := Task{Identifier: 4, Priority: 40}
	task5 := Task{Identifier: 5, Priority: 50}

	scheduler := NewScheduler()
	scheduler.AddTask(task1)
	scheduler.AddTask(task2)
	scheduler.AddTask(task3)
	scheduler.AddTask(task4)
	scheduler.AddTask(task5)

	task := scheduler.GetTask()
	assert.Equal(t, task5, task)

	task = scheduler.GetTask()
	assert.Equal(t, task4, task)

	scheduler.ChangeTaskPriority(1, 100)
	task1.Priority = 100

	task = scheduler.GetTask()
	assert.Equal(t, task1, task)

	task = scheduler.GetTask()
	assert.Equal(t, task3, task)
}
