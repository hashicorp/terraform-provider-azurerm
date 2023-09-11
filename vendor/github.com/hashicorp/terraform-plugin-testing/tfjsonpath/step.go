// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfjsonpath

// step represents a traversal type indicating the underlying Go type
// representation for a Terraform JSON value.
type step interface{}

// MapStep represents a traversal for map[string]any
type MapStep string

// SliceStep represents a traversal for []any
type SliceStep int
