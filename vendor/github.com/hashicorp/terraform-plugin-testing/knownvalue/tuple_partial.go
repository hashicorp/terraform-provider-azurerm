// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

var _ Check = tuplePartial{}

type tuplePartial struct {
	value map[int]Check
}

// CheckValue determines whether the passed value is of type []any, and
// contains matching slice entries in the same sequence.
func (v tuplePartial) CheckValue(other any) error {
	otherVal, ok := other.([]any)

	if !ok {
		return fmt.Errorf("expected []any value for TuplePartial check, got: %T", other)
	}

	var keys []int

	for k := range v.value {
		keys = append(keys, k)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	for _, k := range keys {
		if len(otherVal) <= k {
			return fmt.Errorf("missing element index %d for TuplePartial check", k)
		}

		if err := v.value[k].CheckValue(otherVal[k]); err != nil {
			return fmt.Errorf("tuple element %d: %s", k, err)
		}
	}

	return nil
}

// String returns the string representation of the value.
func (v tuplePartial) String() string {
	var b bytes.Buffer

	b.WriteString("[")

	var keys []int

	var tupleVals []string

	for k := range v.value {
		keys = append(keys, k)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	for _, k := range keys {
		tupleVals = append(tupleVals, fmt.Sprintf("%d:%s", k, v.value[k]))
	}

	b.WriteString(strings.Join(tupleVals, " "))

	b.WriteString("]")

	return b.String()
}

// TuplePartial returns a Check for asserting partial equality between the
// supplied map[int]Check and the value passed to the CheckValue method. The
// map keys represent the zero-ordered element indices within the tuple that is
// being checked. Only the elements at the indices defined within the
// supplied map[int]Check are checked.
func TuplePartial(value map[int]Check) tuplePartial {
	return tuplePartial{
		value: value,
	}
}
