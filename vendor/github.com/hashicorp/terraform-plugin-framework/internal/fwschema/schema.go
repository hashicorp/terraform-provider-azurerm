// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwschema

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/totftypes"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Schema is the core interface required for data sources, providers, and
// resources.
type Schema interface {
	// Implementations should include the tftypes.AttributePathStepper
	// interface methods for proper path and data handling.
	tftypes.AttributePathStepper

	// AttributeAtPath should return the Attribute at the given path or return
	// an error.
	AttributeAtPath(context.Context, path.Path) (Attribute, diag.Diagnostics)

	// AttributeAtTerraformPath should return the Attribute at the given
	// Terraform path or return an error.
	AttributeAtTerraformPath(context.Context, *tftypes.AttributePath) (Attribute, error)

	// GetAttributes should return the attributes of a schema. This is named
	// differently than Attributes to prevent a conflict with the tfsdk.Schema
	// field name.
	GetAttributes() map[string]Attribute

	// GetBlocks should return the blocks of a schema. This is named
	// differently than Blocks to prevent a conflict with the tfsdk.Schema
	// field name.
	GetBlocks() map[string]Block

	// GetDeprecationMessage should return a non-empty string if a schema
	// is deprecated. This is named differently than DeprecationMessage to
	// prevent a conflict with the tfsdk.Schema field name.
	GetDeprecationMessage() string

	// GetDescription should return a non-empty string if a schema has a
	// plaintext description. This is named differently than Description
	// to prevent a conflict with the tfsdk.Schema field name.
	GetDescription() string

	// GetMarkdownDescription should return a non-empty string if a schema has
	// a Markdown description. This is named differently than
	// MarkdownDescription to prevent a conflict with the tfsdk.Schema field
	// name.
	GetMarkdownDescription() string

	// GetVersion should return the version of a schema. This is named
	// differently than Version to prevent a conflict with the tfsdk.Schema
	// field name.
	GetVersion() int64

	// Type should return the framework type of the schema.
	Type() attr.Type

	// TypeAtPath should return the framework type of the Attribute at the
	// the given path or return an error.
	TypeAtPath(context.Context, path.Path) (attr.Type, diag.Diagnostics)

	// AttributeTypeAtPath should return the framework type of the Attribute at
	// the given Terraform path or return an error.
	TypeAtTerraformPath(context.Context, *tftypes.AttributePath) (attr.Type, error)
}

// SchemaApplyTerraform5AttributePathStep is a helper function to perform base
// tftypes.AttributePathStepper handling using the GetAttributes and GetBlocks
// methods.
func SchemaApplyTerraform5AttributePathStep(s Schema, step tftypes.AttributePathStep) (any, error) {
	name, ok := step.(tftypes.AttributeName)

	if !ok {
		return nil, fmt.Errorf("cannot apply AttributePathStep %T to schema", step)
	}

	if attr, ok := s.GetAttributes()[string(name)]; ok {
		return attr, nil
	}

	if block, ok := s.GetBlocks()[string(name)]; ok {
		return block, nil
	}

	return nil, fmt.Errorf("could not find attribute or block %q in schema", name)
}

// SchemaAttributeAtPath is a helper function to perform base type handling using
// the AttributeAtTerraformPath method.
func SchemaAttributeAtPath(ctx context.Context, s Schema, p path.Path) (Attribute, diag.Diagnostics) {
	var diags diag.Diagnostics

	tftypesPath, tftypesDiags := totftypes.AttributePath(ctx, p)

	diags.Append(tftypesDiags...)

	if diags.HasError() {
		return nil, diags
	}

	attribute, err := s.AttributeAtTerraformPath(ctx, tftypesPath)

	if err != nil {
		diags.AddAttributeError(
			p,
			"Invalid Schema Path",
			"When attempting to get the framework attribute associated with a schema path, an unexpected error was returned. "+
				"This is always an issue with the provider. Please report this to the provider developers.\n\n"+
				fmt.Sprintf("Path: %s\n", p)+
				fmt.Sprintf("Original Error: %s", err),
		)
		return nil, diags
	}

	return attribute, diags
}

