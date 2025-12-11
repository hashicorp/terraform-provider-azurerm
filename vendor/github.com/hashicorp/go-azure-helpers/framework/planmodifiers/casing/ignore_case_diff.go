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

func (i ignoreCaseStringPlanModifier) Description(ctx context.Context) string {
	return "TODO"
}

func (i ignoreCaseStringPlanModifier) MarkdownDescription(ctx context.Context) string {
	return "TODO"
}

func (i ignoreCaseStringPlanModifier) PlanModifyString(ctx context.Context, request planmodifier.StringRequest, response *planmodifier.StringResponse) {
	if strings.EqualFold(request.PlanValue.ValueString(), request.StateValue.ValueString()) {
		response.PlanValue = request.StateValue
	}
}
