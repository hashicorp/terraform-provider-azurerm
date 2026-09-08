// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package set provides a basic generic set implementation.
//
// https://en.wikipedia.org/wiki/Set_(mathematics)
package set

import (
	"fmt"
	"iter"
	"sort"
)

type nothing struct{}

var sentinel = nothing{}

// New creates a new Set with initial underlying capacity of size.
//
// A Set will automatically grow or shrink its capacity as items are added or
// removed.
//
// T may be any comparable type. Keep in mind that pointer types or structs
// containing pointer fields will be compared using shallow equality. For deep
// equality use HashSet instead.
func New[T comparable](size int) *Set[T] {
	return &Set[T]{
		items: make(map[T]nothing, max(0, size)),
	}
}

// From creates a new Set containing each item in items.
//
// T may be any comparable type. Keep in mind that pointer types or structs
// containing pointer fields will be compared using shallow equality. For deep
// equality use HashSet instead.
func From[T comparable](items []T) *Set[T] {
	s := New[T](len(items))
	s.InsertSlice(items)
	return s
}

// FromFunc creates a new Set containing a conversion of each item in items.
//
// T may be any comparable type. Keep in mind that pointer types or structs
// containing pointer fields will be compared using shallow equality. For deep
// equality use HashSet instead.
func FromFunc[A any, T comparable](items []A, conversion func(A) T) *Set[T] {
	s := New[T](len(items))
	for _, item := range items {
		s.Insert(conversion(item))
	}
	return s
}

// Set is a simple, generic implementation of the set mathematical data structure.
// It is optimized for correctness and convenience, as a replacement for the use
// of map[interface{}]struct{}.
type Set[T comparable] struct {
	items map[T]nothing
}

// Insert item into s.
//
// Return true if s was modified (item was not already in s), false otherwise.
func (s *Set[T]) Insert(item T) bool {
	if _, exists := s.items[item]; exists {
		return false
	}
	if s.items == nil {
		s.items = make(map[T]nothing)
	}
	s.items[item] = sentinel
	return true
}

// InsertSlice will insert each item in items into s.
//
// Return true if s was modified (at least one item was not already in s), false otherwise.
func (s *Set[T]) InsertSlice(items []T) bool {
	modified := false
	for _, item := range items {
		if s.Insert(item) {
			modified = true
		}
	}
	return modified
}

// InsertSet will insert each element of col into s.
//
// Return true if s was modified (at least one item of col was not already in s), false otherwise.
func (s *Set[T]) InsertSet(col Collection[T]) bool {
	modified := false
	for item := range col.Items() {
		if s.Insert(item) {
			modified = true
		}
	}
	return modified
}

// Remove will remove item from s.
//
// Return true if s was modified (item was present), false otherwise.
func (s *Set[T]) Remove(item T) bool {
	if _, exists := s.items[item]; !exists {
		return false
	}
	delete(s.items, item)
	return true
}

// RemoveSlice will remove each item in items from s.
//
// Return true if s was modified (any item was present), false otherwise.
func (s *Set[T]) RemoveSlice(items []T) bool {
	modified := false
	for _, item := range items {
		if s.Remove(item) {
			modified = true
		}
	}
	return modified
}

// RemoveSet will remove each element of col from s.
//
// Return true if s was modified (any item of o was present in s), false otherwise.
func (s *Set[T]) RemoveSet(col Collection[T]) bool {
	return removeSet(s, col)
}

// RemoveFunc will remove each element from s that satisfies condition f.
//
// Return true if s was modified, false otherwise.
func (s *Set[T]) RemoveFunc(f func(T) bool) bool {
	return removeFunc(s, f)
}

// Contains returns whether item is present in s.
func (s *Set[T]) Contains(item T) bool {
	_, exists := s.items[item]
	return exists
}

// ContainsSlice returns whether all elements in items are present in s.
func (s *Set[T]) ContainsSlice(items []T) bool {
	return containsSlice(s, items)
}

// Subset returns whether col is a subset of s.
func (s *Set[T]) Subset(col Collection[T]) bool {
	return subset(s, col)
}

