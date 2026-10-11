// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
)

type NestedAttribute interface {
	Attribute
	fwschema.NestedAttribute
}
