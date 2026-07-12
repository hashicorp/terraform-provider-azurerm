// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package rules

import "fmt"

var (
	_ PropertyRule = limitsOnNonCollection{}
	_ Fixer        = limitsOnNonCollection{}
)

// limitsOnNonCollection flags a property that sets MinItems or MaxItems while not
// being a list or set. Those constraints only apply to TypeList and TypeSet, so
// on any other type they have no effect and should be removed.
type limitsOnNonCollection struct{}

func (limitsOnNonCollection) ID() string   { return "SL007" }
func (limitsOnNonCollection) Name() string { return "limits-on-non-collection" }

func (limitsOnNonCollection) Description() string {
	return "MinItems and MaxItems only apply to TypeList and TypeSet; they should not be set on other types."
}

func (limitsOnNonCollection) DefaultSeverity() Severity { return SeverityError }

func (limitsOnNonCollection) FixHint() string { return "remove MinItems/MaxItems" }

func (r limitsOnNonCollection) CheckProperty(ctx PropertyContext) []Finding {
	if IsCollection(ctx.Schema) {
		return nil
	}
	if ctx.Schema.MinItems == 0 && ctx.Schema.MaxItems == 0 {
		return nil
	}

	return []Finding{propertyFindingFix(r, ctx,
		fmt.Sprintf("property %q (%s) sets MinItems/MaxItems, which only apply to TypeList/TypeSet", ctx.Path, ctx.Schema.Type),
		fmt.Sprintf("remove MinItems/MaxItems from %q", ctx.Path),
	)}
}