// Subset returns whether col is a proper subset of s.
func (s *Set[T]) ProperSubset(col Collection[T]) bool {
	if len(s.items) <= col.Size() {
		return false
	}
	return s.Subset(col)
}

// Size returns the cardinality of s.
func (s *Set[T]) Size() int {
	return len(s.items)
}

// Empty returns true if s contains no elements, false otherwise.
func (s *Set[T]) Empty() bool {
	return s.Size() == 0
}

// Union returns a set that contains all elements of s and col combined.
func (s *Set[T]) Union(col Collection[T]) Collection[T] {
	size := max(s.Size(), col.Size())
	result := New[T](size)
	insert(result, s)
	insert(result, col)
	return result
}

// Difference returns a set that contains elements of s that are not in col.
func (s *Set[T]) Difference(col Collection[T]) Collection[T] {
	result := New[T](max(0, s.Size()-col.Size()))
	for item := range s.items {
		if !col.Contains(item) {
			result.items[item] = sentinel
		}
	}
	return result
}

// Intersect returns a set that contains elements that are present in both s and col.
func (s *Set[T]) Intersect(col Collection[T]) Collection[T] {
	result := New[T](0)
	intersect(result, s, col)
	return result
}

// Copy creates a copy of s.
func (s *Set[T]) Copy() *Set[T] {
	result := New[T](s.Size())
	for item := range s.items {
		result.items[item] = sentinel
	}
	return result
}

// Slice creates a copy of s as a slice. Elements are in no particular order.
func (s *Set[T]) Slice() []T {
	result := make([]T, 0, s.Size())
	for item := range s.items {
		result = append(result, item)
	}
	return result
}

// String creates a string representation of s, using "%v" printf formating to transform
// each element into a string. The result contains elements sorted by their lexical
// string order.
func (s *Set[T]) String() string {
	return s.StringFunc(func(element T) string {
		return fmt.Sprintf("%v", element)
	})
}

// StringFunc creates a string representation of s, using f to transform each element
// into a string. The result contains elements sorted by their lexical string order.
func (s *Set[T]) StringFunc(f func(element T) string) string {
	l := make([]string, 0, s.Size())
	for item := range s.items {
		l = append(l, f(item))
	}
	sort.Strings(l)
	return fmt.Sprintf("%s", l)
}

// Equal returns whether s and o contain the same elements.
func (s *Set[T]) Equal(o *Set[T]) bool {
	if len(s.items) != len(o.items) {
		return false
	}
	for item := range s.items {
		if !o.Contains(item) {
			return false
		}
	}
	return true
}

// EqualSet returns whether s and col contain the same elements.
func (s *Set[T]) EqualSet(col Collection[T]) bool {
	return equalSet(s, col)
}

// EqualSlice returns whether s and items contain the same elements.
//
// The items slice may contain duplicates.
//
// If the items slice is known to contain no duplicates, EqualSliceSet may be
// used instead as a faster implementation.
//
// To detect if a slice is a subset of s, use ContainsSlice.
func (s *Set[T]) EqualSlice(items []T) bool {
	other := From[T](items)
	return s.Equal(other)
}

// EqualSliceSet returns whether s and items contain exactly the same elements.
//
// If items contains duplicates EqualSliceSet will return false. The elements of
// items are assumed to be set-like. For comparing s to a slice that may contain
// duplicate elements, use EqualSlice instead.
//
// To detect if a slice is a subset of s, use ContainsSlice.
func (s *Set[T]) EqualSliceSet(items []T) bool {
	if len(items) != s.Size() {
		return false
	}
	for _, item := range items {
		if !s.Contains(item) {
			return false
		}
	}
	return true
}

// MarshalJSON implements the json.Marshaler interface.
func (s *Set[T]) MarshalJSON() ([]byte, error) {
	return marshalJSON[T](s)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (s *Set[T]) UnmarshalJSON(data []byte) error {
	return unmarshalJSON[T](s, data)
}

// Items returns a generator function for iterating each element in s by using
// the range keyword.
//
//	for element := range s.Items() { ... }
func (s *Set[T]) Items() iter.Seq[T] {
	return func(yield func(T) bool) {
		for item := range s.items {
			if !yield(item) {
				return
			}
		}
	}
}
