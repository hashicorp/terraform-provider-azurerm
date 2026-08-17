package validators

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatorfuncerr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type resourceId struct {
	id resourceids.ResourceId
}

var _ validator.String = &resourceId{}

var _ function.StringParameterValidator = &resourceId{}

func AzureResourceManagerId(id resourceids.ResourceId) resourceId {
	return resourceId{
		id: id,
	}
}

func (r resourceId) Description(ctx context.Context) string {
	return "validates that the provided string is an Azure resource ID"
}

func (r resourceId) MarkdownDescription(ctx context.Context) string {
	return r.Description(ctx)
}

func (r resourceId) ValidateString(_ context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()

	parser := resourceids.NewParserFromResourceIdType(r.id)
	parsed, err := parser.Parse(value, false)
	if err != nil {
		response.Diagnostics.AddError("ID validation error", err.Error())
		return
	}

	for i, segment := range r.id.Segments() {
		if _, ok := parsed.Parsed[segment.Name]; !ok {
			response.Diagnostics.AddError("id segment error", fmt.Sprintf("expected the segment %d (type %q / name %q) to have a value but it didn't", i, segment.Type, segment.Name))
			return
		}
	}
}

func (r resourceId) ValidateParameterString(_ context.Context, request function.StringParameterValidatorRequest, response *function.StringParameterValidatorResponse) {
	if request.Value.IsNull() || request.Value.IsUnknown() {
		return
	}

	value := request.Value.ValueString()

	parser := resourceids.NewParserFromResourceIdType(r.id)
	parsed, err := parser.Parse(value, false)
	if err != nil {
		response.Error = validatorfuncerr.InvalidParameterValueFuncError(request.ArgumentPosition, "ID validation error", err.Error())
		return
	}

	for i, segment := range r.id.Segments() {
		if _, ok := parsed.Parsed[segment.Name]; !ok {
			response.Error = validatorfuncerr.InvalidParameterValueMatchFuncError(request.ArgumentPosition, "id segment error", fmt.Sprintf("expected the segment %d (type %q / name %q) to have a value but it didn't", i, segment.Type, segment.Name))
			return
		}
	}
}
