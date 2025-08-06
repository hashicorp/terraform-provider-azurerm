// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package types

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type Tuple = basetypes.TupleValue

// TupleNull creates a Tuple with a null value. Determine whether the value is
// null via the Tuple type IsNull method.
func TupleNull(elementTypes []attr.Type) basetypes.TupleValue {
	return basetypes.NewTupleNull(elementTypes)
}

// TupleUnknown creates a Tuple with an unknown value. Determine whether the
// value is unknown via the Tuple type IsUnknown method.
func TupleUnknown(elementTypes []attr.Type) basetypes.TupleValue {
	return basetypes.NewTupleUnknown(elementTypes)
}

// TupleValue creates a Tuple with a known value. Access the value via the Tuple type Elements method.
func TupleValue(elementTypes []attr.Type, elements []attr.Value) (basetypes.TupleValue, diag.Diagnostics) {
	return basetypes.NewTupleValue(elementTypes, elements)
}

// TupleValueMust creates a Tuple with a known value, converting any diagnostics
// into a panic at runtime. Access the value via the Tuple type Elements method.
//
// This creation function is only recommended to create Tuple values which will
// not potentially affect practitioners, such as testing, or exhaustively
// tested provider logic.
func TupleValueMust(elementTypes []attr.Type, elements []attr.Value) basetypes.TupleValue {
	return basetypes.NewTupleValueMust(elementTypes, elements)
}
