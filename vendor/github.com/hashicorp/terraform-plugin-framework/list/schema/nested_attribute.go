// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
)

// Nested attributes are only compatible with protocol version 6.
type NestedAttribute interface {
	Attribute
	fwschema.NestedAttribute
}
