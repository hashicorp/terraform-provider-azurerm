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

// Schema returns the *tfprotov6.Schema equivalent of a Schema.
func Schema(ctx context.Context, s fwschema.Schema) (*tfprotov6.Schema, error) {
	if s == nil {
		return nil, nil
	}

	result := &tfprotov6.Schema{
		Version: s.GetVersion(),
	}

	var attrs []*tfprotov6.SchemaAttribute
	var blocks []*tfprotov6.SchemaNestedBlock

	for name, attr := range s.GetAttributes() {
		a, err := SchemaAttribute(ctx, name, tftypes.NewAttributePath().WithAttributeName(name), attr)

		if err != nil {
			return nil, err
		}

		attrs = append(attrs, a)
	}

	for name, block := range s.GetBlocks() {
		proto6, err := Block(ctx, name, tftypes.NewAttributePath().WithAttributeName(name), block)

		if err != nil {
			return nil, err
		}

		blocks = append(blocks, proto6)
	}

	sort.Slice(attrs, func(i, j int) bool {
		if attrs[i] == nil {
			return true
		}

		if attrs[j] == nil {
			return false
		}

		return attrs[i].Name < attrs[j].Name
	})

	sort.Slice(blocks, func(i, j int) bool {
		if blocks[i] == nil {
			return true
		}

		if blocks[j] == nil {
			return false
		}

		return blocks[i].TypeName < blocks[j].TypeName
	})

	result.Block = &tfprotov6.SchemaBlock{
		// core doesn't do anything with version, as far as I can tell,
		// so let's not set it.
		Attributes: attrs,
		BlockTypes: blocks,
		Deprecated: s.GetDeprecationMessage() != "",
	}

	if s.GetDescription() != "" {
		result.Block.Description = s.GetDescription()
		result.Block.DescriptionKind = tfprotov6.StringKindPlain
	}

	if s.GetMarkdownDescription() != "" {
		result.Block.Description = s.GetMarkdownDescription()
		result.Block.DescriptionKind = tfprotov6.StringKindMarkdown
	}

	return result, nil
}
