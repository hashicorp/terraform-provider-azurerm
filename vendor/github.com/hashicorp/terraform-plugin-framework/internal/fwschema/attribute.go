// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwschema

import (
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-plugin-framework/attr"
)

// Attribute is the core interface required for implementing Terraform
// schema functionality that can accept a value. Refer to NestedAttribute for
// the additional interface that defines nested attributes.
//
// Refer to the internal/fwschema/fwxschema package for optional interfaces
// that define framework-specific functionality, such a plan modification and
// validation.
type Attribute interface {
	// Implementations should include the tftypes.AttributePathStepper
	// interface methods for proper path and data handling.
	tftypes.AttributePathStepper

	// Equal should return true if the other attribute is exactly equivalent.
	Equal(o Attribute) bool

	// GetDeprecationMessage should return a non-empty string if an attribute
	// is deprecated. This is named differently than DeprecationMessage to
	// prevent a conflict with the tfsdk.Attribute field name.
	GetDeprecationMessage() string

	// GetDescription should return a non-empty string if an attribute
	// has a plaintext description. This is named differently than Description
	// to prevent a conflict with the tfsdk.Attribute field name.
	GetDescription() string

	// GetMarkdownDescription should return a non-empty string if an attribute
	// has a Markdown description. This is named differently than
	// MarkdownDescription to prevent a conflict with the tfsdk.Attribute field
	// name.
	GetMarkdownDescription() string

	// GetType should return the framework type of an attribute. This is named
	// differently than Type to prevent a conflict with the tfsdk.Attribute
	// field name.
	GetType() attr.Type

	// IsComputed should return true if the attribute configuration value is
	// computed. This is named differently than Computed to prevent a conflict
	// with the tfsdk.Attribute field name.
	IsComputed() bool

	// IsOptional should return true if the attribute configuration value is
	// optional. This is named differently than Optional to prevent a conflict
	// with the tfsdk.Attribute field name.
	IsOptional() bool

	// IsRequired should return true if the attribute configuration value is
	// required. This is named differently than Required to prevent a conflict
	// with the tfsdk.Attribute field name.
	IsRequired() bool

	// IsSensitive should return true if the attribute configuration value is
	// sensitive. This is named differently than Sensitive to prevent a
	// conflict with the tfsdk.Attribute field name.
	IsSensitive() bool

	// IsWriteOnly should return true if the attribute configuration value is
	// write-only. This is named differently than WriteOnly to prevent a
	// conflict with the tfsdk.Attribute field name.
	//
	// Write-only attributes are a managed-resource schema concept only.
	IsWriteOnly() bool

	// IsOptionalForImport should return true if the identity attribute is optional to be set by
	// the practitioner when importing by identity. This is named differently than OptionalForImport
	// to prevent a conflict with the relevant field name.
	IsOptionalForImport() bool

	// IsRequiredForImport should return true if the identity attribute must be set by
	// the practitioner when importing by identity. This is named differently than RequiredForImport
	// to prevent a conflict with the relevant field name.
	IsRequiredForImport() bool
}

// AttributesEqual is a helper function to perform equality testing on two
// Attribute. Attribute Equal implementations should still compare the concrete
// types in addition to using this helper.
func AttributesEqual(a, b Attribute) bool {
	if !a.GetType().Equal(b.GetType()) {
		return false
	}

	if a.GetDeprecationMessage() != b.GetDeprecationMessage() {
		return false
	}

	if a.GetDescription() != b.GetDescription() {
		return false
	}

	if a.GetMarkdownDescription() != b.GetMarkdownDescription() {
		return false
	}

	if a.IsRequired() != b.IsRequired() {
		return false
	}

	if a.IsOptional() != b.IsOptional() {
		return false
	}

	if a.IsComputed() != b.IsComputed() {
		return false
	}

	if a.IsSensitive() != b.IsSensitive() {
		return false
	}

	if a.IsWriteOnly() != b.IsWriteOnly() {
		return false
	}

	if a.IsOptionalForImport() != b.IsOptionalForImport() {
		return false
	}

	if a.IsRequiredForImport() != b.IsRequiredForImport() {
		return false
	}

	return true
}
