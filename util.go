package main

import "errors"

type Queue[T any] struct {
	data []T
}

func (queue *Queue[T]) pop() (*T, error) {
	if len(queue.data) == 0 {
		return nil, errors.New("no items in queue")
	}
	itemToPop := queue.data[0]
	(*queue).data = queue.data[1:]
	return &itemToPop, nil
}

func (queue *Queue[T]) push(stack []T, item T) []T {
	stack = append(stack, item)
	return stack
}

func (queue *Queue[T]) isEmpty() bool {
	return len(queue.data) == 0
}
