// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package set

type stack[T any] struct {
	top *object[T]
}

type object[T any] struct {
	item T
	next *object[T]
}

func makeStack[T any]() *stack[T] {
	return new(stack[T])
}

func (s *stack[T]) push(item T) {
	obj := &object[T]{
		item: item,
		next: s.top,
	}
	s.top = obj
}

func (s *stack[T]) pop() T {
	obj := s.top
	s.top = obj.next
	obj.next = nil
	return obj.item
}

func (s *stack[T]) empty() bool {
	return s.top == nil
}
