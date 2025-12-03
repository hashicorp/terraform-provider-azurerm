// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mdparser

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/data/models"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/util"
)

const (
	BlcokNotDefined        = "block is not defined in the documentation"
	IncorrectlyBlockMarked = "The document incorrectly implies this field is a block"
)

// ParseMarkdownSection is the main entry point for parsing markdown content into DocumentProperties.
// It performs a three-phase parsing pipeline:
// 1. Tokenization (via newMarkFromString)
// 2. Structure building (via buildField)
// 3. Linking (via buildStruct)
func ParseMarkdownSection(content []string) *models.DocumentProperties {
	// Join lines back into a single string for parsing
	fullContent := strings.Join(content, "\n")

	mark := newMarkFromString(fullContent)
	mark.buildField()
	mark.buildStruct()

	result := models.NewDocumentProperties()

	// Copy top-level fields
	for name, field := range mark.fields {
		result.Objects[name] = field
		result.Names = append(result.Names, name)
	}

	// Copy block definitions
	for _, block := range mark.blocks {
		if blockProp := convertBlockToProperty(&block); blockProp != nil {
			result.BlockDefinitions[block.Name] = blockProp
		}
	}

	return result
}

// convertBlockToProperty converts a markBlock to a DocumentProperty
func convertBlockToProperty(block *markBlock) *models.DocumentProperty {
	prop := &models.DocumentProperty{
		Name:   block.Name,
		Block:  true,
		Line:   block.HeadLine,
		Nested: block.asProperties(),
	}

	return prop
}

// buildField performs Phase 2: Build field and block structures from items
func (m *mark) buildField() {
	var inBlock bool
	var block markBlock

	for _, item := range m.items {
		content := item.content()
		switch item.itemType {
		case itemField:
			f := newFieldFromLine(content)
			f.Line = item.fromLine
			item.field = f
			if inBlock {
				block.addField(f)
			} else {
				if arg, ok := m.fields[f.Name]; ok {
					if arg.ParseErrors == nil {
						arg.ParseErrors = []string{}
					}
					arg.ParseErrors = append(arg.ParseErrors, "duplicate fields declared")
				} else {
					m.fields[f.Name] = f
				}
			}
		case itemBlockHead:
			if inBlock {
				m.addBlock(block)
			}
			names := extractBlockNames(item.lines[0])
			// of/within block
			var of string
			for _, sep := range []string{" of ", " within "} {
				if idx := strings.Index(content, sep); idx > 0 {
					of = util.FirstCodeValue(content[idx:])
				}
			}

			block = markBlock{
				Names:    names,
				Name:     names[0],
				Of:       of,
				HeadLine: item.fromLine,
			}
			inBlock = true
		case itemSeparator:
			if inBlock {
				m.addBlock(block)
			}
			inBlock = false
		}
	}

	if inBlock {
		m.addBlock(block)
	}
}

// buildStruct performs Phase 3: Link block-type fields to their definitions
func (m *mark) buildStruct() {
	fillField := func(f *models.DocumentProperty, parent string) {
		if f.Block {
			if b, msg := m.blockOfName(f.BlockTypeName, parent); b != nil {
				f.Nested = b.asProperties()
				if msg != "" {
					if f.ParseErrors == nil {
						f.ParseErrors = []string{}
					}
					f.ParseErrors = append(f.ParseErrors, msg)
				}
			} else {
				if b2, _ := m.blockOfName(f.Name, parent); b2 != nil {
					if f.ParseErrors == nil {
						f.ParseErrors = []string{}
					}
					f.ParseErrors = append(f.ParseErrors, fmt.Sprintf("misspell of name from `%s` to `%s`", f.Name, f.BlockTypeName))
				} else {
					if f.ParseErrors == nil {
						f.ParseErrors = []string{}
					}
					f.ParseErrors = append(f.ParseErrors, fmt.Sprintf("block `%s` not defined in documentation", f.Name))
				}
			}
		}
	}

	for _, f := range m.fields {
		fillField(f, "")
	}

	for _, b := range m.blocks {
		for _, f := range b.Fields {
			fillField(f, b.Name)
		}
	}
}

// blockOfName finds a block definition by name, optionally within a parent block
// Returns the block and an error message if there are issues (like duplicates)
func (m *mark) blockOfName(name string, parent string) (*markBlock, string) {
	var res []*markBlock
	for i := range m.blocks {
		b := &m.blocks[i]
		for _, n2 := range b.Names {
			if n2 == name {
				res = append(res, b)
			}
		}
	}

	if len(res) == 0 {
		return nil, ""
	}

	if parent != "" {
		for _, item := range res {
			if item.Of == parent {
				return item, ""
			}
		}
	}

	var msg string
	if len(res) > 1 {
		uniqueDefinitions := make(map[string]*markBlock)
		for _, block := range res {
			// Include the parent context (Of) in the key to distinguish blocks in different scopes
			key := fmt.Sprintf("%s:%s:%d", block.Name, block.Of, len(block.Fields))
			if existing, exists := uniqueDefinitions[key]; exists {
				if !blocksHaveSameDefinition(existing, block) {
					msg = fmt.Sprintf("duplicate block exists as name `%s`", name)
					break
				}
			} else {
				uniqueDefinitions[key] = block
			}
		}
	}
	return res[0], msg
}
