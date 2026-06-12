// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto6

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
)

// IdentitySchemaAttribute returns the *tfprotov6.ResourceIdentitySchemaAttribute equivalent of an
// Attribute. Errors will be tftypes.AttributePathErrors based on `path`. `name` is the name of the attribute.
func IdentitySchemaAttribute(ctx context.Context, name string, path *tftypes.AttributePath, a fwschema.Attribute) (*tfprotov6.ResourceIdentitySchemaAttribute, error) {
	if _, ok := a.(fwschema.NestedAttribute); ok {
		return nil, path.NewErrorf("identity schemas don't support NestedAttribute")
	}

	if a.GetType() == nil {
		return nil, path.NewErrorf("must have Type set")
	}

	if !a.IsRequiredForImport() && !a.IsOptionalForImport() {
		return nil, path.NewErrorf("must have RequiredForImport or OptionalForImport set")
	}

	identitySchemaAttribute := &tfprotov6.ResourceIdentitySchemaAttribute{
		Name:              name,
		RequiredForImport: a.IsRequiredForImport(),
		OptionalForImport: a.IsOptionalForImport(),
		Type:              a.GetType().TerraformType(ctx),

		// Unlike other schema attributes, identity attributes only have a single description field which
		// is assumed to be markdown. Both a.GetDescription() and a.GetMarkdownDescription() will return
		// the same string, so we just chose one here.
		Description: a.GetDescription(),
	}

	return identitySchemaAttribute, nil
}
