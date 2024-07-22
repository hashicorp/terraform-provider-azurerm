// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package reflect

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/attr/xattr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

// Struct builds a new struct using the data in `object`, as long as `object`
// is a `tftypes.Object`. It will take the struct type from `target`, which
// must be a struct type.
//
// The properties on `target` must be tagged with a "tfsdk" label containing
// the field name to map to that property. Every property must be tagged, and
// every property must be present in the type of `object`, and all the
// attributes in the type of `object` must have a corresponding property.
// Properties that don't map to object attributes must have a `tfsdk:"-"` tag,
// explicitly defining them as not part of the object. This is to catch typos
// and other mistakes early.
//
// Struct is meant to be called from Into, not directly.
func Struct(ctx context.Context, typ attr.Type, object tftypes.Value, target reflect.Value, opts Options, path path.Path) (reflect.Value, diag.Diagnostics) {
	var diags diag.Diagnostics

	// this only works with object values, so make sure that constraint is
	// met
	if target.Kind() != reflect.Struct {
		diags.Append(diag.WithPath(path, DiagIntoIncompatibleType{
			Val:        object,
			TargetType: target.Type(),
			Err:        fmt.Errorf("expected a struct type, got %s", target.Type()),
		}))
		return target, diags
	}
	if !object.Type().Is(tftypes.Object{}) {
		diags.Append(diag.WithPath(path, DiagIntoIncompatibleType{
			Val:        object,
			TargetType: target.Type(),
			Err:        fmt.Errorf("cannot reflect %s into a struct, must be an object", object.Type().String()),
		}))
		return target, diags
	}
	attrsType, ok := typ.(attr.TypeWithAttributeTypes)
	if !ok {
		diags.Append(diag.WithPath(path, DiagIntoIncompatibleType{
			Val:        object,
			TargetType: target.Type(),
			Err:        fmt.Errorf("cannot reflect object using type information provided by %T, %T must be an attr.TypeWithAttributeTypes", typ, typ),
		}))
		return target, diags
	}

	// collect a map of fields that are in the object passed in
	var objectFields map[string]tftypes.Value
	err := object.As(&objectFields)
	if err != nil {
		diags.Append(diag.WithPath(path, DiagIntoIncompatibleType{
			Val:        object,
			TargetType: target.Type(),
			Err:        err,
		}))
		return target, diags
	}

	// collect a map of fields that are defined in the tags of the struct
	// passed in
	targetFields, err := getStructTags(ctx, target, path)
	if err != nil {
		diags.Append(diag.WithPath(path, DiagIntoIncompatibleType{
			Val:        object,
			TargetType: target.Type(),
			Err:        fmt.Errorf("error retrieving field names from struct tags: %w", err),
		}))
		return target, diags
	}

	// we require an exact, 1:1 match of these fields to avoid typos
	// leading to surprises, so let's ensure they have the exact same
	// fields defined
	var objectMissing, targetMissing []string
	for field := range targetFields {
		if _, ok := objectFields[field]; !ok {
			objectMissing = append(objectMissing, field)
		}
	}
	for field := range objectFields {
		if _, ok := targetFields[field]; !ok {
			targetMissing = append(targetMissing, field)
		}
	}
	if len(objectMissing) > 0 || len(targetMissing) > 0 {
		var missing []string
		if len(objectMissing) > 0 {
			missing = append(missing, fmt.Sprintf("Struct defines fields not found in object: %s.", commaSeparatedString(objectMissing)))
		}
		if len(targetMissing) > 0 {
			missing = append(missing, fmt.Sprintf("Object defines fields not found in struct: %s.", commaSeparatedString(targetMissing)))
		}
		diags.Append(diag.WithPath(path, DiagIntoIncompatibleType{
			Val:        object,
			TargetType: target.Type(),
			Err:        fmt.Errorf("mismatch between struct and object: %s", strings.Join(missing, " ")),
		}))
		return target, diags
	}

	attrTypes := attrsType.AttributeTypes()

	// now that we know they match perfectly, fill the struct with the
	// values in the object
	result := reflect.New(target.Type()).Elem()
	for field, structFieldPos := range targetFields {
		attrType, ok := attrTypes[field]
		if !ok {
			diags.Append(diag.WithPath(path, DiagIntoIncompatibleType{
				Val:        object,
				TargetType: target.Type(),
				Err:        fmt.Errorf("could not find type information for attribute in supplied attr.Type %T", typ),
			}))
			return target, diags
		}
		structField := result.Field(structFieldPos)
		fieldVal, fieldValDiags := BuildValue(ctx, attrType, objectFields[field], structField, opts, path.AtName(field))
		diags.Append(fieldValDiags...)

		if diags.HasError() {
			return target, diags
		}
		structField.Set(fieldVal)
	}
	return result, diags
}

