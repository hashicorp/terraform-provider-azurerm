// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mdparser

import (
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/data/models"
)

// itemType represents the type of markdown item
type itemType int

const (
	itemDefault itemType = iota
	itemHeader
	itemMeteInfo
	itemExample
	itemField
	itemBlockHead
	itemNote
	itemSeparator
	itemPlainText
)

// markItem represents a single markdown item with its type and content
type markItem struct {
	fromLine int
	toLine   int
	lines    []string
	itemType itemType
	field    *models.DocumentProperty
}

func (m *markItem) content() string {
	return strings.Join(m.lines, "\n")
}

func (m *markItem) addLine(num int, line string) {
	m.lines = append(m.lines, line)
	m.toLine = num
}

func newMarkItem(fromLine int, content string, typ itemType) *markItem {
	return &markItem{
		fromLine: fromLine,
		lines:    []string{content},
		itemType: typ,
	}
}

// markBlock represents a block definition in the documentation
type markBlock struct {
	Names    []string
	Of       string
	Name     string
	HeadLine int
	Fields   []*models.DocumentProperty
	asProp   *models.DocumentProperties
}

func (b *markBlock) asProperties() *models.DocumentProperties {
	if b.asProp == nil {
		res := models.NewDocumentProperties()
		for _, f := range b.Fields {
			if _, ok := res.Objects[f.Name]; ok {
				if f.ParseErrors == nil {
					f.ParseErrors = []string{}
				}
				f.ParseErrors = append(f.ParseErrors, DuplicateFieldsFound)
			}
			res.Objects[f.Name] = f
			res.Names = append(res.Names, f.Name)
		}
		b.asProp = res
	}
	return b.asProp
}

func (b *markBlock) addField(f *models.DocumentProperty) {
	b.Fields = append(b.Fields, f)
}

// mark represents the parsed markdown structure
type mark struct {
	items  []*markItem
	blocks []markBlock
	fields map[string]*models.DocumentProperty
}

func (m *mark) lastItem() *markItem {
	if len(m.items) > 0 {
		return m.items[len(m.items)-1]
	}
	return nil
}

func (m *mark) addItem(item *markItem) {
	m.items = append(m.items, item)
}

func (m *mark) addLineOrItem(num int, line string, typ itemType) {
	last := m.lastItem()
	if last != nil && last.itemType == typ {
		last.addLine(num, line)
	} else {
		m.addItem(newMarkItem(num, line, typ))
	}
}

func (m *mark) addBlock(b markBlock) {
	m.blocks = append(m.blocks, b)
}

// newMarkFromString performs Phase 1: Parse markdown into structured items (tokenization)
func newMarkFromString(content string) *mark {
	lines := strings.Split(content, "\n")
	result := &mark{
		fields: map[string]*models.DocumentProperty{},
	}

	for idx, line := range lines {
		switch {
		case strings.HasPrefix(line, "###"):
			result.addItem(newMarkItem(idx, line, itemHeader))
			continue
		case strings.HasPrefix(line, "##"):
			result.addItem(newMarkItem(idx, line, itemHeader))
			continue
		case strings.HasPrefix(line, "#"):
			result.addItem(newMarkItem(idx, line, itemHeader))
			continue
		case strings.HasPrefix(line, "*"):
			result.addItem(newMarkItem(idx, line, itemField))
		case strings.HasPrefix(line, "---"):
			if idx == 0 {
				result.addItem(newMarkItem(idx, line, itemMeteInfo)) // TODO: remove? this seems to be the example
			} else {
				last := result.lastItem()
				if last != nil && last.itemType == itemMeteInfo {
					last.addLine(idx, line)
				} else {
					result.addItem(newMarkItem(idx, line, itemSeparator))
				}
			}
		case strings.HasPrefix(line, "```"):
			result.addLineOrItem(idx, line, itemExample)
		case strings.HasPrefix(line, "->"), strings.HasPrefix(line, "~>"), strings.HasPrefix(line, "!>"):
			result.addItem(newMarkItem(idx, line, itemNote))
		case isBlockHead(line):
			result.addItem(newMarkItem(idx, line, itemBlockHead))
		default:
			// plain text
			last := result.lastItem()
			if last == nil {
				result.addItem(newMarkItem(idx, line, itemPlainText))
				continue
			}
			switch last.itemType {
			case itemField, itemMeteInfo, itemExample, itemPlainText:
				last.addLine(idx, line)
			default:
				if strings.TrimSpace(line) == "" {
					last.addLine(idx, line)
				} else {
					result.addItem(newMarkItem(idx, line, itemPlainText))
				}
			}
		}
	}

	return result
}
