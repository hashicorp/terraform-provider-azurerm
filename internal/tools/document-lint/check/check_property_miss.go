// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package check

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/model"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/util"
)

type MissType int

func (m MissType) String() string {
	return []string{"ok", "documentation", "schema", "doc attribute", "block declare"}[m]
}

const (
	NotMiss MissType = iota
	MissInDoc
	MissInCode
	MissInDocAttr
	MissBlockDeclare // document block declares wrong-formatted
	Misspelling
	MissWrongPlace // field nested in the wrong block
)

type propertyMissDiff struct {
	checkBase
	MissType    MissType
	correctName string // for misspelling diff only
}

func newPropertyMiss(checkBase checkBase, missType MissType) *propertyMissDiff {
	return &propertyMissDiff{checkBase: checkBase, MissType: missType}
}

func (c propertyMissDiff) String() string {
	switch c.MissType {
	case MissBlockDeclare:
		return fmt.Sprintf("%s blocks should be declared like '%s'", c.checkBase.Str(), util.ItalicCode("One or more `xxx` block as defined below."))
	case MissWrongPlace:
		return fmt.Sprintf("%s should be nested in %s", c.checkBase.Str(), util.ItalicCode(util.XPathDir(c.correctName)))
	case Misspelling:
		// it can be in the wrong place
		return fmt.Sprintf("%s does not exist in the schema - should this be %s?", c.checkBase.Str(), util.FixedCode(util.XPathBase(c.correctName)))
	}
	return fmt.Sprintf("%s does not exist in the %s or is poorly formatted", c.checkBase.Str(), c.MissType)
}

func (c propertyMissDiff) Fix(line string) (result string, err error) {
	if c.MissType == Misspelling && c.correctName != "" {
		return strings.ReplaceAll(line,
			fmt.Sprintf("`%s`", c.mdField.Name),
			fmt.Sprintf("`%s`", util.XPathBase(c.correctName)),
		), nil
	}
	return line, nil
}

var _ Checker = (*propertyMissDiff)(nil)

func newMissItem(path string, f *model.Field, typ MissType) Checker {
	base := newCheckBase(0, path, f)
	if f != nil {
		base.line = f.Line
	}
	return newPropertyMiss(base, typ)
}

func newMissInCode(path string, f *model.Field) Checker {
	return newMissItem(path, f, MissInCode)
}

// miss in doc will fill a mock `f`
func newMissInDoc(path string, f *model.Field) Checker {
	return newMissItem(path, f, MissInDoc)
}

func newMissBlockDeclare(path string, f *model.Field) Checker {
	return newMissItem(path, f, MissBlockDeclare)
}

func newMisspelling(c *propertyMissDiff, d *propertyMissDiff) Checker {
	// special logic. if the base name equals, treat as a wrong placed field
	if util.XPathBase(c.checkBase.key) == util.XPathBase(d.checkBase.key) && strings.Contains(d.checkBase.key, ".") {
		item := newPropertyMiss(c.checkBase, MissWrongPlace)
		item.correctName = d.checkBase.key
		return item
	}
	item := newPropertyMiss(c.checkBase, Misspelling)
	item.correctName = d.checkBase.key
	return item
}

// missing in doc/code can be a misspelling in document. do have a check

func mergeMisspelling(checks []Checker) (res []Checker) {
	var missInDoc, missInCode []*propertyMissDiff
	for _, c := range checks {
		if p, ok := c.(*propertyMissDiff); ok {
			if p.MissType == MissInDoc {
				missInDoc = append(missInDoc, p)
			} else if p.MissType == MissInCode {
				missInCode = append(missInCode, p)
			}
		}
	}
	// check if missed name be like
	filterOut := map[*propertyMissDiff]struct{}{}
	for _, c := range missInCode {
		for _, d := range missInDoc {
			if dist := levenshteinDist(c.MDField().Name, d.mdField.Name); dist <= 3 {
				// if the edit distances less than 3, we think it's a misspelling
				filterOut[c] = struct{}{}
				filterOut[d] = struct{}{}
				res = append(res, newMisspelling(c, d))
			}
		}
	}
	for _, c := range checks {
		if miss, ok := c.(*propertyMissDiff); !ok {
			res = append(res, c)
		} else {
			if _, ok = filterOut[miss]; !ok {
				res = append(res, c)
			}
		}
	}
	return res
}

func levenshteinDist(str1, str2 string) int {
	column := make([]int, len(str1)+1)
	for y := 1; y <= len(str1); y++ {
		column[y] = y
	}

	for x := 1; x <= len(str2); x++ {
		column[0] = x
		lastKey := x - 1
		for y := 1; y <= len(str1); y++ {
			oldKey := column[y]
			var incr int
			if str1[y-1] != str2[x-1] {
				incr = 1
			}

			column[y] = minimumOf3(column[y]+1, column[y-1]+1, lastKey+incr)
			lastKey = oldKey
		}
	}
	return column[len(str1)]
}

func minimumOf3(a, b, c int) int {
	if a > b {
		a = b
	}
	if a > c {
		a = c
	}
	return a
}
