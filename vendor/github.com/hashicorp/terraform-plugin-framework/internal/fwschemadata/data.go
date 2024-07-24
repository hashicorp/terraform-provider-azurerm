// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwschemadata

import (
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Data is the shared storage implementation for schema-based values, such as
// configuration, plan, and state.
type Data struct {
	// Description contains the human friendly type of the data. Used in error
	// diagnostics.
	Description DataDescription

	// Schema contains the data structure and types for the value.
	Schema fwschema.Schema

	// TerraformValue contains the terraform-plugin-go value implementation.
	//
	// TODO: In the future this may be migrated to attr.Value, or more
	// succinctly, types.Object.
	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/172
	TerraformValue tftypes.Value
}
