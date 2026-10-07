package casing

import (
	"context"

	"github.com/hashicorp/go-azure-helpers/framework/location"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type normaliseLocationStringPlanModifier struct{}

var _ planmodifier.String = &normaliseLocationStringPlanModifier{}

func NormaliseLocationStringPlanModifier() planmodifier.String {
	return &normaliseLocationStringPlanModifier{}
}

func (i normaliseLocationStringPlanModifier) Description(_ context.Context) string {
	return "Normalises Azure locations to their lowercase, zero white-space versions."
}

func (i normaliseLocationStringPlanModifier) MarkdownDescription(ctx context.Context) string {
	return i.Description(ctx)
}

func (i normaliseLocationStringPlanModifier) PlanModifyString(_ context.Context, request planmodifier.StringRequest, response *planmodifier.StringResponse) {
	response.PlanValue = types.StringValue(location.NormalizeNilable(request.PlanValue.ValueStringPointer()))
}
