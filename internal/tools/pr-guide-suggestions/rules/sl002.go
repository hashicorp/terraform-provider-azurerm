// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package rules

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/pr-guide-suggestions/schematree"
)

var sl002 = &Rule{
	ID:       "SL002",
	Name:     "single-property-block",
	Severity: Warning,
	Fixable:  true,
	Check:    checkSL002,
}

func checkSL002(res *schematree.Result, report ReportFunc) {
	for _, n := range res.All {
		if n.Schema == nil || n.Schema.Int(fieldMaxItems) != 1 {
			continue
		}
		if len(n.Children) != 1 {
			continue
		}

		child := onlyChild(n.Children)
		msg := fmt.Sprintf("block %q has a single nested property %q (MaxItems 1); consider flattening it", n.Path, child.Name)

		fix := fmt.Sprintf("flatten %q into a single top-level %q property", n.Path, n.Name+"_"+child.Name)
		if child.Name == "enabled" && child.Schema != nil && child.Schema.ValueType() == typeBool {
			fix = fmt.Sprintf("replace the block %q with a single top-level boolean %q", n.Path, n.Name+"_enabled")
		}

		report(n, msg, &Fix{Suggestion: fix})
	}
}

func onlyChild(children map[string]*schematree.Node) *schematree.Node {
	for _, c := range children {
		return c
	}
	return nil
}
