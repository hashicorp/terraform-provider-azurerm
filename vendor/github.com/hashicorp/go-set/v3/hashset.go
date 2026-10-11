// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package set

import (
	"fmt"
	"iter"
	"sort"
)

// Hash represents the output type of a Hash() function defined on a type.
//
// A Hash could be string-like or int-like. A string hash could be something like
// and md5, sha1, or GoString() representation of a type. An int hash could be
// something like the prime multiple hash code of a type.
type Hash interface {
	~string | ~int | ~uint | ~int64 | ~uint64 | ~int32 | ~uint32 | ~int16 | ~uint16 | ~int8 | ~uint8
}

// Hasher represents a type that implements a Hash() method. Types that wish to
// cache a hash value with an internal field should implement Hash accordingly.
type Hasher[H Hash] interface {
	Hash() H
}

// HasherFunc creates a closure around the T.Hash function so that the type can
// be used as the HashFunc for a HashSet.
func HasherFunc[T Hasher[H], H Hash]() HashFunc[T, H] {
	return func(t T) H {
		return t.Hash()
	}
}

// HashFunc represents a function that that produces a hash value when applied
// to a given T. Typically this will be implemented as T.Hash but by separating
// HashFunc a HashSet can be made to make use of any hash implementation.
type HashFunc[T any, H Hash] func(T) H

// HashSet is a generic implementation of the mathematical data structure, oriented
// around the use of a HashFunc to make hash values from other types.
type HashSet[T any, H Hash] struct {
	fn    HashFunc[T, H]
	items map[H]T
}

// NewHashSet creates a HashSet with underlying capacity of size and will compute
// hash values from the T.Hash method.
func NewHashSet[T Hasher[H], H Hash](size int) *HashSet[T, H] {
	return NewHashSetFunc[T, H](size, HasherFunc[T, H]())
}

// NewHashSetFunc creates a HashSet with underlying capacity of size and uses
// the given hashing function to compute hashes on elements.
//
// A HashSet will automatically grow or shrink its capacity as items are added
// or removed.
func NewHashSetFunc[T any, H Hash](size int, fn HashFunc[T, H]) *HashSet[T, H] {
	return &HashSet[T, H]{
		fn:    fn,
		items: make(map[H]T, max(0, size)),
	}
}

// HashSetFrom creates a new HashSet containing each element in items.
//
// T must implement HashFunc[H], where H is of type Hash. This allows custom types
// that include non-comparable fields to provide their own hash algorithm.
func HashSetFrom[T Hasher[H], H Hash](items []T) *HashSet[T, H] {
	s := NewHashSet[T, H](len(items))
	s.InsertSlice(items)
	return s
}

// NewHashSetFromFunc creates a new HashSet containing each element in items.
func HashSetFromFunc[T any, H Hash](items []T, hash HashFunc[T, H]) *HashSet[T, H] {
	s := NewHashSetFunc[T, H](len(items), hash)
	s.InsertSlice(items)
	return s
}

// Insert item into s.
//
// Return true if s was modified (item was not already in s), false otherwise.
func (s *HashSet[T, H]) Insert(item T) bool {
	key := s.fn(item)
	if _, exists := s.items[key]; exists {
		return false
	}
	s.items[key] = item
	return true
}

