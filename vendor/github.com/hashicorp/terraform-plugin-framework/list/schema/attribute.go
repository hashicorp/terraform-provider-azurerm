// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

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
//   - MapAttribute
//   - NumberAttribute
//   - StringAttribute
//
// In practitioner configurations, an equals sign (=) is required to set
// the value. [Configuration Reference]
//
// [Configuration Reference]: https://developer.hashicorp.com/terraform/language/syntax/configuration
type Attribute interface {
	fwschema.Attribute
}
