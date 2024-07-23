// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package metaschema

import (
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
)

// Attribute define a value field inside the Schema. Implementations in this
// package include:
//   - BoolAttribute
//   - Float64Attribute
//   - Int64Attribute
//   - ListAttribute
//   - MapAttribute
//   - NumberAttribute
//   - ObjectAttribute
//   - SetAttribute
//   - StringAttribute
//
// Additionally, the NestedAttribute interface extends Attribute with nested
// attributes. Only supported in protocol version 6. Implementations in this
// package include:
//   - ListNestedAttribute
//   - MapNestedAttribute
//   - SetNestedAttribute
//   - SingleNestedAttribute
//
// In practitioner configurations, an equals sign (=) is required to set
// the value. [Configuration Reference]
//
// [Configuration Reference]: https://developer.hashicorp.com/terraform/language/syntax/configuration
type Attribute interface {
	fwschema.Attribute
}
