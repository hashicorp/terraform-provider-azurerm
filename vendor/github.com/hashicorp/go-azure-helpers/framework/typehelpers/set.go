package typehelpers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// DecodeObjectSet returns the decoded list of a List of Objects or an empty struct
func DecodeObjectSet[T any](ctx context.Context, input SetNestedObjectValueOf[T]) (result []T, diags diag.Diagnostics) {
	if !input.IsNull() && len(input.Elements()) > 0 {
		l := make([]T, 0)
		diags.Append(input.ElementsAs(ctx, &l, false)...)
		if diags.HasError() {
			return result, diags
		}

		result = l
	}

	return result, diags
}
