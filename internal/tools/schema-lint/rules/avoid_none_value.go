// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package rules

import (
	"fmt"
	"strings"
)

var _ PropertyRule = avoidNoneValue{}

type avoidNoneValue struct{}

// noneValues are paired with nonsense decoys that share its length and
// letter-case shape: a finite-set (enum) validator accepts the special value but
// rejects every decoy, whereas a format or length validator accepts the
// same-shaped decoys too. This lets the rule flag only genuine enum validators.
var noneValues = []struct {
	value  string
	decoys []string
}{
	{value: "None", decoys: []string{"Nabe", "Vise", "Lome"}},
	{value: "Off", decoys: []string{"Eft", "Uzz", "Alk"}},
	{value: "Default", decoys: []string{"Marlend", "Bequlor", "Savnite"}},
	{value: "Disabled", decoys: []string{"Ronketal", "Verbunth", "Malforte"}},
}

func (avoidNoneValue) ID() string   { return "SL004" }
func (avoidNoneValue) Name() string { return "avoid-none-value" }

func (avoidNoneValue) Description() string {
	return "A user-settable enum property should not accept None/Off/Default/Disabled as valid values; omit them and normalise in Create/Read."
}

func (avoidNoneValue) DefaultSeverity() Severity { return SeverityWarning }

func (r avoidNoneValue) CheckProperty(ctx PropertyContext) []Finding {
	// Computed-only properties should return the value from Azure (including
	// None), so only user-settable properties are checked.
	if !ctx.Schema.Optional && !ctx.Schema.Required {
		return nil
	}

	accepted := acceptedSpecialValues(ctx.Schema.AcceptsValue)
	if len(accepted) == 0 {
		return nil
	}

	return []Finding{propertyFinding(r, ctx, fmt.Sprintf("property %q accepts %s via validation; omit these values (use Terraform null) and normalise in Create/Read", ctx.Path, strings.Join(accepted, ", ")))}
}

// acceptedSpecialValues returns the special sentinel values that a finite-set
// (enum) validator accepts, using the property's accept predicate. It returns
// nil when accept is nil (no string validator) or when the validator is not an
// enum (it also accepts the same-shaped decoys).
func acceptedSpecialValues(accept func(string) bool) []string {
	if accept == nil {
		return nil
	}

	var out []string
	for _, s := range noneValues {
		if !accept(s.value) {
			continue
		}
		if acceptsAny(accept, s.decoys) {
			continue
		}
		out = append(out, s.value)
	}

	return out
}

// acceptsAny reports whether accept returns true for any of the given values.
func acceptsAny(accept func(string) bool, values []string) bool {
	for _, v := range values {
		if accept(v) {
			return true
		}
	}

	return false
}
