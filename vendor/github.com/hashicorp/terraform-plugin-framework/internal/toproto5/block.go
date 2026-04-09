// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto5

import (
	"context"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Block returns the *tfprotov5.SchemaNestedBlock equivalent of a Block.
// Errors will be tftypes.AttributePathErrors based on `path`. `name` is the
// name of the attribute.
func Block(ctx context.Context, name string, path *tftypes.AttributePath, b fwschema.Block) (*tfprotov5.SchemaNestedBlock, error) {
	schemaNestedBlock := &tfprotov5.SchemaNestedBlock{
		Block: &tfprotov5.SchemaBlock{
			Deprecated: b.GetDeprecationMessage() != "",
		},
		TypeName: name,
	}

	if b.GetDescription() != "" {
		schemaNestedBlock.Block.Description = b.GetDescription()
		schemaNestedBlock.Block.DescriptionKind = tfprotov5.StringKindPlain
	}

	if b.GetMarkdownDescription() != "" {
		schemaNestedBlock.Block.Description = b.GetMarkdownDescription()
		schemaNestedBlock.Block.DescriptionKind = tfprotov5.StringKindMarkdown
	}

	nm := b.GetNestingMode()
	switch nm {
	case fwschema.BlockNestingModeList:
		schemaNestedBlock.Nesting = tfprotov5.SchemaNestedBlockNestingModeList
	case fwschema.BlockNestingModeSet:
		schemaNestedBlock.Nesting = tfprotov5.SchemaNestedBlockNestingModeSet
	case fwschema.BlockNestingModeSingle:
		schemaNestedBlock.Nesting = tfprotov5.SchemaNestedBlockNestingModeSingle
	default:
		return nil, path.NewErrorf("unrecognized nesting mode %v", nm)
	}

	nestedBlockObject := b.GetNestedObject()

	for attrName, attr := range nestedBlockObject.GetAttributes() {
		attrPath := path.WithAttributeName(attrName)
		attrProto5, err := SchemaAttribute(ctx, attrName, attrPath, attr)

		if err != nil {
			return nil, err
		}

		schemaNestedBlock.Block.Attributes = append(schemaNestedBlock.Block.Attributes, attrProto5)
	}

	for blockName, block := range nestedBlockObject.GetBlocks() {
		blockPath := path.WithAttributeName(blockName)
		blockProto5, err := Block(ctx, blockName, blockPath, block)

		if err != nil {
			return nil, err
		}

		schemaNestedBlock.Block.BlockTypes = append(schemaNestedBlock.Block.BlockTypes, blockProto5)
	}

	sort.Slice(schemaNestedBlock.Block.Attributes, func(i, j int) bool {
		if schemaNestedBlock.Block.Attributes[i] == nil {
			return true
		}

		if schemaNestedBlock.Block.Attributes[j] == nil {
			return false
		}

		return schemaNestedBlock.Block.Attributes[i].Name < schemaNestedBlock.Block.Attributes[j].Name
	})

	sort.Slice(schemaNestedBlock.Block.BlockTypes, func(i, j int) bool {
		if schemaNestedBlock.Block.BlockTypes[i] == nil {
			return true
		}

		if schemaNestedBlock.Block.BlockTypes[j] == nil {
			return false
		}

		return schemaNestedBlock.Block.BlockTypes[i].TypeName < schemaNestedBlock.Block.BlockTypes[j].TypeName
	})

	return schemaNestedBlock, nil
}
