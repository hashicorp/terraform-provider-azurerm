// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package types

import "github.com/hashicorp/terraform-plugin-framework/types/basetypes"

type String = basetypes.StringValue

// StringNull creates a String with a null value. Determine whether the value is
// null via the String type IsNull method.
func StringNull() basetypes.StringValue {
	return basetypes.NewStringNull()
}

// StringUnknown creates a String with an unknown value. Determine whether the
// value is unknown via the String type IsUnknown method.
func StringUnknown() basetypes.StringValue {
	return basetypes.NewStringUnknown()
}

// StringValue creates a String with a known value. Access the value via the String
// type ValueString method.
func StringValue(value string) basetypes.StringValue {
	return basetypes.NewStringValue(value)
}

// StringPointerValue creates a String with a null value if nil or a known value.
func StringPointerValue(value *string) basetypes.StringValue {
	return basetypes.NewStringPointerValue(value)
}
