// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package check

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/model"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/util"
)

type circularRef struct {
	checkBase
	fieldName string
	md        *model.ResourceDoc
}

func newCircularRef(name string, md *model.ResourceDoc) *circularRef {
	base := newCheckBase(0, "", nil)
	return &circularRef{
		checkBase: base,
		fieldName: name,
		md:        md,
	}
}

// Fix implements Checker.
func (*circularRef) Fix(line string) (result string, err error) {
	return line, nil
}

// Key implements Checker.
// ShouldSkip implements Checker.
func (*circularRef) ShouldSkip() bool {
	return false
}

// String implements Checker.
func (c *circularRef) String() string {
	return fmt.Sprintf("0 document has circular reference in block name: %s", util.Bold(c.fieldName))
}

var _ Checker = (*circularRef)(nil)
