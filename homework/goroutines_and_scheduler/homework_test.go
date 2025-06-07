package main

import (
	"container/heap"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Task struct {
	Identifier int
	Priority   int
}

type maxTaskPriorityHeap []Task

type Scheduler struct {
	h     maxTaskPriorityHeap
	index map[int]int
}

func (s Scheduler) Len() int {
	return len(s.h)
}

func (s Scheduler) Less(i, j int) bool {
	return s.h[i].Priority > s.h[j].Priority
}

func (s Scheduler) Swap(i, j int) {
	s.index[s.h[i].Identifier] = j
	s.index[s.h[j].Identifier] = i
	s.h[i], s.h[j] = s.h[j], s.h[i]
}

func (s *Scheduler) Push(val interface{}) {
	s.h = append(s.h, val.(Task))
}

func (s *Scheduler) Pop() interface{} {
	size := len(s.h)
	val := s.h[size-1]
	s.h = s.h[:size-1]

	delete(s.index, val.Identifier)

	return val
}

func NewScheduler() Scheduler {
	return Scheduler{maxTaskPriorityHeap{}, make(map[int]int)}
}

func (s *Scheduler) AddTask(task Task) {
	s.index[task.Identifier] = len(s.h)
	heap.Push(s, task)
}

func (s *Scheduler) ChangeTaskPriority(taskID int, newPriority int) {
	idx, ok := s.index[taskID]
	if !ok {
		return
	}

	s.h[idx].Priority = newPriority
	heap.Fix(s, idx)
}

func (s *Scheduler) GetTask() Task {
	return heap.Pop(s).(Task)
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
