package convert

import (
	"context"
	"fmt"
	"reflect"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Expand converts a terraform-plugin-framework object into a go-azure-sdk (i.e native Go) object.
// It will write any diagnostics back to the supplied diag.Diagnostics pointer
func Expand(ctx context.Context, fwObject any, apiObject any, diags *diag.Diagnostics) {
	source, target, d := convert(fwObject, apiObject)
	if d.HasError() {
		diags.Append(d...)
		return
	}

	sourcePath := path.Empty()
	targetPath := path.Empty()

	if source.IsValid() && target.IsValid() {
		if sourceType, targetType := source.Type(), target.Type(); sourceType.Kind() == reflect.Struct && targetType.Kind() == reflect.Struct {
			diags.Append(expandStruct(ctx, sourcePath, fwObject, targetPath, apiObject)...)
			return
		}
	}

	diags.Append(expand(ctx, sourcePath, source, targetPath, target)...)
}

// expand does the heavy lifting via reflection to convert the tfObject into Go types for use with go-azure-sdk
func expand(ctx context.Context, sourcePath path.Path, source reflect.Value, targetPath path.Path, target reflect.Value) diag.Diagnostics {
	var diags diag.Diagnostics

	sourceVal, ok := source.Interface().(attr.Value)
	if !ok {
		diags.AddError("ExpandError", fmt.Sprintf("%s does not implement attr.Value", source.Kind()))
		return diags
	}

	if sourceVal.IsNull() || sourceVal.IsUnknown() {
		// Nothing to convert
		return diags
	}

	switch t := sourceVal.(type) {
	// Primitives
	case basetypes.BoolValuable:
		{
			diags.Append(expandBool(ctx, t, target)...)
			return diags
		}

	case basetypes.Float64Valuable:
		{
			diags.Append(expandFloat64(ctx, t, target)...)
			return diags
		}

	case basetypes.Int64Valuable:
		{
			diags.Append(expandInt64(ctx, t, target)...)
			return diags
		}

	case basetypes.StringValuable:
		{
			diags.Append(expandString(ctx, t, target)...)
			return diags
		}

	// complex / structs
	case basetypes.ObjectValuable:
		{
			diags.Append(expandObject(ctx, sourcePath, t, targetPath, target)...)
			return diags
		}

	case basetypes.ListValuable:
		{
			diags.Append(expandList(ctx, sourcePath, t, targetPath, target)...)
			return diags
		}
	case basetypes.SetValuable:
		{
			diags.Append(expandSet(ctx, sourcePath, t, targetPath, target)...)
			return diags
		}
	case basetypes.MapValuable:
		{
			diags.Append(expandMap(ctx, t, target)...)
			return diags
		}

	}

	return diags
}

func expandBool(ctx context.Context, source basetypes.BoolValuable, target reflect.Value) diag.Diagnostics {
	var diags diag.Diagnostics

	v, d := source.ToBoolValue(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	switch t := target.Type(); target.Kind() {
	case reflect.Bool:
		{
			target.SetBool(v.ValueBool())
			return diags
		}

	case reflect.Ptr:
		switch tElem := t.Elem(); tElem.Kind() {
		case reflect.Bool:
			{
				target.Set(reflect.ValueOf(v.ValueBoolPointer()))
				return diags
			}
		}
	}

	return diags
}

func expandFloat64(ctx context.Context, source basetypes.Float64Valuable, target reflect.Value) diag.Diagnostics {
	var diags diag.Diagnostics

	v, d := source.ToFloat64Value(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	switch t := target.Type(); target.Kind() {
	case reflect.Float64:
		{
			target.SetFloat(v.ValueFloat64())
			return diags
		}

	case reflect.Ptr:
		switch tElem := t.Elem(); tElem.Kind() {
		case reflect.Float64:
			{
				target.Set(reflect.ValueOf(v.ValueFloat64Pointer()))
				return diags
			}
		}
	}

	return diags
}

func expandInt64(ctx context.Context, source basetypes.Int64Valuable, target reflect.Value) diag.Diagnostics {
	var diags diag.Diagnostics

	v, d := source.ToInt64Value(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	switch t := target.Type(); target.Kind() {
	case reflect.Int64:
		{
			target.SetInt(v.ValueInt64())
			return diags
		}

	case reflect.Ptr:
		switch tElem := t.Elem(); tElem.Kind() {
		case reflect.Int64:
			{
				target.Set(reflect.ValueOf(v.ValueInt64Pointer()))
				return diags
			}
		}
	}

	return diags
}

func expandString(ctx context.Context, source basetypes.StringValuable, target reflect.Value) diag.Diagnostics {
	var diags diag.Diagnostics

	v, d := source.ToStringValue(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	switch t := target.Type(); target.Kind() {
	case reflect.String:
		{
			target.SetString(v.ValueString())
			return diags
		}

	case reflect.Ptr:
		switch tElem := t.Elem(); tElem.Kind() {
		case reflect.String:
			target.Set(reflect.ValueOf(v.ValueStringPointer()))
			return diags
		}
	}

	// TODO do something for unexpected things, or is it fine to just silently continue?

	return diags
}

func expandObject(ctx context.Context, sourcePath path.Path, source basetypes.ObjectValuable, targetPath path.Path, target reflect.Value) diag.Diagnostics {
	diags := diag.Diagnostics{}

	_, d := source.ToObjectValue(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	switch t := target.Type(); target.Kind() {
	case reflect.Struct:
		{
			if nestedObjectSource, ok := source.(typehelpers.NestedObjectValue); ok {
				diags.Append(expandNestedObjectToStruct(ctx, sourcePath, nestedObjectSource, targetPath, t, target)...)
				return diags
			}
		}

	case reflect.Ptr:
		switch elem := t.Elem(); elem.Kind() {
		case reflect.Struct:
			if nestedObjectSource, ok := source.(typehelpers.NestedObjectValue); ok {
				diags.Append(expandNestedObjectToStruct(ctx, sourcePath, nestedObjectSource, targetPath, elem, target)...)
				return diags
			}
		}
	}

	// TODO do something for unexpected usage, or is it fine to just silently continue?

	return diags
}

func expandList(ctx context.Context, sourcePath path.Path, source basetypes.ListValuable, targetPath path.Path, target reflect.Value) diag.Diagnostics {
	diags := diag.Diagnostics{}

	v, d := source.ToListValue(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	switch v.ElementType(ctx).(type) {
	case basetypes.StringTypable, basetypes.Int64Typable, basetypes.Float64Typable, basetypes.BoolTypable:
		{
			diags.Append(expandListOfPrimitive(ctx, v, target)...)
			return diags
		}
	case basetypes.ObjectTypable:
		{
			if s, ok := source.(typehelpers.NestedObjectCollectionValue); ok {
				diags.Append(expandNestedObjectCollection(ctx, sourcePath, s, targetPath, target)...)
			}
		}
	case basetypes.MapTypable:
		{
			// TODO?
			diags.AddError("unsupported list type", "lists of maps are not currently supported")
		}
	}

	return diags
}

func expandListOfPrimitive(ctx context.Context, source basetypes.ListValue, target reflect.Value) diag.Diagnostics {
	diags := diag.Diagnostics{}

	switch target.Kind() {
	case reflect.Slice:
		switch elem := target.Type().Elem(); elem.Kind() {
		case reflect.String:
			{
				s := make([]string, 0)
				diags.Append(source.ElementsAs(ctx, &s, false)...)
				target.Set(reflect.ValueOf(s))
			}

		case reflect.Bool:
			{
				b := make([]bool, 0)
				diags.Append(source.ElementsAs(ctx, &b, false)...)
				target.Set(reflect.ValueOf(b))
			}

		case reflect.Float64:
			{
				f := make([]float64, 0)
				diags.Append(source.ElementsAs(ctx, &f, false)...)
				target.Set(reflect.ValueOf(f))
			}

		case reflect.Int64:
			{
				i := make([]int64, 0)
				diags.Append(source.ElementsAs(ctx, &i, false)...)
				target.Set(reflect.ValueOf(i))
			}
		}
	}

	return diags
}

func expandListOfObject(ctx context.Context, sourcePath path.Path, source typehelpers.NestedObjectCollectionValue, targetPath path.Path, target reflect.Value) diag.Diagnostics {
	diags := diag.Diagnostics{}

	switch t := target.Type(); target.Kind() {
	case reflect.Struct:
		{
			diags.Append(expandNestedObjectToStruct(ctx, sourcePath, source, targetPath, t, target)...)
			return diags
		}

	case reflect.Ptr:
		{
			switch elem := t.Elem(); elem.Kind() {
			case reflect.Struct:
				{
					diags.Append(expandNestedObjectToStruct(ctx, sourcePath, source, targetPath, elem, target)...)
					return diags
				}
			}
		}

	case reflect.Map:
		{
			switch elem := t.Elem(); elem.Kind() {
			case reflect.Struct, reflect.Ptr:
				// Maps and pointers to maps can be treated the same here
				{
					diags.Append(expandNestedObjectToMap(ctx, sourcePath, source, targetPath, elem, target)...)
				}
			}
		}

	case reflect.Slice:
		{
			switch elem := t.Elem(); elem.Kind() {
			case reflect.Struct:
				{
					diags.Append(expandNestedObjectToSlice(ctx, sourcePath, source, targetPath, t, elem, target)...)
				}
			}
		}
	}

	return diags
}

func expandSet(ctx context.Context, sourcePath path.Path, source basetypes.SetValuable, targetPath path.Path, target reflect.Value) diag.Diagnostics {
	diags := diag.Diagnostics{}

	v, d := source.ToSetValue(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	switch v.ElementType(ctx).(type) {
	case basetypes.StringTypable, basetypes.Int64Typable, basetypes.Float64Typable, basetypes.BoolTypable:
		{
			diags.Append(expandSetOfPrimitive(ctx, v, target)...)
			return diags
		}
	case basetypes.ObjectTypable:
		{
			if s, ok := source.(typehelpers.NestedObjectCollectionValue); ok {
				diags.Append(expandListOfObject(ctx, sourcePath, s, targetPath, target)...)
			}
		}
	case basetypes.MapTypable:
		{
			// TODO?
			diags.AddError("unsupported set type", "sets of maps are not currently supported")
		}
	}

	return diags
}

func expandSetOfPrimitive(ctx context.Context, source basetypes.SetValue, target reflect.Value) diag.Diagnostics {
	diags := diag.Diagnostics{}

	switch target.Kind() {
	case reflect.Slice:
		switch elem := target.Type().Elem(); elem.Kind() {
		case reflect.String:
			{
				s := make([]string, 0)
				diags.Append(source.ElementsAs(ctx, &s, false)...)
				target.Set(reflect.ValueOf(s))
			}

		case reflect.Bool:
			{
				b := make([]bool, 0)
				diags.Append(source.ElementsAs(ctx, &b, false)...)
				target.Set(reflect.ValueOf(b))
			}

		case reflect.Float64:
			{
				f := make([]float64, 0)
				diags.Append(source.ElementsAs(ctx, &f, false)...)
				target.Set(reflect.ValueOf(f))
			}

		case reflect.Int64:
			{
				i := make([]int64, 0)
				diags.Append(source.ElementsAs(ctx, &i, false)...)
				target.Set(reflect.ValueOf(i))
			}
		}
	}

	return diags
}

func expandNestedObjectCollection(ctx context.Context, sourcePath path.Path, source typehelpers.NestedObjectCollectionValue, targetPath path.Path, target reflect.Value) diag.Diagnostics {
	diags := diag.Diagnostics{}

	switch targetType := target.Type(); target.Kind() {
	case reflect.Struct:
		pathAtListIndex := sourcePath.AtListIndex(0)
		diags.Append(expandNestedObjectToStruct(ctx, pathAtListIndex, source, targetPath, targetType, target)...)
		return diags

	case reflect.Pointer:
		switch tElem := targetType.Elem(); tElem.Kind() {
		case reflect.Struct:
			pathAtListIndex := sourcePath.AtListIndex(0)
			diags.Append(expandNestedObjectToStruct(ctx, pathAtListIndex, source, targetPath, tElem, target)...)
			return diags
		}

	case reflect.Interface:
		pathAtListIndex := sourcePath.AtListIndex(0)
		diags.Append(expandNestedObjectToStruct(ctx, pathAtListIndex, source, targetPath, targetType, target)...)
		return diags

	case reflect.Map:
		switch tElem := targetType.Elem(); tElem.Kind() {
		case reflect.Struct:
			diags.Append(expandNestedObjectToMap(ctx, sourcePath, source, targetPath, tElem, target)...)
			return diags

		case reflect.Pointer:
			diags.Append(expandNestedObjectToMap(ctx, sourcePath, source, targetPath, tElem, target)...)
			return diags
		}

	case reflect.Slice:
		switch tElem := targetType.Elem(); tElem.Kind() {
		case reflect.Struct:
			diags.Append(expandNestedObjectToSlice(ctx, sourcePath, source, targetPath, targetType, tElem, target)...)
			return diags

		case reflect.Pointer:
			switch elem := tElem.Elem(); elem.Kind() {
			case reflect.Struct:
				diags.Append(expandNestedObjectToSlice(ctx, sourcePath, source, targetPath, targetType, elem, target)...)
				return diags
			}

		case reflect.Interface:
			diags.Append(expandNestedObjectToSlice(ctx, sourcePath, source, targetPath, targetType, tElem, target)...)
			return diags
		}
	}

	diags.AddError("Incompatible types", fmt.Sprintf("nestedObjectCollection[%s] cannot be expanded to %s", source.Type(ctx).(attr.TypeWithElementType).ElementType(), target.Kind()))
	return diags
}

func expandNestedObjectToStruct(ctx context.Context, sourcePath path.Path, source typehelpers.NestedObjectValue, targetPath path.Path, targetType reflect.Type, targetValue reflect.Value) diag.Diagnostics {
	var diags diag.Diagnostics

	v, d := source.ToObjectPtr(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	into := reflect.New(targetType)
	diags.Append(expandStruct(ctx, sourcePath, v, targetPath, into.Interface())...)
	if diags.HasError() {
		return diags
	}

	if targetValue.Type().Kind() == reflect.Struct {
		targetValue.Set(into.Elem())
	} else {
		targetValue.Set(into)
	}

	return diags
}

func expandMap(ctx context.Context, source basetypes.MapValuable, target reflect.Value) diag.Diagnostics {
	diags := diag.Diagnostics{}

	v, d := source.ToMapValue(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	switch v.ElementType(ctx).(type) {
	case basetypes.StringTypable, basetypes.Int64Typable, basetypes.Float64Typable, basetypes.BoolTypable:
		{
			diags.Append(expandMapOfPrimitive(ctx, v, target)...)
			return diags
		}
	}

	return diags
}

func expandMapOfPrimitive(ctx context.Context, source basetypes.MapValue, target reflect.Value) diag.Diagnostics {
	diags := diag.Diagnostics{}

	switch target.Kind() {
	case reflect.Map:
		switch tMapKey := target.Type().Key(); tMapKey.Kind() {
		case reflect.String: // key
			switch tMapElem := target.Type().Elem(); tMapElem.Kind() {
			case reflect.String:
				{
					var to map[string]string
					diags.Append(source.ElementsAs(ctx, &to, false)...)
					if diags.HasError() {
						return diags
					}

					target.Set(reflect.ValueOf(to))
					return diags
				}

			case reflect.Int64:
				{
					var to map[string]int64
					diags.Append(source.ElementsAs(ctx, &to, false)...)
					if diags.HasError() {
						return diags
					}

					target.Set(reflect.ValueOf(to))
					return diags
				}

			case reflect.Float64:
				{
					var to map[string]float64
					diags.Append(source.ElementsAs(ctx, &to, false)...)
					if diags.HasError() {
						return diags
					}

					target.Set(reflect.ValueOf(to))
					return diags
				}

			case reflect.Bool:
				{
					var to map[string]bool
					diags.Append(source.ElementsAs(ctx, &to, false)...)
					if diags.HasError() {
						return diags
					}

					target.Set(reflect.ValueOf(to))
					return diags
				}

			case reflect.Ptr:
				{
					switch k := tMapElem.Elem().Kind(); k {
					case reflect.String:
						var to map[string]*string
						diags.Append(source.ElementsAs(ctx, &to, false)...)
						if diags.HasError() {
							return diags
						}

						target.Set(reflect.ValueOf(to))
						return diags

					case reflect.Int64:
						{
							var to map[string]*int64
							diags.Append(source.ElementsAs(ctx, &to, false)...)
							if diags.HasError() {
								return diags
							}

							target.Set(reflect.ValueOf(to))
							return diags
						}

					case reflect.Float64:
						{
							var to map[string]*float64
							diags.Append(source.ElementsAs(ctx, &to, false)...)
							if diags.HasError() {
								return diags
							}

							target.Set(reflect.ValueOf(to))
							return diags
						}

					case reflect.Bool:
						{
							var to map[string]*bool
							diags.Append(source.ElementsAs(ctx, &to, false)...)
							if diags.HasError() {
								return diags
							}

							target.Set(reflect.ValueOf(to))
							return diags
						}
					}
				}
			}
		}
	}

	return diags
}

// protocol v6 only
func expandNestedObjectToMap(ctx context.Context, sourcePath path.Path, source typehelpers.NestedObjectCollectionValue, targetPath path.Path, targetType reflect.Type, targetValue reflect.Value) diag.Diagnostics {
	diags := diag.Diagnostics{}

	from, d := source.ToObjectSlice(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	if targetType.Kind() == reflect.Ptr {
		targetType = targetType.Elem()
	}

	// Create a new target slice and expand each element.
	f := reflect.ValueOf(from)
	m := reflect.MakeMap(targetValue.Type())
	for i := 0; i < f.Len(); i++ {
		// Create a new target structure and walk its fields.
		target := reflect.New(targetType)
		diags.Append(expandStruct(ctx, sourcePath, f.Index(i).Interface(), targetPath, target.Interface())...)
		if diags.HasError() {
			return diags
		}

		key, d := extractMapKeyValue(f.Index(i).Interface())
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}

		// Set value (or pointer) in the target map.
		if target.Type().Elem().Kind() == reflect.Struct {
			m.SetMapIndex(key, target.Elem())
		} else {
			m.SetMapIndex(key, target)
		}
	}

	targetValue.Set(m)
	return diags
}

func expandNestedObjectToSlice(ctx context.Context, sourcePath path.Path, source typehelpers.NestedObjectCollectionValue, targetPath path.Path, targetType reflect.Type, targetElemType reflect.Type, targetValue reflect.Value) diag.Diagnostics {
	diags := diag.Diagnostics{}
	from, d := source.ToObjectSlice(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	f := reflect.ValueOf(from)
	n := f.Len()
	t := reflect.MakeSlice(targetType, n, n)
	for i := 0; i < n; i++ {
		target := reflect.New(targetElemType)
		diags.Append(expandStruct(ctx, sourcePath, f.Index(i).Interface(), targetPath, target.Interface())...)
		if diags.HasError() {
			return diags
		}

		if target.Type().Elem().Kind() == reflect.Struct {
			t.Index(i).Set(target.Elem())
		} else {
			t.Index(i).Set(target)
		}
	}

	targetValue.Set(t)
	return diags
}

func expandStruct(ctx context.Context, sourcePath path.Path, source any, targetPath path.Path, target any) diag.Diagnostics {
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
		if fieldName == "MapBlockKey" {
			continue
		}

		targetFieldVal := findField(ctx, fieldName, sourceVal, targetVal, field.Tag.Get("convert"))
		if !targetFieldVal.IsValid() {
			continue
		}

		if !targetFieldVal.CanSet() {
			continue
		}

		diags.Append(expand(ctx, sourcePath, sourceVal.Field(i), targetPath, targetFieldVal)...)
		if diags.HasError() {
			diags.AddError("Expanding", fmt.Sprintf("could not expand (%s)", fieldName))
			return diags
		}
	}

	return diags
}

// func expandConvertStruct(ctx context.Context, sourcePath path.Path, source any, targetPath path.Path, target any) diag.Diagnostics {
//
// }

// findField looks for the matching API struct name in the target struct
func findField(ctx context.Context, fieldName string, _ reflect.Value, target reflect.Value, tagHint string) reflect.Value {
	// specific apiName struct tag take precedence
	if v := target.FieldByName(tagHint); v.IsValid() {
		return v
	}

	if v := target.FieldByName(fieldName); v.IsValid() {
		return v
	}

	// TODO - Can we do case-insensitive matching? (And should we?)

	// TODO - resource manager suffix trimming find? e.g. ThingProperties == Thing

	return target.FieldByName(fieldName)
}

// extractMapKeyValue gets the reflect.Value of the Key of a map entry
func extractMapKeyValue(source any) (reflect.Value, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	valFrom := reflect.ValueOf(source)
	if kind := valFrom.Kind(); kind == reflect.Ptr {
		valFrom = valFrom.Elem()
	}

	for i, typFrom := 0, valFrom.Type(); i < typFrom.NumField(); i++ {
		field := typFrom.Field(i)
		if field.PkgPath != "" {
			continue // unexported fields shouldn't be considered.
		}

		if field.Name == "MapBlockKey" {
			fieldVal := valFrom.Field(i)

			if v, ok := fieldVal.Interface().(basetypes.StringValue); ok {
				return reflect.ValueOf(v.ValueString()), diags
			}

			fieldType := fieldVal.Type()
			method, found := fieldType.MethodByName("ValueString")
			if found {
				result := fieldType.Method(method.Index).Func.Call([]reflect.Value{fieldVal})
				if len(result) > 0 {
					return result[0], diags
				}
			}

			return valFrom.Field(i), diags
		}
	}

	diags.AddError("convert", "unable to find map block key")

	return reflect.Zero(reflect.TypeOf("")), diags
}
