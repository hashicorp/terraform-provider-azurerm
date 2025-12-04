// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package convert

import (
	"context"
	"fmt"
	"reflect"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Flatten converts a terraform-plugin-framework object into a go-azure-sdk (native Go) object
// it will write any diagnostics back to the supplied diag.Diagnostics pointer
func Flatten(ctx context.Context, apiObject any, fwObject any, diags *diag.Diagnostics) {
	source, target, d := convert(apiObject, fwObject)
	if d.HasError() {
		diags.Append(d...)
		return
	}

	sourcePath := path.Empty()
	targetPath := path.Empty()

	if source.IsValid() && target.IsValid() {
		if sourceType, targetType := source.Type(), target.Type(); sourceType.Kind() == reflect.Struct && targetType.Kind() == reflect.Struct {
			diags.Append(flattenStruct(ctx, sourcePath, apiObject, targetPath, fwObject)...)
			return
		}
	}

	diags.Append(flatten(ctx, sourcePath, source, targetPath, target)...)
}

// flatten does the heavy lifting via reflection to convert the API values to TF
func flatten(ctx context.Context, sourcePath path.Path, source reflect.Value, targetPath path.Path, target reflect.Value) diag.Diagnostics {
	diags := diag.Diagnostics{}

	targetKindStr := target.Kind().String()

	targetVal, ok := target.Interface().(attr.Value)
	if !ok {
		diags.AddError("FlattenError", fmt.Sprintf("%s does not implement attr.Value", targetKindStr))
		return diags
	}

	targetValType := targetVal.Type(ctx)
	switch k := source.Kind(); k {
	case reflect.Bool:
		diags.Append(flattenBool(ctx, source, targetValType, target, false)...)
		return diags

	case reflect.String:
		diags.Append(flattenString(ctx, source, targetValType, target, false)...)
		return diags

	case reflect.Int64:
		diags.Append(flattenInt64(ctx, source, targetValType, target, false)...)
		return diags

	case reflect.Float64:
		diags.Append(flattenFloat(ctx, source, targetValType, target, false)...)
		return diags

	case reflect.Ptr:
		diags.Append(flattenPtr(ctx, sourcePath, source, targetPath, targetValType, target)...)

	case reflect.Slice:
		diags.Append(flattenSlice(ctx, sourcePath, source, targetPath, targetValType, target)...)
		return diags

	case reflect.Map:
		diags.Append(flattenMap(ctx, sourcePath, source, targetPath, targetValType, target)...)

	case reflect.Struct:
		diags.Append(flattenObject(ctx, sourcePath, source, targetPath, targetValType, target, false)...)
		return diags

	default:
		diags.AddError("FlattenError", fmt.Sprintf("%s is not a supported type for %s", source.Kind(), targetPath.String()))
	}

	return diags
}

func flattenStruct(ctx context.Context, sourcePath path.Path, source any, targetPath path.Path, target any) diag.Diagnostics {
	diags := diag.Diagnostics{}

	sourceVal, targetVal, d := convert(source, target)
	if d.HasError() {
		diags.Append(d...)
		return diags
	}

	for i, sourceType := 0, sourceVal.Type(); i < sourceType.NumField(); i++ {
		field := sourceType.Field(i)
		if field.PkgPath != "" {
			continue
		}
		fieldName := field.Name

		targetFieldName := ""
		for i := 0; i < targetVal.NumField(); i++ {
			if targetVal.Type().Field(i).Tag.Get("convert") == fieldName {
				targetFieldName = targetVal.Type().Field(i).Name
				break
			}
		}

		targetFieldVal := findField(ctx, fieldName, sourceVal, targetVal, targetFieldName)
		if !targetFieldVal.IsValid() {
			continue
		}

		if !targetFieldVal.CanSet() {
			continue
		}

		diags.Append(flatten(ctx, sourcePath.AtName(fieldName), sourceVal.Field(i), targetPath.AtName(fieldName), targetFieldVal)...)
		if diags.HasError() {
			diags.AddError("Flattening", fmt.Sprintf("could not flatten (%s)", fieldName))
			return diags
		}
	}

	return diags
}

func flattenBool(ctx context.Context, source reflect.Value, targetType attr.Type, target reflect.Value, returnNull bool) diag.Diagnostics {
	diags := diag.Diagnostics{}

	switch t := targetType.(type) {
	case basetypes.BoolTypable:
		boolValue := types.BoolNull()
		if !returnNull {
			boolValue = types.BoolValue(source.Bool())
		}

		v, d := t.ValueFromBool(ctx, boolValue)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}

		target.Set(reflect.ValueOf(v))
		return diags

	default:
		diags.AddError("Flatten Error", fmt.Sprintf("flattening bool, but was passed %T", targetType))
	}

	return diags
}

