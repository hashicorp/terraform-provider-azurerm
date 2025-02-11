// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema/fwxschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure the implementation satisfies the desired interfaces.
var (
	_ Attribute                                    = DynamicAttribute{}
	_ fwschema.AttributeWithValidateImplementation = DynamicAttribute{}
	_ fwschema.AttributeWithDynamicDefaultValue    = DynamicAttribute{}
	_ fwxschema.AttributeWithDynamicPlanModifiers  = DynamicAttribute{}
	_ fwxschema.AttributeWithDynamicValidators     = DynamicAttribute{}
)

// DynamicAttribute represents a schema attribute that is a dynamic, rather
// than a single static type. Static types are always preferable over dynamic
// types in Terraform as practitioners will receive less helpful configuration
// assistance from validation error diagnostics and editor integrations. When
// retrieving the value for this attribute, use types.Dynamic as the value type
// unless the CustomType field is set.
//
// The concrete value type for a dynamic is determined at runtime in this order:
//  1. By Terraform, if defined in the configuration (if Required or Optional).
//  2. By the provider (if Computed).
//
// Once the concrete value type has been determined, it must remain consistent between
// plan and apply or Terraform will return an error.
type DynamicAttribute struct {
	// CustomType enables the use of a custom attribute type in place of the
	// default basetypes.DynamicType. When retrieving data, the basetypes.DynamicValuable
	// associated with this custom type must be used in place of types.Dynamic.
	CustomType basetypes.DynamicTypable

	// Required indicates whether the practitioner must enter a value for
	// this attribute or not. Required and Optional cannot both be true,
	// and Required and Computed cannot both be true.
	Required bool

	// Optional indicates whether the practitioner can choose to enter a value
	// for this attribute or not. Optional and Required cannot both be true.
	Optional bool

	// Computed indicates whether the provider may return its own value for
	// this Attribute or not. Required and Computed cannot both be true. If
	// Required and Optional are both false, Computed must be true, and the
	// attribute will be considered "read only" for the practitioner, with
	// only the provider able to set its value.
	Computed bool

	// Sensitive indicates whether the value of this attribute should be
	// considered sensitive data. Setting it to true will obscure the value
	// in CLI output. Sensitive does not impact how values are stored, and
	// practitioners are encouraged to store their state as if the entire
	// file is sensitive.
	Sensitive bool

	// Description is used in various tooling, like the language server, to
	// give practitioners more information about what this attribute is,
	// what it's for, and how it should be used. It should be written as
	// plain text, with no special formatting.
	Description string

	// MarkdownDescription is used in various tooling, like the
	// documentation generator, to give practitioners more information
	// about what this attribute is, what it's for, and how it should be
	// used. It should be formatted using Markdown.
	MarkdownDescription string

	// DeprecationMessage defines warning diagnostic details to display when
	// practitioner configurations use this Attribute. The warning diagnostic
	// summary is automatically set to "Attribute Deprecated" along with
	// configuration source file and line information.
	//
	// Set this field to a practitioner actionable message such as:
	//
	//  - "Configure other_attribute instead. This attribute will be removed
	//    in the next major version of the provider."
	//  - "Remove this attribute's configuration as it no longer is used and
	//    the attribute will be removed in the next major version of the
	//    provider."
	//
	// In Terraform 1.2.7 and later, this warning diagnostic is displayed any
	// time a practitioner attempts to configure a value for this attribute and
	// certain scenarios where this attribute is referenced.
	//
	// In Terraform 1.2.6 and earlier, this warning diagnostic is only
	// displayed when the Attribute is Required or Optional, and if the
	// practitioner configuration sets the value to a known or unknown value
	// (which may eventually be null). It has no effect when the Attribute is
	// Computed-only (read-only; not Required or Optional).
	//
	// Across any Terraform version, there are no warnings raised for
	// practitioner configuration values set directly to null, as there is no
	// way for the framework to differentiate between an unset and null
	// configuration due to how Terraform sends configuration information
	// across the protocol.
	//
	// Additional information about deprecation enhancements for read-only
	// attributes can be found in:
	//
	//  - https://github.com/hashicorp/terraform/issues/7569
	//
	DeprecationMessage string

	// Validators define value validation functionality for the attribute. All
	// elements of the slice of AttributeValidator are run, regardless of any
	// previous error diagnostics.
	//
	// Many common use case validators can be found in the
	// github.com/hashicorp/terraform-plugin-framework-validators Go module.
	//
	// If the Type field points to a custom type that implements the
	// xattr.TypeWithValidate interface, the validators defined in this field
	// are run in addition to the validation defined by the type.
	Validators []validator.Dynamic

	// PlanModifiers defines a sequence of modifiers for this attribute at
	// plan time. Schema-based plan modifications occur before any
	// resource-level plan modifications.
	//
	// Schema-based plan modifications can adjust Terraform's plan by:
	//
	//  - Requiring resource recreation. Typically used for configuration
	//    updates which cannot be done in-place.
	//  - Setting the planned value. Typically used for enhancing the plan
	//    to replace unknown values. Computed must be true or Terraform will
	//    return an error. If the plan value is known due to a known
	//    configuration value, the plan value cannot be changed or Terraform
	//    will return an error.
	//
	// Any errors will prevent further execution of this sequence or modifiers.
	PlanModifiers []planmodifier.Dynamic

	// Default defines a proposed new state (plan) value for the attribute
	// if the configuration value is null. Default prevents the framework
	// from automatically marking the value as unknown during planning when
	// other proposed new state changes are detected. If the attribute is
	// computed and the value could be altered by other changes then a default
	// should be avoided and a plan modifier should be used instead.
	Default defaults.Dynamic
}