// SchemaAttributeAtTerraformPath is a helper function to perform base type
// handling using the tftypes.AttributePathStepper interface.
func SchemaAttributeAtTerraformPath(ctx context.Context, s Schema, p *tftypes.AttributePath) (Attribute, error) {
	rawType, remaining, err := tftypes.WalkAttributePath(s, p)

	if err != nil {
		return nil, checkErrForDynamic(rawType, remaining, err)
	}

	switch typ := rawType.(type) {
	case attr.Type:
		return nil, ErrPathInsideAtomicAttribute
	case Attribute:
		return typ, nil
	case Block:
		return nil, ErrPathIsBlock
	case NestedAttributeObject:
		return nil, ErrPathInsideAtomicAttribute
	case NestedBlockObject:
		return nil, ErrPathInsideAtomicAttribute
	case UnderlyingAttributes:
		return nil, ErrPathInsideAtomicAttribute
	default:
		return nil, fmt.Errorf("got unexpected type %T", rawType)
	}
}

// SchemaBlockPathExpressions returns a slice of all path expressions which
// represent a Block according to the Schema.
func SchemaBlockPathExpressions(ctx context.Context, s Schema) path.Expressions {
	result := path.Expressions{}

	for name, block := range s.GetBlocks() {
		result = append(result, BlockPathExpressions(ctx, block, path.MatchRoot(name))...)
	}

	return result
}

// SchemaTypeAtPath is a helper function to perform base type handling using
// the TypeAtTerraformPath method.
func SchemaTypeAtPath(ctx context.Context, s Schema, p path.Path) (attr.Type, diag.Diagnostics) {
	var diags diag.Diagnostics

	tftypesPath, tftypesDiags := totftypes.AttributePath(ctx, p)

	diags.Append(tftypesDiags...)

	if diags.HasError() {
		return nil, diags
	}

	attrType, err := s.TypeAtTerraformPath(ctx, tftypesPath)

	if err != nil {
		diags.AddAttributeError(
			p,
			"Invalid Schema Path",
			"When attempting to get the framework type associated with a schema path, an unexpected error was returned. "+
				"This is always an issue with the provider. Please report this to the provider developers.\n\n"+
				fmt.Sprintf("Path: %s\n", p)+
				fmt.Sprintf("Original Error: %s", err),
		)
		return nil, diags
	}

	return attrType, diags
}

// SchemaTypeAtTerraformPath is a helper function to perform base type handling
// using the tftypes.AttributePathStepper interface.
func SchemaTypeAtTerraformPath(ctx context.Context, s Schema, p *tftypes.AttributePath) (attr.Type, error) {
	rawType, remaining, err := tftypes.WalkAttributePath(s, p)

	if err != nil {
		return nil, checkErrForDynamic(rawType, remaining, err)
	}

	switch typ := rawType.(type) {
	case attr.Type:
		return typ, nil
	case Attribute:
		return typ.GetType(), nil
	case Block:
		return typ.Type(), nil
	case NestedAttributeObject:
		return typ.Type(), nil
	case NestedBlockObject:
		return typ.Type(), nil
	case Schema:
		return typ.Type(), nil
	case UnderlyingAttributes:
		return typ.Type(), nil
	default:
		return nil, fmt.Errorf("got unexpected type %T", rawType)
	}
}

// SchemaType is a helper function to perform base type handling using the
// GetAttributes and GetBlocks methods.
func SchemaType(s Schema) attr.Type {
	attrTypes := map[string]attr.Type{}

	for name, attr := range s.GetAttributes() {
		attrTypes[name] = attr.GetType()
	}

	for name, block := range s.GetBlocks() {
		attrTypes[name] = block.Type()
	}

	return types.ObjectType{AttrTypes: attrTypes}
}

// checkErrForDynamic is a helper function that will always return an error. It will return
// an `ErrPathInsideDynamicAttribute` error if rawType:
//   - Is a dynamic type
//   - Is an attribute that has a dynamic type
func checkErrForDynamic(rawType any, remaining *tftypes.AttributePath, err error) error {
	if rawType == nil {
		return fmt.Errorf("%v still remains in the path: %w", remaining, err)
	}

	// Check to see if we tried walking into a dynamic type (types.DynamicType)
	_, isDynamic := rawType.(basetypes.DynamicTypable)
	if isDynamic {
		// If the type is dynamic there is no schema information underneath it, return an error to allow calling logic to safely skip
		return ErrPathInsideDynamicAttribute
	}

	// Check to see if we tried walking into an attribute with a dynamic type (schema.DynamicAttribute)
	attr, ok := rawType.(Attribute)
	if ok {
		_, isDynamic := attr.GetType().(basetypes.DynamicTypable)
		if isDynamic {
			// If the attribute is dynamic there are no nested attributes underneath it, return an error to allow calling logic to safely skip
			return ErrPathInsideDynamicAttribute
		}
	}

	return fmt.Errorf("%v still remains in the path: %w", remaining, err)
}
