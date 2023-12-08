// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfjsonpath

import (
	"fmt"
)

// Path represents exact traversal steps specifying a value inside
// Terraform JSON data. These steps always start from a MapStep with a key
// specifying the name of a top-level JSON object or array.
//
// The [terraform-json] library serves as the de facto documentation
// for JSON format of Terraform data.
//
// Use the New() function to create a Path with an initial AtMapKey() step.
// Path functionality follows a builder pattern, which allows for chaining method
// calls to construct a full path. The available traversal steps after Path
// creation are:
//
//   - AtSliceIndex(): Step into a slice at a specific 0-based index
//   - AtMapKey(): Step into a map at a specific key
//
// For example, to represent the first element of a JSON array
// underneath a "some_array" property of this JSON value:
//
//	   {
//	     "some_array": [true]
//	   }
//
//	 The path code would be represented by:
//
//		tfjsonpath.New("some_array").AtSliceIndex(0)
//
// [terraform-json]: (https://pkg.go.dev/github.com/hashicorp/terraform-json)
type Path struct {
	steps []step
}

// New creates a new path with an initial MapStep.
func New(name string) Path {
	return Path{
		steps: []step{
			MapStep(name),
		},
	}
}

// AtSliceIndex returns a copied Path with a new SliceStep at the end.
func (s Path) AtSliceIndex(index int) Path {
	newSteps := append(s.steps, SliceStep(index))
	s.steps = newSteps
	return s
}

// AtMapKey returns a copied Path with a new MapStep at the end.
func (s Path) AtMapKey(key string) Path {
	newSteps := append(s.steps, MapStep(key))
	s.steps = newSteps
	return s
}

// Traverse returns the element found when traversing the given
// object using the specified Path. The object is an unmarshalled
// JSON object representing Terraform data.
//
// Traverse returns an error if the value specified by the Path
// is not found in the given object or if the given object does not
// conform to format of Terraform JSON data.
func Traverse(object any, attrPath Path) (any, error) {
	_, ok := object.(map[string]any)

	if !ok {
		return nil, fmt.Errorf("cannot convert given object to map[string]any")
	}

	result := object

	for _, step := range attrPath.steps {
		switch s := step.(type) {
		case MapStep:
			mapObj, ok := result.(map[string]any)
			if !ok {
				return nil, fmt.Errorf("path not found: cannot convert object at MapStep %s to map[string]any", string(s))
			}
			result, ok = mapObj[string(s)]
			if !ok {
				return nil, fmt.Errorf("path not found: specified key %s not found in map", string(s))
			}

		case SliceStep:
			sliceObj, ok := result.([]any)
			if !ok {
				return nil, fmt.Errorf("path not found: cannot convert object at SliceStep %d to []any", s)
			}

			if int(s) >= len(sliceObj) {
				return nil, fmt.Errorf("path not found: SliceStep index %d is out of range with slice length %d", s, len(sliceObj))
			}

			result = sliceObj[s]
		}
	}

	return result, nil
}