// InsertSlice will insert each item in items into s.
//
// Return true if s was modified (at least one item was not already in s), false otherwise.
func (s *HashSet[T, H]) InsertSlice(items []T) bool {
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
func (s *HashSet[T, H]) InsertSet(col Collection[T]) bool {
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
func (s *HashSet[T, H]) Remove(item T) bool {
	key := s.fn(item)
	if _, exists := s.items[key]; !exists {
		return false
	}
	delete(s.items, key)
	return true
}

// RemoveSlice will remove each item in items from s.
//
// Return true if s was modified (any item was present), false otherwise.
func (s *HashSet[T, H]) RemoveSlice(items []T) bool {
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
// Return true if s was modified (any item of col was present in s), false otherwise.
func (s *HashSet[T, H]) RemoveSet(col Collection[T]) bool {
	return removeSet(s, col)
}

// RemoveFunc will remove each element from s that satisfies condition f.
//
// Return true if s was modified, false otherwise.
func (s *HashSet[T, H]) RemoveFunc(f func(item T) bool) bool {
	return removeFunc(s, f)
}

// Contains returns whether item is present in s.
func (s *HashSet[T, H]) Contains(item T) bool {
	hash := s.fn(item)
	_, exists := s.items[hash]
	return exists
}

// ContainsSlice returns true if all elements of items are in s, otherwise
// false.
func (s *HashSet[T, H]) ContainsSlice(items []T) bool {
	for _, item := range items {
		hash := s.fn(item)
		if _, exists := s.items[hash]; !exists {
			return false
		}
	}
	return true
}

// Subset returns whether col is a subset of s.
func (s *HashSet[T, H]) Subset(col Collection[T]) bool {
	return subset(s, col)
}

// ProperSubset returns whether col is a proper subset of s.
func (s *HashSet[T, H]) ProperSubset(col Collection[T]) bool {
	if len(s.items) <= col.Size() {
		return false
	}
	return s.Subset(col)
}

// Size returns the cardinality of s.
func (s *HashSet[T, H]) Size() int {
	return len(s.items)
}

// Empty returns true if s contains no elements, false otherwise.
func (s *HashSet[T, H]) Empty() bool {
	return s.Size() == 0
}

// Union returns a set that contains all elements of s and col combined.
//
// Elements in s take priority in the event of colliding hash values.
func (s *HashSet[T, H]) Union(col Collection[T]) Collection[T] {
	result := NewHashSetFunc[T, H](s.Size(), s.fn)
	insert(result, s)
	insert(result, col)
	return result
}

// Difference returns a set that contains elements of s that are not in col.
func (s *HashSet[T, H]) Difference(col Collection[T]) Collection[T] {
	result := NewHashSetFunc[T, H](max(0, s.Size()-col.Size()), s.fn)
	for item := range s.Items() {
		if !col.Contains(item) {
			result.Insert(item)
		}
	}
	return result
}

// Intersect returns a set that contains elements that are present in both s and col.
func (s *HashSet[T, H]) Intersect(col Collection[T]) Collection[T] {
	result := NewHashSetFunc[T, H](0, s.fn)
	intersect(result, s, col)
	return result
}

// Copy creates a shallow copy of s.
func (s *HashSet[T, H]) Copy() *HashSet[T, H] {
	result := NewHashSetFunc[T, H](s.Size(), s.fn)
	for key, item := range s.items {
		result.items[key] = item
	}
	return result
}

// Slice creates a copy of s as a slice.
//
// The result is not ordered.
func (s *HashSet[T, H]) Slice() []T {
	result := make([]T, 0, s.Size())
	for _, item := range s.items {
		result = append(result, item)
	}
	return result
}

// String creates a string representation of s, using "%v" printf formatting to transform
// each element into a string. The result contains elements sorted by their lexical
// string order.
func (s *HashSet[T, H]) String() string {
	return s.StringFunc(func(element T) string {
		return fmt.Sprintf("%v", element)
	})
}

// StringFunc creates a string representation of s, using f to transform each element
// into a string. The result contains elements sorted by their string order.
func (s *HashSet[T, H]) StringFunc(f func(element T) string) string {
	l := make([]string, 0, s.Size())
	for _, item := range s.items {
		l = append(l, f(item))
	}
	sort.Strings(l)
	return fmt.Sprintf("%s", l)
}

// Equal returns whether s and o contain the same elements.
func (s *HashSet[T, H]) Equal(o *HashSet[T, H]) bool {
	if len(s.items) != len(o.items) {
		return false
	}
	for _, item := range s.items {
		if !o.Contains(item) {
			return false
		}
	}
	return true
}

// EqualSet returns whether s and col contain the same elements.
func (s *HashSet[T, H]) EqualSet(col Collection[T]) bool {
	return equalSet(s, col)
}

// EqualSlice returns whether s and items contain the same elements.
//
// The items slice may contain duplicates.
//
// If the items slice is known to contain no duplicates, EqualSliceSet can be
// used instead as a faster implementation.
//
// To detect if a slice is a subset of s, use ContainsSlice.
func (s *HashSet[T, H]) EqualSlice(items []T) bool {
	other := HashSetFromFunc[T, H](items, s.fn)
	return s.Equal(other)
}

// EqualSliceSet returns whether s and items contain exactly the same elements.
//
// If items contains duplicates EqualSliceSet will return false. The elements of
// items are assumed to be set-like. For comparing s to a slice that may contain
// duplicate elements, use EqualSlice instead.
//
// To detect if a slice is a subset of s, use ContainsSlice.
func (s *HashSet[T, H]) EqualSliceSet(items []T) bool {
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
func (s *HashSet[T, H]) MarshalJSON() ([]byte, error) {
	return marshalJSON[T](s)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (s *HashSet[T, H]) UnmarshalJSON(data []byte) error {
	return unmarshalJSON[T](s, data)
}

// Items returns a generator function for iterating each element in s by using
// the range keyword.
//
//	for element := range s.Items() { ... }
func (s *HashSet[T, H]) Items() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, item := range s.items {
			if !yield(item) {
				return
			}
		}
	}
}
