package tasks

import (
	"../models"

	"testing"
)

func TestCreateQueue(t *testing.T) {
	q := NewQueue(3)
	if q.Capacity() != 3 {
		t.Errorf("expected empty queue to have a specified capacity")
	}
	if q.Size() != 0 {
		t.Errorf("expected empty queue to have size 0")
	}
}

func TestPushQueue(t *testing.T) {
	q := NewQueue(2)
	err := q.Push(models.Monitor{})
	if err != nil {
		t.Errorf("expected to be able to push to a new queue %v", err)
	}
	if q.Size() != 1 {
		t.Errorf("expected queue to have size 1 after pushing an item")
	}
}

func TestPushOverflow(t *testing.T) {
	q := NewQueue(1)
	q.Push(models.Monitor{})
	err := q.Push(models.Monitor{})
	if err == nil {
		t.Errorf("expected to get an error once the queue's capacity is reached")
	}
}

func TestPopQueue(t *testing.T) {
	q := NewQueue(1)
	q.Push(models.Monitor{})
	_, err := q.Pop()
	if err != nil {
		t.Errorf("expected to be able to pop single item from queue with size 1")
	}
	if q.Size() != 0 {
		t.Errorf("expected popping to empty queue with size 1")
	}
}

func TestPopUnderflow(t *testing.T) {
	q := NewQueue(1)
	_, err := q.Pop()
	if err == nil {
		t.Errorf("expected to get an error trying to pop from an empty queue")
	}
}

func TestOrderOfOperations(t *testing.T) {
	q := NewQueue(2)
	q.Push(models.Monitor{})
	q.Push(models.Monitor{})
	q.Pop()
	err := q.Push(models.Monitor{})
	if err != nil {
		t.Errorf("expected to be able to push to a queue after popping")
	}
	q.Pop()
	err = q.Push(models.Monitor{})
	if err != nil {
		t.Errorf("expected to be able to push to a queue after popping")
	}
}
