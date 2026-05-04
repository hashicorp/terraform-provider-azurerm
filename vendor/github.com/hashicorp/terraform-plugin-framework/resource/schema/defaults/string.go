// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package defaults

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// String is a schema default value for types.String attributes.
type String interface {
	Describer

	// DefaultString should set the default value.
	DefaultString(context.Context, StringRequest, *StringResponse)
}

type StringRequest struct {
	// Path contains the path of the attribute for setting the
	// default value. Use this path for any response diagnostics.
	Path path.Path
}

type StringResponse struct {
	// Diagnostics report errors or warnings related to setting the
	// default value resource configuration. An empty slice
	// indicates success, with no warnings or errors generated.
	Diagnostics diag.Diagnostics

	// PlanValue is the planned new state for the attribute.
	PlanValue types.String
}
