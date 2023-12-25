// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package check

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/md"
)

type ForceNewType int

func (f ForceNewType) String() string {
	if f <= 2 {
		return []string{"ok", "should be ForceNew", "is NOT ForceNew"}[f]
	}
	return ""
}

const (
	ForceNewDefault ForceNewType = iota
	ShouldBeForceNew
	ShouldBeNotForceNew
)

type forceNewDiff struct {
	checkBase
	ForceNew ForceNewType
}

func newForceNewDiff(checkBase checkBase, forceNew ForceNewType) *forceNewDiff {
	return &forceNewDiff{checkBase: checkBase, ForceNew: forceNew}
}

func (c forceNewDiff) String() string {
	return fmt.Sprintf("%s %s ", c.checkBase.Str(), c.ForceNew)
}

func (c forceNewDiff) Fix(line string) (result string, err error) {
	switch c.ForceNew {
	case ShouldBeForceNew:
		line = strings.TrimRight(line, " ")
		if strings.HasSuffix(line, ",") {
			line = line[:len(line)-1] + "."
		} else if !strings.HasSuffix(line, ".") {
			line += "."
		}
		line += " Changing this forces a new resource to be created."
	case ShouldBeNotForceNew:
		line = md.ForceNewReg.ReplaceAllString(line, "")
	}
	return line, nil
}

var _ Checker = (*forceNewDiff)(nil)
