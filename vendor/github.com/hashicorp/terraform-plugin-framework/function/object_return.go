// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package function

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwfunction"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwtype"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure the implementation satisifies the desired interfaces.
var (
	_ Return                                      = ObjectReturn{}
	_ fwfunction.ReturnWithValidateImplementation = ObjectReturn{}
)

// ObjectReturn represents a function return that is mapping of defined
// attribute names to values. When setting the value for this return, use
// [types.Object] or a compatible Go struct as the value type unless the
// CustomType field is set. The AttributeTypes field must be set.
type ObjectReturn struct {
	// AttributeTypes is the mapping of underlying attribute names to attribute
	// types. This field must be set.
	//
	// Attribute types that contain a collection with a nested dynamic type (i.e. types.List[types.Dynamic]) are not supported.
	// If underlying dynamic collection values are required, replace this return definition with
	// DynamicReturn instead.
	AttributeTypes map[string]attr.Type

	// CustomType enables the use of a custom data type in place of the
	// default [basetypes.ObjectType]. When setting data, the
	// [basetypes.ObjectValuable] implementation associated with this custom
	// type must be used in place of [types.Object].
	CustomType basetypes.ObjectTypable
}

// GetType returns the return data type.
func (r ObjectReturn) GetType() attr.Type {
	if r.CustomType != nil {
		return r.CustomType
	}

	return basetypes.ObjectType{
		AttrTypes: r.AttributeTypes,
	}
}

// NewResultData returns a new result data based on the type.
func (r ObjectReturn) NewResultData(ctx context.Context) (ResultData, *FuncError) {
	value := basetypes.NewObjectUnknown(r.AttributeTypes)

	if r.CustomType == nil {
		return NewResultData(value), nil
	}

	valuable, diags := r.CustomType.ValueFromObject(ctx, value)

	return NewResultData(valuable), FuncErrorFromDiags(ctx, diags)
}

// ValidateImplementation contains logic for validating the
// provider-defined implementation of the Return to prevent unexpected
// errors or panics. This logic runs during the GetProviderSchema RPC and
// should never include false positives.
func (p ObjectReturn) ValidateImplementation(ctx context.Context, req fwfunction.ValidateReturnImplementationRequest, resp *fwfunction.ValidateReturnImplementationResponse) {
	if p.CustomType == nil && fwtype.ContainsCollectionWithDynamic(p.GetType()) {
		resp.Diagnostics.Append(fwtype.ReturnCollectionWithDynamicTypeDiag())
	}
}
