// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package check

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/util"
)

type RequiredMiss int

func (r RequiredMiss) String() string {
	return []string{"ok", "required", "optional", "computed"}[r]
}

const (
	RequiredOK       RequiredMiss = iota // no need to fix requiredness
	ShouldBeRequired                     // code is required, but doc be optional or not specify
	ShouldBeOptional                     // code is optional, but doc be required or not specify
	ShouldBeComputed
)

type requireDiff struct {
	checkBase
	RequiredMiss RequiredMiss
}

func newRequireDiff(checkBase checkBase, requiredMiss RequiredMiss) *requireDiff {
	return &requireDiff{checkBase: checkBase, RequiredMiss: requiredMiss}
}

func (c requireDiff) String() string {
	if c.RequiredMiss == ShouldBeComputed {
		return fmt.Sprintf("%s Fields listed under Attributes Reference should not be marked Required/Optional", c.checkBase.Str())
	}
	str := fmt.Sprintf("%s should be %s", c.checkBase.Str(), util.Blue(c.RequiredMiss))
	return str
}

func (c requireDiff) Fix(line string) (result string, err error) {
	from, to := "(Required)", "(Optional)"

	switch c.RequiredMiss {
	case ShouldBeComputed:
		// remove from both. may cause two space, have to remove one
		var idx, size int
		if idx = strings.Index(line, from); idx > 0 {
			size = len(from)
		} else if idx = strings.Index(line, to); idx > 0 {
			size = len(to)
		}
		if idx > 0 {
			if line[idx-1] == ' ' && line[idx+size] == ' ' {
				idx -= 1
				size += 1
			}
			line = line[:idx] + line[idx+size:]
		}
	case ShouldBeOptional, ShouldBeRequired:
		if c.RequiredMiss == ShouldBeRequired {
			from, to = to, from
		}
		if strings.Contains(line, from) {
			line = strings.Replace(line, from, to, 1)
		} else {
			// add after the first -
			if idx := strings.Index(line, " - "); idx > 0 {
				line = line[:idx+3] + to + " " + line[idx+3:]
			} else {
				// no dash add after second `
				idx = strings.Index(line, "`")
				idx += strings.Index(line[idx+1:], "`") + 1
				line = line[:idx+1] + " " + to + line[idx+1:]
			}
		}
	}
	return line, nil
}

var _ Checker = (*requireDiff)(nil)
