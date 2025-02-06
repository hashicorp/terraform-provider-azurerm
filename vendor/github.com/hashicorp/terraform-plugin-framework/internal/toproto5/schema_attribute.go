// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto5

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// SchemaAttribute returns the *tfprotov5.SchemaAttribute equivalent of an
// Attribute. Errors will be tftypes.AttributePathErrors based on `path`.
// `name` is the name of the attribute.
func SchemaAttribute(ctx context.Context, name string, path *tftypes.AttributePath, a fwschema.Attribute) (*tfprotov5.SchemaAttribute, error) {
	if _, ok := a.(fwschema.NestedAttribute); ok {
		return nil, path.NewErrorf("protocol version 5 cannot have Attributes set")
	}

	if a.GetType() == nil {
		return nil, path.NewErrorf("must have Type set")
	}

	if !a.IsRequired() && !a.IsOptional() && !a.IsComputed() {
		return nil, path.NewErrorf("must have Required, Optional, or Computed set")
	}

	schemaAttribute := &tfprotov5.SchemaAttribute{
		Name:      name,
		Required:  a.IsRequired(),
		Optional:  a.IsOptional(),
		Computed:  a.IsComputed(),
		Sensitive: a.IsSensitive(),
		Type:      a.GetType().TerraformType(ctx),
	}

	if a.GetDeprecationMessage() != "" {
		schemaAttribute.Deprecated = true
	}

	if a.GetDescription() != "" {
		schemaAttribute.Description = a.GetDescription()
		schemaAttribute.DescriptionKind = tfprotov5.StringKindPlain
	}

	if a.GetMarkdownDescription() != "" {
		schemaAttribute.Description = a.GetMarkdownDescription()
		schemaAttribute.DescriptionKind = tfprotov5.StringKindMarkdown
	}

	return schemaAttribute, nil
}
