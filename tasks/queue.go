package tasks

import (
	"../models"

	"errors"
)

// Queue contains monitors to run to check sites for changes in a First-In,
// First-Out order.
type Queue struct {
	size      uint
	capacity  uint
	container []models.Monitor
}

// NewQueue constructs a Queue with a given capacity.
func NewQueue(cap uint) Queue {
	return Queue{
		size:      0,
		capacity:  cap,
		container: make([]models.Monitor, cap),
	}
}

// Capacity gets the max capacity of the Queue.
func (q Queue) Capacity() uint {
	return q.capacity
}

// Size gets the number of items current in the queue.
func (q Queue) Size() uint {
	return q.size
}

// Push inserts a new Monitor into a free position in the queue.
func (q *Queue) Push(m models.Monitor) error {
	if q.size == q.capacity {
		return errors.New("queue is full")
	}
	q.container[q.size] = m
	q.size++
	return nil
}

// Pop removes and returns the element at the front of the queue.
// It also shifts all elements forward rather than having the elements in
// the container array wrap around its bounds.
func (q *Queue) Pop() (models.Monitor, error) {
	if q.size == 0 {
		return models.Monitor{}, errors.New("queue is empty")
	}
	monitor := q.container[0]
	var i uint = 1
	for ; i < q.size; i++ {
		q.container[i-1] = q.container[i]
	}
	q.size--
	return monitor, nil
}
