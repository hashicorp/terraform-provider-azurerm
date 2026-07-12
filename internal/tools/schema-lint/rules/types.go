// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package rules

import "github.com/hashicorp/terraform-provider-azurerm/internal/tools/schema-api/providerjson"

// Schema value type strings as produced by (schema.ValueType).String() in the
// terraform-plugin-sdk.
//
// NOTE: providerjson also defines type constants, but some of them (e.g.
// SchemaTypeString = "String") do not match the actual (schema.ValueType).String()
// output. The linter therefore defines and uses these corrected values.
const (
	TypeBool   = "TypeBool"
	TypeInt    = "TypeInt"
	TypeFloat  = "TypeFloat"
	TypeString = "TypeString"
	TypeList   = "TypeList"
	TypeMap    = "TypeMap"
	TypeSet    = "TypeSet"
)

// IsCollection reports whether the schema node is a list or a set.
func IsCollection(s providerjson.SchemaJSON) bool {
	return s.Type == TypeList || s.Type == TypeSet
}

// BlockElem returns the nested resource schema for a block-typed node (a
// TypeList or TypeSet whose Elem is a resource), and true when the node is a
// nested block.
//
// When the schema is loaded from the live provider, Elem is a
// *providerjson.ResourceJSON; a value providerjson.ResourceJSON (as produced by
// JSON unmarshalling) is also handled for safety.
func BlockElem(s providerjson.SchemaJSON) (providerjson.ResourceJSON, bool) {
	if !IsCollection(s) {
		return providerjson.ResourceJSON{}, false
	}

	switch e := s.Elem.(type) {
	case *providerjson.ResourceJSON:
		if e != nil {
			return *e, true
		}
	case providerjson.ResourceJSON:
		return e, true
	}

	return providerjson.ResourceJSON{}, false
}
