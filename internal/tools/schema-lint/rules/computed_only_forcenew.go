// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package rules

import "fmt"

var (
	_ PropertyRule = computedOnlyForceNew{}
	_ Fixer        = computedOnlyForceNew{}
)

// computedOnlyForceNew flags a Computed-only property (neither Optional nor
// Required) that also sets ForceNew. ForceNew has no effect on a value the user
// cannot set, so this usually indicates a mistake.
type computedOnlyForceNew struct{}

func (computedOnlyForceNew) ID() string   { return "SL004" }
func (computedOnlyForceNew) Name() string { return "computed-only-forcenew" }

func (computedOnlyForceNew) Description() string {
	return "A Computed-only property (not Optional or Required) should not set ForceNew, which has no effect."
}

func (computedOnlyForceNew) DefaultSeverity() Severity { return SeverityError }

func (computedOnlyForceNew) FixHint() string { return "remove ForceNew" }

func (r computedOnlyForceNew) CheckProperty(ctx PropertyContext) []Finding {
	if ctx.Schema.Computed && !ctx.Schema.Optional && !ctx.Schema.Required && ctx.Schema.ForceNew {
		return []Finding{propertyFindingFix(r, ctx,
			fmt.Sprintf("property %q is Computed-only but sets ForceNew, which has no effect", ctx.Path),
			fmt.Sprintf("remove ForceNew from the Computed-only property %q", ctx.Path),
		)}
	}

	return nil
}
