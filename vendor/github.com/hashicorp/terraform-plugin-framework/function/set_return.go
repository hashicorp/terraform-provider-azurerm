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
	_ Return                                      = SetReturn{}
	_ fwfunction.ReturnWithValidateImplementation = SetReturn{}
)

// SetReturn represents a function return that is an unordered collection of a
// single element type. Either the ElementType or CustomType field must be set.
//
// When setting the value for this return:
//
//   - If CustomType is set, use its associated value type.
//   - Otherwise, use [types.Set] or a Go slice value type compatible with the
//     element type.
type SetReturn struct {
	// ElementType is the type for all elements of the set. This field must be
	// set.
	//
	// Element types that contain a dynamic type (i.e. types.Dynamic) are not supported.
	// If underlying dynamic values are required, replace this return definition with
	// DynamicReturn instead.
	ElementType attr.Type

	// CustomType enables the use of a custom data type in place of the
	// default [basetypes.SetType]. When setting data, the
	// [basetypes.SetValuable] implementation associated with this custom
	// type must be used in place of [types.Set].
	CustomType basetypes.SetTypable
}

// GetType returns the return data type.
func (r SetReturn) GetType() attr.Type {
	if r.CustomType != nil {
		return r.CustomType
	}

	return basetypes.SetType{
		ElemType: r.ElementType,
	}
}

// NewResultData returns a new result data based on the type.
func (r SetReturn) NewResultData(ctx context.Context) (ResultData, *FuncError) {
	value := basetypes.NewSetUnknown(r.ElementType)

	if r.CustomType == nil {
		return NewResultData(value), nil
	}

	valuable, diags := r.CustomType.ValueFromSet(ctx, value)

	return NewResultData(valuable), FuncErrorFromDiags(ctx, diags)
}

// ValidateImplementation contains logic for validating the
// provider-defined implementation of the Return to prevent unexpected
// errors or panics. This logic runs during the GetProviderSchema RPC and
// should never include false positives.
func (p SetReturn) ValidateImplementation(ctx context.Context, req fwfunction.ValidateReturnImplementationRequest, resp *fwfunction.ValidateReturnImplementationResponse) {
	if p.CustomType == nil && fwtype.ContainsCollectionWithDynamic(p.GetType()) {
		resp.Diagnostics.Append(fwtype.ReturnCollectionWithDynamicTypeDiag())
	}
}