// ApplyTerraform5AttributePathStep always returns an error as it is not
// possible to step further into a DynamicAttribute.
func (a DynamicAttribute) ApplyTerraform5AttributePathStep(step tftypes.AttributePathStep) (interface{}, error) {
	return a.GetType().ApplyTerraform5AttributePathStep(step)
}

// Equal returns true if the given Attribute is a DynamicAttribute
// and all fields are equal.
func (a DynamicAttribute) Equal(o fwschema.Attribute) bool {
	if _, ok := o.(DynamicAttribute); !ok {
		return false
	}

	return fwschema.AttributesEqual(a, o)
}

// GetDeprecationMessage returns the DeprecationMessage field value.
func (a DynamicAttribute) GetDeprecationMessage() string {
	return a.DeprecationMessage
}

// GetDescription returns the Description field value.
func (a DynamicAttribute) GetDescription() string {
	return a.Description
}

// GetMarkdownDescription returns the MarkdownDescription field value.
func (a DynamicAttribute) GetMarkdownDescription() string {
	return a.MarkdownDescription
}

// GetType returns types.DynamicType or the CustomType field value if defined.
func (a DynamicAttribute) GetType() attr.Type {
	if a.CustomType != nil {
		return a.CustomType
	}

	return types.DynamicType
}

// IsComputed returns the Computed field value.
func (a DynamicAttribute) IsComputed() bool {
	return a.Computed
}

// IsOptional returns the Optional field value.
func (a DynamicAttribute) IsOptional() bool {
	return a.Optional
}

// IsRequired returns the Required field value.
func (a DynamicAttribute) IsRequired() bool {
	return a.Required
}

// IsSensitive returns the Sensitive field value.
func (a DynamicAttribute) IsSensitive() bool {
	return a.Sensitive
}

// DynamicDefaultValue returns the Default field value.
func (a DynamicAttribute) DynamicDefaultValue() defaults.Dynamic {
	return a.Default
}

// DynamicPlanModifiers returns the PlanModifiers field value.
func (a DynamicAttribute) DynamicPlanModifiers() []planmodifier.Dynamic {
	return a.PlanModifiers
}

// DynamicValidators returns the Validators field value.
func (a DynamicAttribute) DynamicValidators() []validator.Dynamic {
	return a.Validators
}

// ValidateImplementation contains logic for validating the
// provider-defined implementation of the attribute to prevent unexpected
// errors or panics. This logic runs during the GetProviderSchema RPC and
// should never include false positives.
func (a DynamicAttribute) ValidateImplementation(ctx context.Context, req fwschema.ValidateImplementationRequest, resp *fwschema.ValidateImplementationResponse) {
	if !a.IsComputed() && a.DynamicDefaultValue() != nil {
		resp.Diagnostics.Append(nonComputedAttributeWithDefaultDiag(req.Path))
	}
}
