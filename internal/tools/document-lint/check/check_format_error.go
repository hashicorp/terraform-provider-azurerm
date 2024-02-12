// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package check

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/md"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/util"
)

type formatErr struct {
	Origin string
	msg    string
	checkBase
}

func newFormatErr(origin, msg string, checkBase checkBase) *formatErr {
	return &formatErr{
		Origin:    origin,
		msg:       msg,
		checkBase: checkBase,
	}
}

func (f formatErr) String() string {
	base := f.checkBase.Str()
	switch {
	case strings.Contains(f.msg, "block is not defined in the documentation"):
		return fmt.Sprintf("%s %s", base, util.IssueLine(f.msg))
	case strings.Contains(f.msg, "duplicate"):
		return fmt.Sprintf("%s %s", base, util.IssueLine(f.msg))
	case strings.Contains(f.msg, md.IncorrectlyBlockMarked):
		return fmt.Sprintf("%s %s", base, util.IssueLine(f.msg))
	case strings.TrimSpace(f.Origin) == "*":
		return fmt.Sprintf("%s Found a list marker with no field name or content. This should be removed", base)
	case strings.HasPrefix(f.Origin, "* ~>"):
		return fmt.Sprintf("%s a %s block should not start with `*`", base, util.Bold("Note"))
	default:
		return fmt.Sprintf("%s should be formatted as: %s", base,
			util.FormatCode("* `field` - (Required/Optional) Xxx..."),
		)
	}
}

var regIncorrectMark = regexp.MustCompile(`(?: blocks?)? as (?:detailed|defined) (below|above)`)

func (f formatErr) Fix(line string) (result string, err error) {
	// some Note lines with a misleading star mark, try to remove it
	switch {
	case strings.HasPrefix(line, "* ~>"):
		line = strings.TrimPrefix(line, "* ")
	case strings.Contains(f.msg, md.IncorrectlyBlockMarked):
		line = regIncorrectMark.ReplaceAllLiteralString(line, "")
	case strings.TrimSpace(line) == "*":
		line = ""
	}

	return line, nil // no fix for format error
}

var _ Checker = (*formatErr)(nil)
