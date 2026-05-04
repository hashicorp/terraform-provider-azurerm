// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package defaults

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Float64 is a schema default value for types.Float64 attributes.
type Float64 interface {
	Describer

	// DefaultFloat64 should set the default value.
	DefaultFloat64(context.Context, Float64Request, *Float64Response)
}

type Float64Request struct {
	// Path contains the path of the attribute for setting the
	// default value. Use this path for any response diagnostics.
	Path path.Path
}

type Float64Response struct {
	// Diagnostics report errors or warnings related to setting the
	// default value resource configuration. An empty slice
	// indicates success, with no warnings or errors generated.
	Diagnostics diag.Diagnostics

	// PlanValue is the planned new state for the attribute.
	PlanValue types.Float64
}
