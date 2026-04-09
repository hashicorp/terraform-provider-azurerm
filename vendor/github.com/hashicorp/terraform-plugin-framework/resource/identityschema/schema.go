// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package identityschema

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

// Schema defines the structure and value types of resource identity data. This type
// is used as the resource.IdentitySchemaResponse type IdentitySchema field, which is
// implemented by the resource.ResourceWithIdentity type IdentitySchema method.
type Schema struct {
	// Attributes is the mapping of underlying attribute names to attribute
	// definitions.
	//
	// Names must only contain lowercase letters, numbers, and underscores.
	Attributes map[string]Attribute

	// Version indicates the current version of the resource identity schema. Resource
	// identity schema versioning enables identity state upgrades in conjunction with the
	// [resource.ResourceWithUpgradeIdentity] interface. Versioning is only
	// required if there is a breaking change involving existing identity state data,
	// such as changing an attribute type in a manner that is incompatible with the Terraform type.
	//
	// Versions are conventionally only incremented by one each release.
	Version int64
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

// GetBlocks returns an empty map as it's not relevant for identity schemas.
func (s Schema) GetBlocks() map[string]fwschema.Block {
	return map[string]fwschema.Block{}
}

// GetDeprecationMessage returns an empty string as identity schemas cannot
// surface deprecation messages.
func (s Schema) GetDeprecationMessage() string {
	return ""
}

// GetDescription returns an empty string as identity schemas cannot
// surface descriptions.
func (s Schema) GetDescription() string {
	return ""
}

// GetMarkdownDescription returns an empty string as identity schemas cannot
// surface descriptions.
func (s Schema) GetMarkdownDescription() string {
	return ""
}

// GetVersion returns the Version field value.
func (s Schema) GetVersion() int64 {
	return s.Version
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
// unexpected errors or panics. This logic runs during the
// GetResourceIdentitySchemas RPC, or via provider-defined unit testing, and should
// never include false positives.
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

	return diags
}

// schemaAttributes is a identity to fwschema type conversion function.
func schemaAttributes(attributes map[string]Attribute) map[string]fwschema.Attribute {
	result := make(map[string]fwschema.Attribute, len(attributes))

	for name, attribute := range attributes {
		result[name] = attribute
	}

	return result
}
