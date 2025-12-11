package values

import "github.com/hashicorp/terraform-plugin-framework/types"

func ValueStringPointer(input types.String) *string {
	if input.IsNull() || input.IsUnknown() {
		return nil
	}

	return input.ValueStringPointer()
}

func ValueInt64Pointer(input types.Int64) *int64 {
	if input.IsNull() || input.IsUnknown() {
		return nil
	}

	return input.ValueInt64Pointer()
}

func ValueFloat64Pointer(input types.Float64) *float64 {
	if input.IsNull() || input.IsUnknown() {
		return nil
	}

	return input.ValueFloat64Pointer()
}

func ValueBoolPointer(input types.Bool) *bool {
	if input.IsNull() || input.IsUnknown() {
		return nil
	}

	return input.ValueBoolPointer()
}
