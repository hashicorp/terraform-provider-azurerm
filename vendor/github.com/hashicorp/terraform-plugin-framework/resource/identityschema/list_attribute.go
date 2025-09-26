// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package identityschema

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwtype"

	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure the implementation satisfies the desired interfaces.
var (
	_ Attribute                                    = ListAttribute{}
	_ fwschema.AttributeWithValidateImplementation = ListAttribute{}
)

// ListAttribute represents a schema attribute that is a list with a single
// element type. When retrieving the value for this attribute, use types.List
// as the value type unless the CustomType field is set. The ElementType field
// must be set.
//
// In identity schemas, ListAttribute is only permitted to have a primitive ElementType,
// which are:
//   - types.BoolType
//   - types.Float32Type
//   - types.Float64Type
//   - types.Int32Type
//   - types.Int64Type
//   - types.NumberType
//   - types.StringType
//
// Terraform configurations configure this attribute using expressions that
// return a list or directly via square brace syntax.
//
//	# list of strings
//	example_attribute = ["first", "second"]
type ListAttribute struct {
	// ElementType is the type for all elements of the list. This field must be
	// set.
	//
	// ElementType must be a primitive, which are:
	//   - types.BoolType
	//   - types.Float32Type
	//   - types.Float64Type
	//   - types.Int32Type
	//   - types.Int64Type
	//   - types.NumberType
	//   - types.StringType
	ElementType attr.Type

	// CustomType enables the use of a custom attribute type in place of the
	// default basetypes.ListType. When retrieving data, the basetypes.ListValuable
	// associated with this custom type must be used in place of types.List.
	CustomType basetypes.ListTypable

	// RequiredForImport indicates whether the practitioner must enter a value for
	// this attribute when importing a managed resource by this identity.
	// RequiredForImport and OptionalForImport cannot both be true.
	RequiredForImport bool

	// OptionalForImport indicates whether the practitioner can choose to enter a value
	// for this attribute when importing a managed resource by this identity.
	// OptionalForImport and RequiredForImport cannot both be true.
	OptionalForImport bool

	// Description is used in various tooling, like the language server or the documentation
	// generator, to give practitioners more information about what this attribute is,
	// what it's for, and how it should be used. It can be written as plain text with no
	// special formatting, or formatted as Markdown.
	Description string
}

// ApplyTerraform5AttributePathStep returns the result of stepping into a list
// index or an error.
func (a ListAttribute) ApplyTerraform5AttributePathStep(step tftypes.AttributePathStep) (interface{}, error) {
	return a.GetType().ApplyTerraform5AttributePathStep(step)
}

// Equal returns true if the given Attribute is a ListAttribute
// and all fields are equal.
func (a ListAttribute) Equal(o fwschema.Attribute) bool {
	if _, ok := o.(ListAttribute); !ok {
		return false
	}

	return fwschema.AttributesEqual(a, o)
}

// GetDeprecationMessage returns an empty string as identity attributes cannot
// surface deprecation messages.
func (a ListAttribute) GetDeprecationMessage() string {
	return ""
}

// GetDescription returns the Description field value. For identity attributes,
// there is only a single description field that is permitted to contain plaintext or Markdown.
func (a ListAttribute) GetDescription() string {
	return a.Description
}

// GetMarkdownDescription returns the Description field value. For identity attributes,
// there is only a single description field that is permitted to contain Markdown or plaintext.
func (a ListAttribute) GetMarkdownDescription() string {
	return a.Description
}

// GetType returns types.ListType or the CustomType field value if defined.
func (a ListAttribute) GetType() attr.Type {
	if a.CustomType != nil {
		return a.CustomType
	}

	return types.ListType{
		ElemType: a.ElementType,
	}
}

// IsComputed returns false as it's not relevant for identity schemas.
func (a ListAttribute) IsComputed() bool {
	return false
}

// IsOptional returns false as it's not relevant for identity schemas.
func (a ListAttribute) IsOptional() bool {
	return false
}

// IsRequired returns false as it's not relevant for identity schemas.
func (a ListAttribute) IsRequired() bool {
	return false
}

// IsSensitive returns false as it's not relevant for identity schemas.
func (a ListAttribute) IsSensitive() bool {
	return false
}

// IsWriteOnly returns false as it's not relevant for identity schemas.
func (a ListAttribute) IsWriteOnly() bool {
	return false
}

// IsRequiredForImport returns the RequiredForImport field value.
func (a ListAttribute) IsRequiredForImport() bool {
	return a.RequiredForImport
}

// IsOptionalForImport returns the OptionalForImport field value.
func (a ListAttribute) IsOptionalForImport() bool {
	return a.OptionalForImport
}

// ValidateImplementation contains logic for validating the
// provider-defined implementation of the attribute to prevent unexpected
// errors or panics. This logic runs during the GetResourceIdentitySchemas RPC and
// should never include false positives.
func (a ListAttribute) ValidateImplementation(_ context.Context, req fwschema.ValidateImplementationRequest, resp *fwschema.ValidateImplementationResponse) {
	if a.CustomType == nil && a.ElementType == nil {
		resp.Diagnostics.Append(fwschema.AttributeMissingElementTypeDiag(req.Path))
		return
	}

	if a.CustomType == nil && !fwtype.IsAllowedPrimitiveType(a.ElementType) {
		resp.Diagnostics.Append(fwschema.AttributeInvalidElementTypeDiag(req.Path, a.ElementType))
	}
}
