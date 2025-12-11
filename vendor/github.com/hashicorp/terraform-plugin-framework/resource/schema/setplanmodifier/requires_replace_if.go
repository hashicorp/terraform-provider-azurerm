// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package setplanmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// RequiresReplaceIf returns a plan modifier that conditionally requires
// resource replacement if:
//
//   - The resource is planned for update.
//   - The plan and state values are not equal.
//   - The given function returns true. Returning false will not unset any
//     prior resource replacement.
//
// Use RequiresReplace if the resource replacement should always occur on value
// changes. Use RequiresReplaceIfConfigured if the resource replacement should
// occur on value changes, but only if there is a configuration value (ignore
// unconfigured drift detection changes).
func RequiresReplaceIf(f RequiresReplaceIfFunc, description, markdownDescription string) planmodifier.Set {
	return requiresReplaceIfModifier{
		ifFunc:              f,
		description:         description,
		markdownDescription: markdownDescription,
	}
}

// requiresReplaceIfModifier is an plan modifier that sets RequiresReplace
// on the attribute if a given function is true.
type requiresReplaceIfModifier struct {
	ifFunc              RequiresReplaceIfFunc
	description         string
	markdownDescription string
}

// Description returns a human-readable description of the plan modifier.
func (m requiresReplaceIfModifier) Description(_ context.Context) string {
	return m.description
}

// MarkdownDescription returns a markdown description of the plan modifier.
func (m requiresReplaceIfModifier) MarkdownDescription(_ context.Context) string {
	return m.markdownDescription
}

// PlanModifySet implements the plan modification logic.
func (m requiresReplaceIfModifier) PlanModifySet(ctx context.Context, req planmodifier.SetRequest, resp *planmodifier.SetResponse) {
	// Do not replace on resource creation.
	if req.State.Raw.IsNull() {
		return
	}

	// Do not replace on resource destroy.
	if req.Plan.Raw.IsNull() {
		return
	}

	// Do not replace if the plan and state values are equal.
	if req.PlanValue.Equal(req.StateValue) {
		return
	}

	ifFuncResp := &RequiresReplaceIfFuncResponse{}

	m.ifFunc(ctx, req, ifFuncResp)

	resp.Diagnostics.Append(ifFuncResp.Diagnostics...)
	resp.RequiresReplace = ifFuncResp.RequiresReplace
}
