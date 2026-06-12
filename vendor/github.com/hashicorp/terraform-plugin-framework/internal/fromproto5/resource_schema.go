// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fromproto5

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// ResourceSchema converts a *tfprotov5.Schema into a resource/schema Schema, used for
// converting protocol schemas (from another provider server, such as SDKv2 or terraform-plugin-go)
// into Framework schemas.
func ResourceSchema(ctx context.Context, s *tfprotov5.Schema) (*resourceschema.Schema, error) {
	if s == nil || s.Block == nil {
		return nil, nil
	}

	attrs, err := ResourceSchemaAttributes(ctx, s.Block.Attributes)
	if err != nil {
		return nil, err
	}

	blocks, err := ResourceSchemaNestedBlocks(ctx, s.Block.BlockTypes)
	if err != nil {
		return nil, err
	}

	return &resourceschema.Schema{
		// MAINTAINER NOTE: At the moment, there isn't a need to copy all of the data from the protocol schema
		// to the resource schema, just enough data to allow provider developers to read and set data.
		Attributes: attrs,
		Blocks:     blocks,
	}, nil
}

func ResourceSchemaAttributes(ctx context.Context, protoAttrs []*tfprotov5.SchemaAttribute) (map[string]resourceschema.Attribute, error) {
	attrs := make(map[string]resourceschema.Attribute, len(protoAttrs))
	for _, protoAttr := range protoAttrs {
		// MAINTAINER NOTE: At the moment, there isn't a need to copy all of the data from the protocol schema
		// to the resource schema, just enough data to allow provider developers to read and set data.
		switch {
		case protoAttr.Type.Is(tftypes.Bool):
			attrs[protoAttr.Name] = resourceschema.BoolAttribute{
				Required:  protoAttr.Required,
				Optional:  protoAttr.Optional,
				Computed:  protoAttr.Computed,
				WriteOnly: protoAttr.WriteOnly,
				Sensitive: protoAttr.Sensitive,
			}
		case protoAttr.Type.Is(tftypes.Number):
			attrs[protoAttr.Name] = resourceschema.NumberAttribute{
				Required:  protoAttr.Required,
				Optional:  protoAttr.Optional,
				Computed:  protoAttr.Computed,
				WriteOnly: protoAttr.WriteOnly,
				Sensitive: protoAttr.Sensitive,
			}
		case protoAttr.Type.Is(tftypes.String):
			attrs[protoAttr.Name] = resourceschema.StringAttribute{
				Required:  protoAttr.Required,
				Optional:  protoAttr.Optional,
				Computed:  protoAttr.Computed,
				WriteOnly: protoAttr.WriteOnly,
				Sensitive: protoAttr.Sensitive,
			}
		case protoAttr.Type.Is(tftypes.DynamicPseudoType):
			attrs[protoAttr.Name] = resourceschema.DynamicAttribute{
				Required:  protoAttr.Required,
				Optional:  protoAttr.Optional,
				Computed:  protoAttr.Computed,
				WriteOnly: protoAttr.WriteOnly,
				Sensitive: protoAttr.Sensitive,
			}
		case protoAttr.Type.Is(tftypes.List{}):
			//nolint:forcetypeassert // Type assertion is guaranteed by the above `(tftypes.Type).Is` function
			l := protoAttr.Type.(tftypes.List)

			elementType, err := basetypes.TerraformTypeToFrameworkType(l.ElementType)
			if err != nil {
				return nil, err
			}

			attrs[protoAttr.Name] = resourceschema.ListAttribute{
				ElementType: elementType,
				Required:    protoAttr.Required,
				Optional:    protoAttr.Optional,
				Computed:    protoAttr.Computed,
				WriteOnly:   protoAttr.WriteOnly,
				Sensitive:   protoAttr.Sensitive,
			}
		case protoAttr.Type.Is(tftypes.Map{}):
			//nolint:forcetypeassert // Type assertion is guaranteed by the above `(tftypes.Type).Is` function
			m := protoAttr.Type.(tftypes.Map)

			elementType, err := basetypes.TerraformTypeToFrameworkType(m.ElementType)
			if err != nil {
				return nil, err
			}

			attrs[protoAttr.Name] = resourceschema.MapAttribute{
				ElementType: elementType,
				Required:    protoAttr.Required,
				Optional:    protoAttr.Optional,
				Computed:    protoAttr.Computed,
				WriteOnly:   protoAttr.WriteOnly,
				Sensitive:   protoAttr.Sensitive,
			}
		case protoAttr.Type.Is(tftypes.Set{}):
			//nolint:forcetypeassert // Type assertion is guaranteed by the above `(tftypes.Type).Is` function
			s := protoAttr.Type.(tftypes.Set)

			elementType, err := basetypes.TerraformTypeToFrameworkType(s.ElementType)
			if err != nil {
				return nil, err
			}

			attrs[protoAttr.Name] = resourceschema.SetAttribute{
				ElementType: elementType,
				Required:    protoAttr.Required,
				Optional:    protoAttr.Optional,
				Computed:    protoAttr.Computed,
				Sensitive:   protoAttr.Sensitive,
			}
		case protoAttr.Type.Is(tftypes.Object{}):
			//nolint:forcetypeassert // Type assertion is guaranteed by the above `(tftypes.Type).Is` function
			o := protoAttr.Type.(tftypes.Object)

			attrTypes := make(map[string]attr.Type, len(o.AttributeTypes))
			for name, tfType := range o.AttributeTypes {
				t, err := basetypes.TerraformTypeToFrameworkType(tfType)
				if err != nil {
					return nil, err
				}
				attrTypes[name] = t
			}

			attrs[protoAttr.Name] = resourceschema.ObjectAttribute{
				AttributeTypes: attrTypes,
				Required:       protoAttr.Required,
				Optional:       protoAttr.Optional,
				Computed:       protoAttr.Computed,
				WriteOnly:      protoAttr.WriteOnly,
				Sensitive:      protoAttr.Sensitive,
			}
		default:
			// MAINTAINER NOTE: Currently the only type not supported by Framework is a tuple, since there
			// is no corresponding attribute to represent it.
			//
			// https://github.com/hashicorp/terraform-plugin-framework/issues/54
			return nil, fmt.Errorf("no supported attribute for %q, type: %T", protoAttr.Name, protoAttr.Type)
		}
	}

	return attrs, nil
}

