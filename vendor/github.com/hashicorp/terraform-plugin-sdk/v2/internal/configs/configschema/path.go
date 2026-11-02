// Copyright IBM Corp. 2019, 2026
// SPDX-License-Identifier: MPL-2.0

package configschema

import (
	"github.com/hashicorp/go-cty/cty"
)

// AttributeByPath looks up the Attribute schema which corresponds to the given
// cty.Path. A nil value is returned if the given path does not correspond to a
// specific attribute.
func (b *Block) AttributeByPath(path cty.Path) *Attribute {
	block := b
	for i, step := range path {
		switch step := step.(type) {
		case cty.GetAttrStep:
			if attr := block.Attributes[step.Name]; attr != nil {
				if i < len(path)-1 { // There's more to the path, but not more to this Attribute.
					return nil
				}
				return attr
			}

			if nestedBlock := block.BlockTypes[step.Name]; nestedBlock != nil {
				block = &nestedBlock.Block
				continue
			}

			return nil
		}
	}
	return nil
}

// BlockByPath looks up the Block schema which corresponds to the given
// cty.Path. A nil value is returned if the given path does not correspond to a
// specific block.
func (b *Block) BlockByPath(path cty.Path) *Block {
	for i, step := range path {
		switch step := step.(type) {
		case cty.GetAttrStep:
			if blockType := b.BlockTypes[step.Name]; blockType != nil {
				if len(blockType.Block.BlockTypes) > 0 && i < len(path)-1 {
					return blockType.Block.BlockByPath(path[i+1:])
				} else if i < len(path)-1 {
					return nil
				}
				return &blockType.Block
			}
		}
	}
	return nil
}
