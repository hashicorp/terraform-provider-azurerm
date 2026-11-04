// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
)

// Attribute define a value field inside a state store type schema. Implementations in this
// package include:
//   - BoolAttribute
//   - DynamicAttribute
//   - Float32Attribute
//   - Float64Attribute
//   - Int32Attribute
//   - Int64Attribute
//   - ListAttribute
//   - MapAttribute
//   - NumberAttribute
//   - ObjectAttribute
//   - SetAttribute
//   - StringAttribute
//
// Additionally, the NestedAttribute interface extends Attribute with nested
// attributes, implementations in this package include:
//   - ListNestedAttribute
//   - MapNestedAttribute
//   - SetNestedAttribute
//   - SingleNestedAttribute
type Attribute interface {
	fwschema.Attribute
}

// schemaAttributes is a state store attribute to fwschema type conversion function.
func schemaAttributes(attributes map[string]Attribute) map[string]fwschema.Attribute {
	result := make(map[string]fwschema.Attribute, len(attributes))

	for name, attribute := range attributes {
		result[name] = attribute
	}

	return result
}