func flattenString(ctx context.Context, source reflect.Value, targetType attr.Type, target reflect.Value, returnNull bool) diag.Diagnostics {
	diags := diag.Diagnostics{}

	switch t := targetType.(type) {
	case basetypes.StringTypable:
		stringValue := types.StringNull()
		if !returnNull {
			stringValue = types.StringValue(source.String())
		}

		v, d := t.ValueFromString(ctx, stringValue)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}

		target.Set(reflect.ValueOf(v))
		return diags

	default:
		diags.AddError("Flatten Error", fmt.Sprintf("flattening string, but was passed %T", targetType))
	}

	return diags
}

func flattenInt64(ctx context.Context, source reflect.Value, targetType attr.Type, target reflect.Value, returnNull bool) diag.Diagnostics {
	diags := diag.Diagnostics{}

	switch t := targetType.(type) {
	case basetypes.Int64Typable:
		int64Value := types.Int64Null()
		if !returnNull {
			int64Value = types.Int64Value(source.Int())
		}

		v, d := t.ValueFromInt64(ctx, int64Value)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}

		target.Set(reflect.ValueOf(v))
		return diags

	default:
		diags.AddError("Flatten Error", fmt.Sprintf("flattening int64, but was passed %T", targetType))
	}

	return diags
}

func flattenFloat(ctx context.Context, source reflect.Value, targetType attr.Type, target reflect.Value, returnNull bool) diag.Diagnostics {
	diags := diag.Diagnostics{}

	switch t := targetType.(type) {
	case basetypes.Float64Typable:
		float64Value := types.Float64Null()
		if !returnNull {
			float64Value = types.Float64Value(source.Float())
		}

		v, d := t.ValueFromFloat64(ctx, float64Value)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}

		target.Set(reflect.ValueOf(v))
		return diags

	default:
		diags.AddError("Flatten Error", fmt.Sprintf("flattening int64, but was passed %T", targetType))
	}

	return diags
}

func flattenPtr(ctx context.Context, sourcePath path.Path, source reflect.Value, targetPath path.Path, targetType attr.Type, target reflect.Value) diag.Diagnostics {
	diags := diag.Diagnostics{}

	switch sourceElem, sourceIsNil := source.Elem(), source.IsNil(); source.Type().Elem().Kind() {
	case reflect.Bool:
		diags.Append(flattenBool(ctx, sourceElem, targetType, target, sourceIsNil)...)
		return diags

	case reflect.String:
		diags.Append(flattenString(ctx, sourceElem, targetType, target, sourceIsNil)...)
		return diags

	case reflect.Int64:
		diags.Append(flattenInt64(ctx, sourceElem, targetType, target, sourceIsNil)...)
		return diags

	case reflect.Float64:
		diags.Append(flattenFloat(ctx, sourceElem, targetType, target, sourceIsNil)...)
		return diags

	case reflect.Struct:
		diags.Append(flattenObject(ctx, sourcePath, source, targetPath, targetType, target, sourceIsNil)...)
	}

	return diags
}

