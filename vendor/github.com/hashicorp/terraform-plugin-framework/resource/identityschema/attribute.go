// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package identityschema

import (
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
)

// Attribute define a value field inside the Schema. Implementations in this
// package include:
//   - BoolAttribute
//   - Float32Attribute
//   - Float64Attribute
//   - Int32Attribute
//   - Int64Attribute
//   - ListAttribute
//   - NumberAttribute
//   - StringAttribute
//
// The available attribute types for a resource identity schema are intentionally
// limited. Nested attributes and blocks are not supported in identity schemas,
// as well as ListAttribute definitions can only have primitive element types of:
//   - types.BoolType
//   - types.Float32Type
//   - types.Float64Type
//   - types.Int32Type
//   - types.Int64Type
//   - types.NumberType
//   - types.StringType
type Attribute interface {
	fwschema.Attribute
}
