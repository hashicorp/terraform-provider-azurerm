// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package models

// SchemaProperty represents a property parsed from Terraform schema (Go code)
type SchemaProperty struct {
	Name        string
	Type        string
	Description string
	Required    bool
	Optional    bool
	Computed    bool
	ForceNew    bool
	Deprecated  bool

	PossibleValues []string
	DefaultValue   interface{}

	// Block related attributes
	Nested     *SchemaProperties
	Block      bool
	NestedType string
}

type SchemaProperties struct {
	Objects map[string]*SchemaProperty
}

func NewSchemaProperties() *SchemaProperties {
	return &SchemaProperties{
		Objects: make(map[string]*SchemaProperty),
	}
}
