// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto6

import (
	"context"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// SchemaAttribute returns the *tfprotov6.SchemaAttribute equivalent of an
// Attribute. Errors will be tftypes.AttributePathErrors based on `path`.
// `name` is the name of the attribute.
func SchemaAttribute(ctx context.Context, name string, path *tftypes.AttributePath, a fwschema.Attribute) (*tfprotov6.SchemaAttribute, error) {
	if !a.IsRequired() && !a.IsOptional() && !a.IsComputed() {
		return nil, path.NewErrorf("must have Required, Optional, or Computed set")
	}

	schemaAttribute := &tfprotov6.SchemaAttribute{
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
		schemaAttribute.DescriptionKind = tfprotov6.StringKindPlain
	}

	if a.GetMarkdownDescription() != "" {
		schemaAttribute.Description = a.GetMarkdownDescription()
		schemaAttribute.DescriptionKind = tfprotov6.StringKindMarkdown
	}

	nestedAttribute, ok := a.(fwschema.NestedAttribute)

	if !ok {
		return schemaAttribute, nil
	}

	object := &tfprotov6.SchemaObject{}
	nm := nestedAttribute.GetNestingMode()
	switch nm {
	case fwschema.NestingModeSingle:
		object.Nesting = tfprotov6.SchemaObjectNestingModeSingle
	case fwschema.NestingModeList:
		object.Nesting = tfprotov6.SchemaObjectNestingModeList
	case fwschema.NestingModeSet:
		object.Nesting = tfprotov6.SchemaObjectNestingModeSet
	case fwschema.NestingModeMap:
		object.Nesting = tfprotov6.SchemaObjectNestingModeMap
	default:
		return nil, path.NewErrorf("unrecognized nesting mode %v", nm)
	}

	for nestedName, nestedA := range nestedAttribute.GetNestedObject().GetAttributes() {
		nestedSchemaAttribute, err := SchemaAttribute(ctx, nestedName, path.WithAttributeName(nestedName), nestedA)

		if err != nil {
			return nil, err
		}

		object.Attributes = append(object.Attributes, nestedSchemaAttribute)
	}

	sort.Slice(object.Attributes, func(i, j int) bool {
		if object.Attributes[i] == nil {
			return true
		}

		if object.Attributes[j] == nil {
			return false
		}

		return object.Attributes[i].Name < object.Attributes[j].Name
	})

	schemaAttribute.NestedType = object
	schemaAttribute.Type = nil

	return schemaAttribute, nil
}
