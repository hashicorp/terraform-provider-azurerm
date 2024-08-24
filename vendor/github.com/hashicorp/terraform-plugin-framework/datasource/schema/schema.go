// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

// Schema must satify the fwschema.Schema interface.
var _ fwschema.Schema = Schema{}

// Schema defines the structure and value types of data source data. This type
// is used as the datasource.SchemaResponse type Schema field, which is
// implemented by the datasource.DataSource type Schema method.
type Schema struct {
	// Attributes is the mapping of underlying attribute names to attribute
	// definitions.
	//
	// Names must only contain lowercase letters, numbers, and underscores.
	// Names must not collide with any Blocks names.
	Attributes map[string]Attribute

	// Blocks is the mapping of underlying block names to block definitions.
	//
	// Names must only contain lowercase letters, numbers, and underscores.
	// Names must not collide with any Attributes names.
	Blocks map[string]Block

	// Description is used in various tooling, like the language server, to
	// give practitioners more information about what this data source is,
	// what it's for, and how it should be used. It should be written as
	// plain text, with no special formatting.
	Description string

	// MarkdownDescription is used in various tooling, like the
	// documentation generator, to give practitioners more information
	// about what this data source is, what it's for, and how it should be
	// used. It should be formatted using Markdown.
	MarkdownDescription string

	// DeprecationMessage defines warning diagnostic details to display when
	// practitioner configurations use this data source. The warning diagnostic
	// summary is automatically set to "Data Source Deprecated" along with
	// configuration source file and line information.
	//
	// Set this field to a practitioner actionable message such as:
	//
	//  - "Use examplecloud_other data source instead. This data source
	//    will be removed in the next major version of the provider."
	//  - "Remove this data source as it no longer is valid and
	//    will be removed in the next major version of the provider."
	//
	DeprecationMessage string
}

// ApplyTerraform5AttributePathStep applies the given AttributePathStep to the
// schema.
func (s Schema) ApplyTerraform5AttributePathStep(step tftypes.AttributePathStep) (any, error) {
	return fwschema.SchemaApplyTerraform5AttributePathStep(s, step)
}

// AttributeAtPath returns the Attribute at the passed path. If the path points
// to an element or attribute of a complex type, rather than to an Attribute,
// it will return an ErrPathInsideAtomicAttribute error.
func (s Schema) AttributeAtPath(ctx context.Context, p path.Path) (fwschema.Attribute, diag.Diagnostics) {
	return fwschema.SchemaAttributeAtPath(ctx, s, p)
}

// AttributeAtPath returns the Attribute at the passed path. If the path points
// to an element or attribute of a complex type, rather than to an Attribute,
// it will return an ErrPathInsideAtomicAttribute error.
func (s Schema) AttributeAtTerraformPath(ctx context.Context, p *tftypes.AttributePath) (fwschema.Attribute, error) {
	return fwschema.SchemaAttributeAtTerraformPath(ctx, s, p)
}

// GetAttributes returns the Attributes field value.
func (s Schema) GetAttributes() map[string]fwschema.Attribute {
	return schemaAttributes(s.Attributes)
}

// GetBlocks returns the Blocks field value.
func (s Schema) GetBlocks() map[string]fwschema.Block {
	return schemaBlocks(s.Blocks)
}

// GetDeprecationMessage returns the DeprecationMessage field value.
func (s Schema) GetDeprecationMessage() string {
	return s.DeprecationMessage
}

// GetDescription returns the Description field value.
func (s Schema) GetDescription() string {
	return s.Description
}

// GetMarkdownDescription returns the MarkdownDescription field value.
func (s Schema) GetMarkdownDescription() string {
	return s.MarkdownDescription
}

// GetVersion always returns 0 as data source schemas cannot be versioned.
func (s Schema) GetVersion() int64 {
	return 0
}

// Type returns the framework type of the schema.
func (s Schema) Type() attr.Type {
	return fwschema.SchemaType(s)
}

// TypeAtPath returns the framework type at the given schema path.
func (s Schema) TypeAtPath(ctx context.Context, p path.Path) (attr.Type, diag.Diagnostics) {
	return fwschema.SchemaTypeAtPath(ctx, s, p)
}

// TypeAtTerraformPath returns the framework type at the given tftypes path.
func (s Schema) TypeAtTerraformPath(ctx context.Context, p *tftypes.AttributePath) (attr.Type, error) {
	return fwschema.SchemaTypeAtTerraformPath(ctx, s, p)
}

// Validate verifies that the schema is not using a reserved field name for a top-level attribute.
//
// Deprecated: Use the ValidateImplementation method instead.
func (s Schema) Validate() diag.Diagnostics {
	return s.ValidateImplementation(context.Background())
}

// ValidateImplementation contains logic for validating the provider-defined
// implementation of the schema and underlying attributes and blocks to prevent
// unexpected errors or panics. This logic runs during the GetProviderSchema
// RPC, or via provider-defined unit testing, and should never include false
// positives.
func (s Schema) ValidateImplementation(ctx context.Context) diag.Diagnostics {
	var diags diag.Diagnostics

	for attributeName, attribute := range s.GetAttributes() {
		req := fwschema.ValidateImplementationRequest{
			Name: attributeName,
			Path: path.Root(attributeName),
		}

		diags.Append(fwschema.IsReservedResourceAttributeName(req.Name, req.Path)...)
		diags.Append(fwschema.ValidateAttributeImplementation(ctx, attribute, req)...)
	}

	for blockName, block := range s.GetBlocks() {
		req := fwschema.ValidateImplementationRequest{
			Name: blockName,
			Path: path.Root(blockName),
		}

		diags.Append(fwschema.IsReservedResourceAttributeName(req.Name, req.Path)...)
		diags.Append(fwschema.ValidateBlockImplementation(ctx, block, req)...)
	}

	return diags
}

// schemaAttributes is a datasource to fwschema type conversion function.
func schemaAttributes(attributes map[string]Attribute) map[string]fwschema.Attribute {
	result := make(map[string]fwschema.Attribute, len(attributes))

	for name, attribute := range attributes {
		result[name] = attribute
	}

	return result
}

// schemaBlocks is a datasource to fwschema type conversion function.
func schemaBlocks(blocks map[string]Block) map[string]fwschema.Block {
	result := make(map[string]fwschema.Block, len(blocks))

	for name, block := range blocks {
		result[name] = block
	}

	return result
}
