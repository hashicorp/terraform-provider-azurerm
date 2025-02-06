// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwschema

// NestingMode is an enum type of the ways nested attributes can be nested in
// an attribute. They can be a list, a set, a map (with string
// keys), or they can be nested directly, like an object.
type NestingMode uint8

const (
	// NestingModeUnknown is an invalid nesting mode, used to catch when a
	// nesting mode is expected and not set.
	NestingModeUnknown NestingMode = 0

	// NestingModeSingle is for attributes that represent a struct or
	// object, a single instance of those attributes directly nested under
	// another attribute.
	NestingModeSingle NestingMode = 1

	// NestingModeList is for attributes that represent a list of objects,
	// with multiple instances of those attributes nested inside a list
	// under another attribute.
	NestingModeList NestingMode = 2

	// NestingModeSet is for attributes that represent a set of objects,
	// with multiple, unique instances of those attributes nested inside a
	// set under another attribute.
	NestingModeSet NestingMode = 3

	// NestingModeMap is for attributes that represent a map of objects,
	// with multiple instances of those attributes, each associated with a
	// unique string key, nested inside a map under another attribute.
	NestingModeMap NestingMode = 4
)
