// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package defaults

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Dynamic is a schema default value for types.Dynamic attributes.
type Dynamic interface {
	Describer

	// DefaultDynamic should set the default value.
	DefaultDynamic(context.Context, DynamicRequest, *DynamicResponse)
}

type DynamicRequest struct {
	// Path contains the path of the attribute for setting the
	// default value. Use this path for any response diagnostics.
	Path path.Path
}

type DynamicResponse struct {
	// Diagnostics report errors or warnings related to setting the
	// default value resource configuration. An empty slice
	// indicates success, with no warnings or errors generated.
	Diagnostics diag.Diagnostics

	// PlanValue is the planned new state for the attribute.
	PlanValue types.Dynamic
}
