// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

// nolint: deadcode unused
func adminPasswordDiffSuppressFunc(_, old, new string, _ *pluginsdk.ResourceData) bool {
	// this is not the greatest hack in the world, this is just a tribute.
	if old == "ignored-as-imported" || new == "ignored-as-imported" {
		return true
	}

	return false
}

type ignoreAdminPasswordDiffSuppressModifier struct{}

func (i ignoreAdminPasswordDiffSuppressModifier) Description(ctx context.Context) string {
	return "Suppresses the diff for Virtual Machines 'admin_password' values for imported resources."
}

func (i ignoreAdminPasswordDiffSuppressModifier) MarkdownDescription(ctx context.Context) string {
	return i.Description(ctx)
}

func (i ignoreAdminPasswordDiffSuppressModifier) PlanModifyString(ctx context.Context, request planmodifier.StringRequest, response *planmodifier.StringResponse) {
	if request.PlanValue.ValueString() == "ignored-as-imported" || request.StateValue.ValueString() == "ignored-as-imported" {
		response.PlanValue = request.StateValue // Set the Plan to the same as the state to prevent showing a diff for this value
	}
}

var _ planmodifier.String = &ignoreAdminPasswordDiffSuppressModifier{}
