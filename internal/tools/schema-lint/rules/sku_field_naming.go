// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package rules

import (
	"fmt"
	"strings"
)

var _ PropertyRule = skuFieldNaming{}

// skuFieldNaming flags resources that expose multiple top-level sku_* arguments
// (e.g. sku_name and sku_capacity), which should be grouped into a single `sku`
// block. A lone sku_* argument is not flagged.
type skuFieldNaming struct{}

func (skuFieldNaming) ID() string   { return "SL008" }
func (skuFieldNaming) Name() string { return "sku-field-naming" }

func (skuFieldNaming) Description() string {
	return "Multiple top-level sku_* arguments should be grouped into a single `sku` block."
}

func (skuFieldNaming) DefaultSeverity() Severity { return SeverityWarning }

func (r skuFieldNaming) CheckProperty(ctx PropertyContext) []Finding {
	// Only applies to top-level arguments.
	if ctx.Parent != nil {
		return nil
	}
	if !strings.HasPrefix(ctx.Name, "sku_") {
		return nil
	}

	// Only warn when multiple sku_* arguments exist at this level; a single
	// sku_* argument is acceptable. Multiple should be grouped into a `sku`
	// block.
	skuFields := 0
	for name := range ctx.Siblings {
		if strings.HasPrefix(name, "sku_") {
			skuFields++
		}
	}
	if skuFields <= 1 {
		return nil
	}

	return []Finding{propertyFinding(r, ctx, fmt.Sprintf("argument %q uses the sku_* format; with multiple sku_* arguments, prefer a single `sku` block", ctx.Path))}
}
