package casing

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type forceLowerCaseStringPlanModifier struct{}

var _ planmodifier.String = &forceLowerCaseStringPlanModifier{}

func ForceLowerCaseStringPlanModifier() planmodifier.String {
	return &forceLowerCaseStringPlanModifier{}
}

func (i forceLowerCaseStringPlanModifier) Description(_ context.Context) string {
	return "TODO"
}

func (i forceLowerCaseStringPlanModifier) MarkdownDescription(ctx context.Context) string {
	return i.Description(ctx)
}

func (i forceLowerCaseStringPlanModifier) PlanModifyString(_ context.Context, request planmodifier.StringRequest, response *planmodifier.StringResponse) {
	response.PlanValue = types.StringValue(strings.ToLower(request.PlanValue.ValueString()))
}