// FromStruct builds an attr.Value as produced by `typ` from the data in `val`.
// `val` must be a struct type, and must have all its properties tagged and be
// a 1:1 match with the attributes reported by `typ`. FromStruct will recurse
// into FromValue for each attribute, using the type of the attribute as
// reported by `typ`.
//
// It is meant to be called through FromValue, not directly.
func FromStruct(ctx context.Context, typ attr.TypeWithAttributeTypes, val reflect.Value, path path.Path) (attr.Value, diag.Diagnostics) {
	var diags diag.Diagnostics
	objTypes := map[string]tftypes.Type{}
	objValues := map[string]tftypes.Value{}

	// collect a map of fields that are defined in the tags of the struct
	// passed in
	targetFields, err := getStructTags(ctx, val, path)
	if err != nil {
		err = fmt.Errorf("error retrieving field names from struct tags: %w", err)
		diags.AddAttributeError(
			path,
			"Value Conversion Error",
			"An unexpected error was encountered trying to convert from struct value. This is always an error in the provider. Please report the following to the provider developer:\n\n"+err.Error(),
		)
		return nil, diags
	}

	attrTypes := typ.AttributeTypes()

	var objectMissing, structMissing []string

	for field := range targetFields {
		if _, ok := attrTypes[field]; !ok {
			objectMissing = append(objectMissing, field)
		}
	}

	for attrName, attrType := range attrTypes {
		if attrType == nil {
			objectMissing = append(objectMissing, attrName)
		}

		if _, ok := targetFields[attrName]; !ok {
			structMissing = append(structMissing, attrName)
		}
	}

	if len(objectMissing) > 0 || len(structMissing) > 0 {
		missing := make([]string, 0, len(objectMissing)+len(structMissing))

		if len(objectMissing) > 0 {
			missing = append(missing, fmt.Sprintf("Struct defines fields not found in object: %s.", commaSeparatedString(objectMissing)))
		}

		if len(structMissing) > 0 {
			missing = append(missing, fmt.Sprintf("Object defines fields not found in struct: %s.", commaSeparatedString(structMissing)))
		}

		diags.AddAttributeError(
			path,
			"Value Conversion Error",
			"An unexpected error was encountered trying to convert from struct into an object. "+
				"This is always an error in the provider. Please report the following to the provider developer:\n\n"+
				fmt.Sprintf("Mismatch between struct and object type: %s\n", strings.Join(missing, " "))+
				fmt.Sprintf("Struct: %s\n", val.Type())+
				fmt.Sprintf("Object type: %s", typ),
		)

		return nil, diags
	}

	for name, fieldNo := range targetFields {
		path := path.AtName(name)
		fieldValue := val.Field(fieldNo)

		// If the attr implements xattr.ValidateableAttribute, or xattr.TypeWithValidate,
		// and the attr does not validate then diagnostics will be added here and returned
		// before reaching the switch statement below.
		attrVal, attrValDiags := FromValue(ctx, attrTypes[name], fieldValue.Interface(), path)
		diags.Append(attrValDiags...)

		if diags.HasError() {
			return nil, diags
		}

		tfObjVal, err := attrVal.ToTerraformValue(ctx)
		if err != nil {
			return nil, append(diags, toTerraformValueErrorDiag(err, path))
		}

		switch t := attrVal.(type) {
		case xattr.ValidateableAttribute:
			resp := xattr.ValidateAttributeResponse{}

			t.ValidateAttribute(ctx,
				xattr.ValidateAttributeRequest{
					Path: path,
				},
				&resp,
			)

			diags.Append(resp.Diagnostics...)

			if diags.HasError() {
				return nil, diags
			}
		default:
			//nolint:staticcheck // xattr.TypeWithValidate is deprecated, but we still need to support it.
			if typeWithValidate, ok := attrTypes[name].(xattr.TypeWithValidate); ok {
				diags.Append(typeWithValidate.Validate(ctx, tfObjVal, path)...)

				if diags.HasError() {
					return nil, diags
				}
			}
		}

		tfObjTyp := tfObjVal.Type()

		// If the original attribute type is tftypes.DynamicPseudoType, the value could end up being
		// a concrete type (like tftypes.String, tftypes.List, etc.). In this scenario, the type used
		// to build the final tftypes.Object must stay as tftypes.DynamicPseudoType
		if attrTypes[name].TerraformType(ctx).Is(tftypes.DynamicPseudoType) {
			tfObjTyp = tftypes.DynamicPseudoType
		}

		objValues[name] = tfObjVal
		objTypes[name] = tfObjTyp
	}

	tfVal := tftypes.NewValue(tftypes.Object{
		AttributeTypes: objTypes,
	}, objValues)

	ret, err := typ.ValueFromTerraform(ctx, tfVal)
	if err != nil {
		return nil, append(diags, valueFromTerraformErrorDiag(err, path))
	}

	switch t := ret.(type) {
	case xattr.ValidateableAttribute:
		resp := xattr.ValidateAttributeResponse{}

		t.ValidateAttribute(ctx,
			xattr.ValidateAttributeRequest{
				Path: path,
			},
			&resp,
		)

		diags.Append(resp.Diagnostics...)

		if diags.HasError() {
			return nil, diags
		}
	default:
		//nolint:staticcheck // xattr.TypeWithValidate is deprecated, but we still need to support it.
		if typeWithValidate, ok := typ.(xattr.TypeWithValidate); ok {
			diags.Append(typeWithValidate.Validate(ctx, tfVal, path)...)

			if diags.HasError() {
				return nil, diags
			}
		}
	}

	return ret, diags
}
