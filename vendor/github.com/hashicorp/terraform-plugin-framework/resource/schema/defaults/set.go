// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package defaults

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Set is a schema default value for types.Set attributes.
type Set interface {
	Describer

	// DefaultSet should set the default value.
	DefaultSet(context.Context, SetRequest, *SetResponse)
}

type SetRequest struct {
	// Path contains the path of the attribute for setting the
	// default value. Use this path for any response diagnostics.
	Path path.Path
}

type SetResponse struct {
	// Diagnostics report errors or warnings related to setting the
	// default value resource configuration. An empty slice
	// indicates success, with no warnings or errors generated.
	Diagnostics diag.Diagnostics

	// PlanValue is the planned new state for the attribute.
	PlanValue types.Set
}