func flattenSlice(ctx context.Context, sourcePath path.Path, source reflect.Value, targetPath path.Path, targetType attr.Type, target reflect.Value) diag.Diagnostics {
	diags := diag.Diagnostics{}

	switch elementType := source.Type().Elem(); elementType.Kind() {
	case reflect.String:
		switch t := targetType.(type) {
		case basetypes.ListTypable:
			{
				if source.IsNil() {
					v, d := t.ValueFromList(ctx, types.ListNull(types.StringType))
					diags.Append(d...)
					if diags.HasError() {
						return diags
					}

					target.Set(reflect.ValueOf(v))
					return diags
				}

				elements := make([]attr.Value, source.Len())
				for i := 0; i < source.Len(); i++ {
					elements[i] = types.StringValue(source.Index(i).String())
				}
				list, d := types.ListValue(types.StringType, elements)
				diags.Append(d...)
				if diags.HasError() {
					return diags
				}

				result, d := t.ValueFromList(ctx, list)
				diags.Append(d...)
				if diags.HasError() {
					return diags
				}
				target.Set(reflect.ValueOf(result))
				return diags
			}

		case basetypes.SetTypable:
			{
				if source.IsNil() {
					v, d := t.ValueFromSet(ctx, types.SetNull(types.StringType))
					diags.Append(d...)
					if diags.HasError() {
						return diags
					}

					target.Set(reflect.ValueOf(v))
					return diags
				}

				elements := make([]attr.Value, source.Len())
				for i := 0; i < source.Len(); i++ {
					elements[i] = types.StringValue(source.Index(i).String())
				}
				set, d := types.SetValue(types.StringType, elements)
				diags.Append(d...)
				if diags.HasError() {
					return diags
				}

				result, d := t.ValueFromSet(ctx, set)
				diags.Append(d...)
				if diags.HasError() {
					return diags
				}
				target.Set(reflect.ValueOf(result))
				return diags
			}
		}
	case reflect.Int64:
		switch t := targetType.(type) {
		case basetypes.ListTypable:
			{
				if source.IsNil() {
					v, d := t.ValueFromList(ctx, types.ListNull(types.Int64Type))
					diags.Append(d...)
					if diags.HasError() {
						return diags
					}

					target.Set(reflect.ValueOf(v))
					return diags
				}

				elements := make([]attr.Value, source.Len())
				for i := 0; i < source.Len(); i++ {
					elements[i] = types.Int64Value(source.Index(i).Int())
				}
				list, d := types.ListValue(types.Int64Type, elements)
				diags.Append(d...)
				if diags.HasError() {
					return diags
				}

				result, d := t.ValueFromList(ctx, list)
				diags.Append(d...)
				if diags.HasError() {
					return diags
				}
				target.Set(reflect.ValueOf(result))
				return diags
			}

		case basetypes.SetTypable:
			{
				if source.IsNil() {
					v, d := t.ValueFromSet(ctx, types.SetNull(types.Int64Type))
					diags.Append(d...)
					if diags.HasError() {
						return diags
					}

					target.Set(reflect.ValueOf(v))
					return diags
				}

				elements := make([]attr.Value, source.Len())
				for i := 0; i < source.Len(); i++ {
					elements[i] = types.Int64Value(source.Index(i).Int())
				}
				set, d := types.SetValue(types.Int64Type, elements)
				diags.Append(d...)
				if diags.HasError() {
					return diags
				}

				result, d := t.ValueFromSet(ctx, set)
				diags.Append(d...)
				if diags.HasError() {
					return diags
				}
				target.Set(reflect.ValueOf(result))
				return diags
			}
		}
	case reflect.Float64:
		switch t := targetType.(type) {
		case basetypes.ListTypable:
			{
				if source.IsNil() {
					v, d := t.ValueFromList(ctx, types.ListNull(types.Float64Type))
					diags.Append(d...)
					if diags.HasError() {
						return diags
					}

					target.Set(reflect.ValueOf(v))
					return diags
				}

				elements := make([]attr.Value, source.Len())
				for i := 0; i < source.Len(); i++ {
					elements[i] = types.Float64Value(source.Index(i).Float())
				}
				list, d := types.ListValue(types.Float64Type, elements)
				diags.Append(d...)
				if diags.HasError() {
					return diags
				}

				result, d := t.ValueFromList(ctx, list)
				diags.Append(d...)
				if diags.HasError() {
					return diags
				}
				target.Set(reflect.ValueOf(result))
				return diags
			}

		case basetypes.SetTypable:
			{
				if source.IsNil() {
					v, d := t.ValueFromSet(ctx, types.SetNull(types.Float64Type))
					diags.Append(d...)
					if diags.HasError() {
						return diags
					}

					target.Set(reflect.ValueOf(v))
					return diags
				}

				elements := make([]attr.Value, source.Len())
				for i := 0; i < source.Len(); i++ {
					elements[i] = types.Float64Value(source.Index(i).Float())
				}
				set, d := types.SetValue(types.Float64Type, elements)
				diags.Append(d...)
				if diags.HasError() {
					return diags
				}

				result, d := t.ValueFromSet(ctx, set)
				diags.Append(d...)
				if diags.HasError() {
					return diags
				}
				target.Set(reflect.ValueOf(result))
				return diags
			}
		}
	case reflect.Bool:
		switch t := targetType.(type) {
		case basetypes.ListTypable:
			{
				if source.IsNil() {
					v, d := t.ValueFromList(ctx, types.ListNull(types.BoolType))
					diags.Append(d...)
					if diags.HasError() {
						return diags
					}

					target.Set(reflect.ValueOf(v))
					return diags
				}

				elements := make([]attr.Value, source.Len())
				for i := 0; i < source.Len(); i++ {
					elements[i] = types.BoolValue(source.Index(i).Bool())
				}
				list, d := types.ListValue(types.BoolType, elements)
				diags.Append(d...)
				if diags.HasError() {
					return diags
				}

				result, d := t.ValueFromList(ctx, list)
				diags.Append(d...)
				if diags.HasError() {
					return diags
				}
				target.Set(reflect.ValueOf(result))
				return diags
			}

		case basetypes.SetTypable:
			{
				if source.IsNil() {
					v, d := t.ValueFromSet(ctx, types.SetNull(types.BoolType))
					diags.Append(d...)
					if diags.HasError() {
						return diags
					}

					target.Set(reflect.ValueOf(v))
					return diags
				}

				elements := make([]attr.Value, source.Len())
				for i := 0; i < source.Len(); i++ {
					elements[i] = types.BoolValue(source.Index(i).Bool())
				}
				set, d := types.SetValue(types.BoolType, elements)
				diags.Append(d...)
				if diags.HasError() {
					return diags
				}

				result, d := t.ValueFromSet(ctx, set)
				diags.Append(d...)
				if diags.HasError() {
					return diags
				}
				target.Set(reflect.ValueOf(result))
				return diags
			}
		}
	case reflect.Ptr:
		switch elementType.Elem().Kind() {
		case reflect.String:
			switch t := targetType.(type) {
			case basetypes.ListTypable:
				{
					if source.IsNil() {
						v, d := t.ValueFromList(ctx, types.ListNull(types.StringType))
						diags.Append(d...)
						if diags.HasError() {
							return diags
						}

						target.Set(reflect.ValueOf(v))
						return diags
					}

					elements := make([]attr.Value, source.Len())
					for i := 0; i < source.Len(); i++ {
						elements[i] = types.StringValue(source.Index(i).String())
					}
					list, d := types.ListValue(types.StringType, elements)
					diags.Append(d...)
					if diags.HasError() {
						return diags
					}

					result, d := t.ValueFromList(ctx, list)
					diags.Append(d...)
					if diags.HasError() {
						return diags
					}
					target.Set(reflect.ValueOf(result))
					return diags
				}

			case basetypes.SetTypable:
				{
					if source.IsNil() {
						v, d := t.ValueFromSet(ctx, types.SetNull(types.StringType))
						diags.Append(d...)
						if diags.HasError() {
							return diags
						}

						target.Set(reflect.ValueOf(v))
						return diags
					}

					elements := make([]attr.Value, source.Len())
					for i := 0; i < source.Len(); i++ {
						elements[i] = types.StringValue(source.Index(i).String())
					}
					set, d := types.SetValue(types.StringType, elements)
					diags.Append(d...)
					if diags.HasError() {
						return diags
					}

					result, d := t.ValueFromSet(ctx, set)
					diags.Append(d...)
					if diags.HasError() {
						return diags
					}
					target.Set(reflect.ValueOf(result))
					return diags
				}
			}
		case reflect.Int64:
			switch t := targetType.(type) {
			case basetypes.ListTypable:
				{
					if source.IsNil() {
						v, d := t.ValueFromList(ctx, types.ListNull(types.Int64Type))
						diags.Append(d...)
						if diags.HasError() {
							return diags
						}

						target.Set(reflect.ValueOf(v))
						return diags
					}

					elements := make([]attr.Value, source.Len())
					for i := 0; i < source.Len(); i++ {
						elements[i] = types.Int64Value(source.Index(i).Int())
					}
					list, d := types.ListValue(types.Int64Type, elements)
					diags.Append(d...)
					if diags.HasError() {
						return diags
					}

					result, d := t.ValueFromList(ctx, list)
					diags.Append(d...)
					if diags.HasError() {
						return diags
					}
					target.Set(reflect.ValueOf(result))
					return diags
				}

			case basetypes.SetTypable:
				{
					if source.IsNil() {
						v, d := t.ValueFromSet(ctx, types.SetNull(types.Int64Type))
						diags.Append(d...)
						if diags.HasError() {
							return diags
						}

						target.Set(reflect.ValueOf(v))
						return diags
					}

					elements := make([]attr.Value, source.Len())
					for i := 0; i < source.Len(); i++ {
						elements[i] = types.Int64Value(source.Index(i).Int())
					}
					set, d := types.SetValue(types.Int64Type, elements)
					diags.Append(d...)
					if diags.HasError() {
						return diags
					}

					result, d := t.ValueFromSet(ctx, set)
					diags.Append(d...)
					if diags.HasError() {
						return diags
					}
					target.Set(reflect.ValueOf(result))
					return diags
				}
			}
		case reflect.Float64:
			switch t := targetType.(type) {
			case basetypes.ListTypable:
				{
					if source.IsNil() {
						v, d := t.ValueFromList(ctx, types.ListNull(types.Float64Type))
						diags.Append(d...)
						if diags.HasError() {
							return diags
						}

						target.Set(reflect.ValueOf(v))
						return diags
					}

					elements := make([]attr.Value, source.Len())
					for i := 0; i < source.Len(); i++ {
						elements[i] = types.Float64Value(source.Index(i).Float())
					}
					list, d := types.ListValue(types.Float64Type, elements)
					diags.Append(d...)
					if diags.HasError() {
						return diags
					}

					result, d := t.ValueFromList(ctx, list)
					diags.Append(d...)
					if diags.HasError() {
						return diags
					}
					target.Set(reflect.ValueOf(result))
					return diags
				}

			case basetypes.SetTypable:
				{
					if source.IsNil() {
						v, d := t.ValueFromSet(ctx, types.SetNull(types.Float64Type))
						diags.Append(d...)
						if diags.HasError() {
							return diags
						}

						target.Set(reflect.ValueOf(v))
						return diags
					}

					elements := make([]attr.Value, source.Len())
					for i := 0; i < source.Len(); i++ {
						elements[i] = types.Float64Value(source.Index(i).Float())
					}
					set, d := types.SetValue(types.Float64Type, elements)
					diags.Append(d...)
					if diags.HasError() {
						return diags
					}

					result, d := t.ValueFromSet(ctx, set)
					diags.Append(d...)
					if diags.HasError() {
						return diags
					}
					target.Set(reflect.ValueOf(result))
					return diags
				}
			}
		case reflect.Bool:
			switch t := targetType.(type) {
			case basetypes.ListTypable:
				{
					if source.IsNil() {
						v, d := t.ValueFromList(ctx, types.ListNull(types.BoolType))
						diags.Append(d...)
						if diags.HasError() {
							return diags
						}

						target.Set(reflect.ValueOf(v))
						return diags
					}

					elements := make([]attr.Value, source.Len())
					for i := 0; i < source.Len(); i++ {
						elements[i] = types.BoolValue(source.Index(i).Bool())
					}
					list, d := types.ListValue(types.BoolType, elements)
					diags.Append(d...)
					if diags.HasError() {
						return diags
					}

					result, d := t.ValueFromList(ctx, list)
					diags.Append(d...)
					if diags.HasError() {
						return diags
					}
					target.Set(reflect.ValueOf(result))
					return diags
				}

			case basetypes.SetTypable:
				{
					if source.IsNil() {
						v, d := t.ValueFromSet(ctx, types.SetNull(types.BoolType))
						diags.Append(d...)
						if diags.HasError() {
							return diags
						}

						target.Set(reflect.ValueOf(v))
						return diags
					}

					elements := make([]attr.Value, source.Len())
					for i := 0; i < source.Len(); i++ {
						elements[i] = types.BoolValue(source.Index(i).Bool())
					}
					set, d := types.SetValue(types.BoolType, elements)
					diags.Append(d...)
					if diags.HasError() {
						return diags
					}

					result, d := t.ValueFromSet(ctx, set)
					diags.Append(d...)
					if diags.HasError() {
						return diags
					}
					target.Set(reflect.ValueOf(result))
					return diags
				}
			}
		}

	case reflect.Struct:
		if t, ok := targetType.(typehelpers.NestedObjectCollectionType); ok {
			diags.Append(flattenSliceOfStructNestedObjectCollection(ctx, sourcePath, source, targetPath, t, target)...)
			return diags
		}
	default:
		diags.AddError("Flatten Slice Error", fmt.Sprintf("unsupported source type: %T", elementType.String()))
	}

	return diags
}

