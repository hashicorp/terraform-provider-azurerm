// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwschema

// BlockNestingMode is an enum type of the ways attributes and blocks can be
// nested in a block. They can be a list or a set.
//
// While the protocol and theoretically Terraform itself support map and group
// nesting modes, this framework intentionally only supports list, set, and
// single blocks as those other modes were not typically implemented or
// tested with Terraform since the older Terraform Plugin SDK did not support
// them.
type BlockNestingMode uint8

const (
	// BlockNestingModeUnknown is an invalid nesting mode, used to catch when a
	// nesting mode is expected and not set.
	BlockNestingModeUnknown BlockNestingMode = 0

	// BlockNestingModeList is for attributes that represent a list of objects,
	// with multiple instances of those attributes nested inside a list
	// under another attribute.
	BlockNestingModeList BlockNestingMode = 1

	// BlockNestingModeSet is for attributes that represent a set of objects,
	// with multiple, unique instances of those attributes nested inside a
	// set under another attribute.
	BlockNestingModeSet BlockNestingMode = 2

	// BlockNestingModeSingle is for attributes that represent a single object.
	// The object cannot be repeated in the practitioner configuration.
	//
	// While the framework implements support for this block nesting mode, it
	// is not thoroughly tested in production Terraform environments beyond the
	// resource timeouts block from the older Terraform Plugin SDK. Use single
	// nested attributes for new implementations instead.
	BlockNestingModeSingle BlockNestingMode = 3
)
