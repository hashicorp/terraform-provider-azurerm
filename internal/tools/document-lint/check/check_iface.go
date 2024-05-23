// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package check

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/model"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/util"
)

type Checker interface {
	Line() int
	Key() string // property key path
	ShouldSkip() bool
	MDField() *model.Field

	// String display diff item information for check
	String() string

	// Fix try to fix this issue with line. return the updated line
	Fix(line string) (result string, err error)
}

type checkBase struct {
	line    int
	key     string
	mdField *model.Field
}

func (c checkBase) ShouldSkip() bool {
	if c.line == 0 || c.MDField() == nil || c.MDField().Skip {
		return true
	}
	return false
}

func (i checkBase) Str() string {
	return fmt.Sprintf("%d %s", i.Line()+1, util.Bold(i.Key()))
}

func (i checkBase) Line() int {
	if i.line == 0 {
		return 1 // if one property missed in document, set to line 1 of doc
	}
	return i.line
}

func (i checkBase) Key() string {
	return i.key
}

func (i checkBase) MDField() *model.Field {
	return i.mdField
}

func newCheckBase(line int, key string, mdField *model.Field) checkBase {
	return checkBase{
		line:    line,
		key:     key,
		mdField: mdField,
	}
}

// for some special diff, we need to show message instead of diff
type diffWithMessage struct {
	checkBase
	msg  string
	skip bool
}

func (i diffWithMessage) Fix(line string) (string, error) {
	return line, nil
}

func (i diffWithMessage) String() string {
	return i.msg
}
func (i diffWithMessage) ShouldSkip() bool {
	return i.skip
}

func newDiffWithMessage(msg string, skip bool) Checker {
	return diffWithMessage{
		checkBase: newCheckBase(0, msg, nil),
		msg:       msg,
		skip:      skip,
	}
}