func ResourceSchemaNestedBlocks(ctx context.Context, protoBlocks []*tfprotov5.SchemaNestedBlock) (map[string]resourceschema.Block, error) {
	nestedBlocks := make(map[string]resourceschema.Block, len(protoBlocks))
	for _, protoBlock := range protoBlocks {
		if protoBlock.Block == nil {
			continue
		}

		attrs, err := ResourceSchemaAttributes(ctx, protoBlock.Block.Attributes)
		if err != nil {
			return nil, err
		}
		blocks, err := ResourceSchemaNestedBlocks(ctx, protoBlock.Block.BlockTypes)
		if err != nil {
			return nil, err
		}

		switch protoBlock.Nesting {
		case tfprotov5.SchemaNestedBlockNestingModeList:
			nestedBlocks[protoBlock.TypeName] = resourceschema.ListNestedBlock{
				NestedObject: resourceschema.NestedBlockObject{
					Attributes: attrs,
					Blocks:     blocks,
				},
			}
		case tfprotov5.SchemaNestedBlockNestingModeSet:
			nestedBlocks[protoBlock.TypeName] = resourceschema.SetNestedBlock{
				NestedObject: resourceschema.NestedBlockObject{
					Attributes: attrs,
					Blocks:     blocks,
				},
			}
		case tfprotov5.SchemaNestedBlockNestingModeSingle:
			nestedBlocks[protoBlock.TypeName] = resourceschema.SingleNestedBlock{
				Attributes: attrs,
				Blocks:     blocks,
			}
		default:
			// MAINTAINER NOTE: Currently the only block type not supported by Framework is a map nested block, since there
			// is no corresponding framework block implementation to represent it.
			return nil, fmt.Errorf("no supported block for nesting mode %v in nested block %q", protoBlock.Nesting, protoBlock.TypeName)
		}
	}

	return nestedBlocks, nil
}