func flattenSliceOfStructNestedObjectCollection(ctx context.Context, sourcePath path.Path, source reflect.Value, targetPath path.Path, targetType typehelpers.NestedObjectCollectionType, target reflect.Value) diag.Diagnostics {
	diags := diag.Diagnostics{}

	if source.IsNil() {
		v, d := targetType.NullValue(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}

		target.Set(reflect.ValueOf(v))
		return diags
	}

	n := source.Len()
	t, d := targetType.NewObjectSlice(ctx, n, n)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	result := reflect.ValueOf(t)
	for i := 0; i < n; i++ {
		tp, d := targetType.NewObjectPtr(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}

		diags.Append(flattenStruct(ctx, sourcePath, source.Index(i).Interface(), targetPath, tp)...)
		if diags.HasError() {
			return diags
		}

		result.Index(i).Set(reflect.ValueOf(tp))
	}

	v, d := targetType.ValueFromObjectSlice(ctx, t)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	target.Set(reflect.ValueOf(v))

	return diags
}

func flattenObject(ctx context.Context, sourcePath path.Path, source reflect.Value, targetPath path.Path, targetType attr.Type, target reflect.Value, returnNull bool) diag.Diagnostics {
	diags := diag.Diagnostics{}

	if t, ok := targetType.(typehelpers.NestedObjectType); ok {
		return flattenStructToNestedObject(ctx, sourcePath, source, targetPath, t, target, returnNull)
	}

	return diags
}

