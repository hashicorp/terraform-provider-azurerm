// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package stringplanmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// RequiresReplace returns a plan modifier that conditionally requires
// resource replacement if:
//
//   - The resource is planned for update.
//   - The plan and state values are not equal.
//
// Use RequiresReplaceIfConfigured if the resource replacement should
// only occur if there is a configuration value (ignore unconfigured drift
// detection changes). Use RequiresReplaceIf if the resource replacement
// should check provider-defined conditional logic.
func RequiresReplace() planmodifier.String {
	return RequiresReplaceIf(
		func(_ context.Context, _ planmodifier.StringRequest, resp *RequiresReplaceIfFuncResponse) {
			resp.RequiresReplace = true
		},
		"If the value of this attribute changes, Terraform will destroy and recreate the resource.",
		"If the value of this attribute changes, Terraform will destroy and recreate the resource.",
	)
}
