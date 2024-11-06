package tags

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func Expand(input types.Map) (result *map[string]string, diags diag.Diagnostics) {
	if input.IsNull() || input.IsUnknown() {
		return
	}

	diags = input.ElementsAs(context.Background(), &result, false)

	return
}

func Flatten(tags *map[string]string) (result basetypes.MapValue, diags diag.Diagnostics) {
	if tags == nil {
		return basetypes.NewMapNull(basetypes.StringType{}), nil
	}

	return types.MapValueFrom(context.Background(), basetypes.StringType{}, tags)
}
