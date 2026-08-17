package casing

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

type ignoreCaseStringPlanModifier struct{}

var _ planmodifier.String = &ignoreCaseStringPlanModifier{}

func IgnoreCaseStringPlanModifier() planmodifier.String {
	return &ignoreCaseStringPlanModifier{}
}

func (i ignoreCaseStringPlanModifier) Description(_ context.Context) string {
	return "Ignores casing differences in changed values unsing folding."
}

func (i ignoreCaseStringPlanModifier) MarkdownDescription(ctx context.Context) string {
	return i.Description(ctx)
}

func (i ignoreCaseStringPlanModifier) PlanModifyString(ctx context.Context, request planmodifier.StringRequest, response *planmodifier.StringResponse) {
	if !request.PlanValue.IsUnknown() {
		if strings.EqualFold(request.PlanValue.ValueString(), request.StateValue.ValueString()) {
			response.PlanValue = request.StateValue
		}
	}
}