func flattenStructToNestedObject(ctx context.Context, sourcePath path.Path, source reflect.Value, targetPath path.Path, targetType typehelpers.NestedObjectType, target reflect.Value, returnNull bool) diag.Diagnostics {
	diags := diag.Diagnostics{}

	if returnNull {
		v, d := targetType.NullValue(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}

		target.Set(reflect.ValueOf(v))
		return diags
	}

	t, d := targetType.NewObjectPtr(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	diags.Append(flattenStruct(ctx, sourcePath, source.Interface(), targetPath, t)...)
	if diags.HasError() {
		return diags
	}

	v, d := targetType.ValueFromObjectPtr(ctx, t)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	target.Set(reflect.ValueOf(v))

	return diags
}

func flattenMap(ctx context.Context, sourcePath path.Path, source reflect.Value, targetPath path.Path, targetType attr.Type, target reflect.Value) diag.Diagnostics {
	diags := diag.Diagnostics{}

	switch mapKeyKind := source.Type().Key().Kind(); mapKeyKind {
	case reflect.String:
		switch mapElem := source.Type().Elem(); mapElem.Kind() {
		case reflect.Struct:
			{
				switch t := targetType.(type) {
				case basetypes.SetTypable:
					if tSet, ok := t.(typehelpers.NestedObjectCollectionType); ok {
						diags.Append(flattenStructMapToObjectList(ctx, sourcePath, source, targetPath, tSet, target)...)
						return diags
					}

				case basetypes.ListTypable:
					if tList, ok := t.(typehelpers.NestedObjectCollectionType); ok {
						diags.Append(flattenStructMapToObjectList(ctx, sourcePath, source, targetPath, tList, target)...)
						return diags
					}
				default:
					diags.AddError("FlattenMap Error", fmt.Sprintf("unsupported map element struct type: %T", target))
					return diags
				}
			}
		case reflect.String:
			{
				switch t := targetType.(type) {
				case basetypes.MapTypable:
					if source.IsNil() {
						v, d := t.ValueFromMap(ctx, types.MapNull(types.StringType))
						diags.Append(d...)
						if diags.HasError() {
							return diags
						}

						target.Set(reflect.ValueOf(v))
						return diags
					}

					f := source.Interface().(map[string]string)
					elems := make(map[string]attr.Value, len(f))
					for k, v := range f {
						elems[k] = types.StringValue(v)
					}

					m, d := types.MapValue(types.StringType, elems)
					diags.Append(d...)
					if diags.HasError() {
						return diags
					}

					result, d := t.ValueFromMap(ctx, m)
					diags.Append(d...)
					if diags.HasError() {
						return diags
					}

					target.Set(reflect.ValueOf(result))
					return diags
				}
			}

		case reflect.Int64, reflect.Int32, reflect.Int:
			{
				switch t := targetType.(type) {
				case basetypes.MapTypable:
					if source.IsNil() {
						v, d := t.ValueFromMap(ctx, types.MapNull(types.Int64Type))
						diags.Append(d...)
						if diags.HasError() {
							return diags
						}

						target.Set(reflect.ValueOf(v))
						return diags
					}

					f := source.Interface().(map[string]int64)
					elems := make(map[string]attr.Value, len(f))
					for k, v := range f {
						elems[k] = types.Int64Value(v)
					}

					m, d := types.MapValue(types.Int64Type, elems)
					diags.Append(d...)
					if diags.HasError() {
						return diags
					}

					result, d := t.ValueFromMap(ctx, m)
					diags.Append(d...)
					if diags.HasError() {
						return diags
					}

					target.Set(reflect.ValueOf(result))
					return diags
				}
			}

		case reflect.Float64, reflect.Float32:
			{
				switch t := targetType.(type) {
				case basetypes.MapTypable:
					if source.IsNil() {
						v, d := t.ValueFromMap(ctx, types.MapNull(types.Float64Type))
						diags.Append(d...)
						if diags.HasError() {
							return diags
						}

						target.Set(reflect.ValueOf(v))
						return diags
					}

					f := source.Interface().(map[string]float64)
					elems := make(map[string]attr.Value, len(f))
					for k, v := range f {
						elems[k] = types.Float64Value(v)
					}

					m, d := types.MapValue(types.Float64Type, elems)
					diags.Append(d...)
					if diags.HasError() {
						return diags
					}

					result, d := t.ValueFromMap(ctx, m)
					diags.Append(d...)
					if diags.HasError() {
						return diags
					}

					target.Set(reflect.ValueOf(result))
					return diags
				}
			}

		case reflect.Bool:
			{
				switch t := targetType.(type) {
				case basetypes.MapTypable:
					if source.IsNil() {
						v, d := t.ValueFromMap(ctx, types.MapNull(types.BoolType))
						diags.Append(d...)
						if diags.HasError() {
							return diags
						}

						target.Set(reflect.ValueOf(v))
						return diags
					}

					f := source.Interface().(map[string]bool)
					elems := make(map[string]attr.Value, len(f))
					for k, v := range f {
						elems[k] = types.BoolValue(v)
					}

					m, d := types.MapValue(types.BoolType, elems)
					diags.Append(d...)
					if diags.HasError() {
						return diags
					}

					result, d := t.ValueFromMap(ctx, m)
					diags.Append(d...)
					if diags.HasError() {
						return diags
					}

					target.Set(reflect.ValueOf(result))
					return diags
				}
			}

		case reflect.Ptr:
			{
				switch mapElemKind := mapElem.Elem().Kind(); mapElemKind {
				case reflect.Struct:
					if t, ok := targetType.(typehelpers.NestedObjectCollectionType); ok {
						diags.Append(flattenStructMapToObjectList(ctx, sourcePath, source, targetPath, t, target)...)
						return diags

					}
				case reflect.Bool:
					switch t := targetType.(type) {
					case basetypes.MapTypable:
						if source.IsNil() {
							v, d := t.ValueFromMap(ctx, types.MapNull(types.BoolType))
							diags.Append(d...)
							if diags.HasError() {
								return diags
							}

							target.Set(reflect.ValueOf(v))
							return diags
						}

						f := source.Interface().(map[string]*bool)
						elems := make(map[string]attr.Value, len(f))
						for k, v := range f {
							elems[k] = types.BoolPointerValue(v)
						}
						m, d := types.MapValue(types.BoolType, elems)
						diags.Append(d...)
						if diags.HasError() {
							return diags
						}

						result, d := t.ValueFromMap(ctx, m)
						diags.Append(d...)
						if diags.HasError() {
							return diags
						}

						target.Set(reflect.ValueOf(result))

						return diags
					}

				case reflect.Int64, reflect.Int32, reflect.Int:
					switch t := targetType.(type) {
					case basetypes.MapTypable:
						if source.IsNil() {
							v, d := t.ValueFromMap(ctx, types.MapNull(types.Int64Type))
							diags.Append(d...)
							if diags.HasError() {
								return diags
							}

							target.Set(reflect.ValueOf(v))
							return diags
						}

						f := source.Interface().(map[string]*int64)
						elems := make(map[string]attr.Value, len(f))
						for k, v := range f {
							elems[k] = types.Int64PointerValue(v)
						}
						m, d := types.MapValue(types.Int64Type, elems)
						diags.Append(d...)
						if diags.HasError() {
							return diags
						}

						result, d := t.ValueFromMap(ctx, m)
						diags.Append(d...)
						if diags.HasError() {
							return diags
						}

						target.Set(reflect.ValueOf(result))

						return diags
					}

				case reflect.Float64:
					switch t := targetType.(type) {
					case basetypes.MapTypable:
						if source.IsNil() {
							v, d := t.ValueFromMap(ctx, types.MapNull(types.Float64Type))
							diags.Append(d...)
							if diags.HasError() {
								return diags
							}

							target.Set(reflect.ValueOf(v))
							return diags
						}

						f := source.Interface().(map[string]*float64)
						elems := make(map[string]attr.Value, len(f))
						for k, v := range f {
							elems[k] = types.Float64PointerValue(v)
						}
						m, d := types.MapValue(types.Float64Type, elems)
						diags.Append(d...)
						if diags.HasError() {
							return diags
						}

						result, d := t.ValueFromMap(ctx, m)
						diags.Append(d...)
						if diags.HasError() {
							return diags
						}

						target.Set(reflect.ValueOf(result))

						return diags
					}

				case reflect.String:
					switch t := targetType.(type) {
					case basetypes.MapTypable:
						if source.IsNil() {
							v, d := t.ValueFromMap(ctx, types.MapNull(types.StringType))
							diags.Append(d...)
							if diags.HasError() {
								return diags
							}

							target.Set(reflect.ValueOf(v))
							return diags
						}

						f := source.Interface().(map[string]*string)
						elems := make(map[string]attr.Value, len(f))
						for k, v := range f {
							elems[k] = types.StringPointerValue(v)
						}
						m, d := types.MapValue(types.StringType, elems)
						diags.Append(d...)
						if diags.HasError() {
							return diags
						}

						result, d := t.ValueFromMap(ctx, m)
						diags.Append(d...)
						if diags.HasError() {
							return diags
						}

						target.Set(reflect.ValueOf(result))

						return diags
					}

				default:
					// TODO - Do we need to support maps of objects?
					diags.AddError("Map Flatten", fmt.Sprintf("maps of type map[string]%s not currently supported", mapElemKind.String()))
					return diags
				}
			}

		default:
			diags.AddError("FlattenMap Error", fmt.Sprintf("unsupported map element type: %T", target))
			return diags
		}

	default:
		// We only support strings for map keys
		diags.AddError("Flatten Error", fmt.Sprintf("unsupported map  key type: %s", mapKeyKind))
		return diags
	}

	return diags
}

func flattenStructMapToObjectList(ctx context.Context, sourcePath path.Path, source reflect.Value, targetPath path.Path, targetType typehelpers.NestedObjectCollectionType, target reflect.Value) diag.Diagnostics {
	diags := diag.Diagnostics{}

	if source.IsNil() {
		v, d := targetType.NullValue(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}

		target.Set(reflect.ValueOf(v))
		return diags
	}

	n := source.Len()
	t, d := targetType.NewObjectSlice(ctx, n, n)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	tVal := reflect.ValueOf(t)
	index := 0

	for _, k := range source.MapKeys() {
		v, d := targetType.NewObjectPtr(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}

		fromInterface := source.MapIndex(k).Interface()
		if source.MapIndex(k).Kind() == reflect.Ptr {
			fromInterface = source.MapIndex(k).Elem().Interface()
		}

		diags.Append(flattenStruct(ctx, sourcePath, fromInterface, targetPath, v)...)
		if diags.HasError() {
			return diags
		}

		d = setMapKey(v, k)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}

		tVal.Index(index).Set(reflect.ValueOf(v))

		index++
	}

	val, d := targetType.ValueFromObjectSlice(ctx, t)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	target.Set(reflect.ValueOf(val))

	return diags
}

// ProtoV6 construct - future use
func setMapKey(target any, key reflect.Value) diag.Diagnostics {
	diags := diag.Diagnostics{}

	targetVal := reflect.ValueOf(target)
	if targetVal.Kind() == reflect.Ptr {
		targetVal = targetVal.Elem()
	}

	if targetVal.Kind() != reflect.Struct {
		diags.AddError("SetMapKey Error", fmt.Sprintf("expected struct, got %T", target))
		return diags
	}

	for i, targetValType := 0, targetVal.Type(); i < targetVal.NumField(); i++ {
		field := targetValType.Field(i)
		if field.PkgPath != "" {
			continue
		}

		if field.Name != "MapBlockKey" {
			continue
		}

		if _, ok := targetVal.Field(i).Interface().(basetypes.StringValue); ok {
			targetVal.Field(i).Set(reflect.ValueOf(basetypes.NewStringValue(key.String())))
			return diags
		}
	}

	return diags
}
