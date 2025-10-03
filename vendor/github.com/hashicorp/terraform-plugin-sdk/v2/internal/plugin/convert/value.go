// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package convert

import (
	"fmt"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func primitiveTfValue(in cty.Value) (*tftypes.Value, error) {
	primitiveType, err := tftypeFromCtyType(in.Type())
	if err != nil {
		return nil, err
	}

	if in.IsNull() {
		return nullTfValue(primitiveType), nil
	}

	if !in.IsKnown() {
		return unknownTfValue(primitiveType), nil
	}

	var val tftypes.Value
	switch in.Type() {
	case cty.String:
		val = tftypes.NewValue(tftypes.String, in.AsString())
	case cty.Bool:
		val = tftypes.NewValue(tftypes.Bool, in.True())
	case cty.Number:
		val = tftypes.NewValue(tftypes.Number, in.AsBigFloat())
	}

	return &val, nil
}

func listTfValue(in cty.Value) (*tftypes.Value, error) {
	listType, err := tftypeFromCtyType(in.Type())
	if err != nil {
		return nil, err
	}

	if in.IsNull() {
		return nullTfValue(listType), nil
	}

	if !in.IsKnown() {
		return unknownTfValue(listType), nil
	}

	vals := make([]tftypes.Value, 0)

	for _, v := range in.AsValueSlice() {
		tfVal, err := ToTfValue(v)
		if err != nil {
			return nil, err
		}
		vals = append(vals, *tfVal)
	}

	out := tftypes.NewValue(listType, vals)

	return &out, nil
}

func mapTfValue(in cty.Value) (*tftypes.Value, error) {
	mapType, err := tftypeFromCtyType(in.Type())
	if err != nil {
		return nil, err
	}

	if in.IsNull() {
		return nullTfValue(mapType), nil
	}

	if !in.IsKnown() {
		return unknownTfValue(mapType), nil
	}

	vals := make(map[string]tftypes.Value)

	for k, v := range in.AsValueMap() {
		tfVal, err := ToTfValue(v)
		if err != nil {
			return nil, err
		}
		vals[k] = *tfVal
	}

	out := tftypes.NewValue(mapType, vals)

	return &out, nil
}

func setTfValue(in cty.Value) (*tftypes.Value, error) {
	setType, err := tftypeFromCtyType(in.Type())
	if err != nil {
		return nil, err
	}

	if in.IsNull() {
		return nullTfValue(setType), nil
	}

	if !in.IsKnown() {
		return unknownTfValue(setType), nil
	}

	vals := make([]tftypes.Value, 0)

	for _, v := range in.AsValueSlice() {
		tfVal, err := ToTfValue(v)
		if err != nil {
			return nil, err
		}
		vals = append(vals, *tfVal)
	}

	out := tftypes.NewValue(setType, vals)

	return &out, nil
}

func objectTfValue(in cty.Value) (*tftypes.Value, error) {
	objType, err := tftypeFromCtyType(in.Type())
	if err != nil {
		return nil, err
	}

	if in.IsNull() {
		return nullTfValue(objType), nil
	}

	if !in.IsKnown() {
		return unknownTfValue(objType), nil
	}

	vals := make(map[string]tftypes.Value)

	for k, v := range in.AsValueMap() {
		tfVal, err := ToTfValue(v)
		if err != nil {
			return nil, err
		}
		vals[k] = *tfVal
	}

	out := tftypes.NewValue(objType, vals)

	return &out, nil
}

func tupleTfValue(in cty.Value) (*tftypes.Value, error) {
	tupleType, err := tftypeFromCtyType(in.Type())
	if err != nil {
		return nil, err
	}

	if in.IsNull() {
		return nullTfValue(tupleType), nil
	}

	if !in.IsKnown() {
		return unknownTfValue(tupleType), nil
	}

	vals := make([]tftypes.Value, 0)

	for _, v := range in.AsValueSlice() {
		tfVal, err := ToTfValue(v)
		if err != nil {
			return nil, err
		}
		vals = append(vals, *tfVal)
	}

	out := tftypes.NewValue(tupleType, vals)

	return &out, nil
}

func ToTfValue(in cty.Value) (*tftypes.Value, error) {
	ty := in.Type()
	switch {
	case ty.IsPrimitiveType():
		return primitiveTfValue(in)
	case ty.IsListType():
		return listTfValue(in)
	case ty.IsObjectType():
		return objectTfValue(in)
	case ty.IsMapType():
		return mapTfValue(in)
	case ty.IsSetType():
		return setTfValue(in)
	case ty.IsTupleType():
		return tupleTfValue(in)
	default:
		return nil, fmt.Errorf("unsupported type %s", ty)
	}
}

func nullTfValue(ty tftypes.Type) *tftypes.Value {
	nullValue := tftypes.NewValue(ty, nil)
	return &nullValue
}

func unknownTfValue(ty tftypes.Type) *tftypes.Value {
	unknownValue := tftypes.NewValue(ty, tftypes.UnknownValue)
	return &unknownValue
}
