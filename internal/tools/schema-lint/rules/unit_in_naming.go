// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package rules

import (
	"fmt"
	"strings"
)

var (
	_ PropertyRule = unitInNaming{}
	_ Fixer        = unitInNaming{}
)

// unitInNaming flags a property whose name ends in a bare unit of measure (for
// example size_mb or time_sec) instead of the provider's "_in_<unit>" naming
// convention (size_in_mb, time_in_sec).
type unitInNaming struct{}

// unitSuffixes are the units of measure that should be written as "_in_<unit>".
// Singular/ambiguous forms (e.g. "day", "hour", "min") are intentionally
// omitted to avoid false positives.
var unitSuffixes = []string{
	"bytes",
	"kib", "mib", "gib", "tib",
	"kb", "mb", "gb", "tb", "pb",
	"kbps", "mbps", "gbps",
	"ms", "seconds", "secs", "sec", "minutes", "hours", "hrs", "days",
}

func (unitInNaming) ID() string   { return "SL013" }
func (unitInNaming) Name() string { return "unit-in-naming" }

func (unitInNaming) Description() string {
	return `A property that ends in a unit of measure should use the "_in_<unit>" convention (e.g. size_in_mb, duration_in_seconds).`
}

func (unitInNaming) DefaultSeverity() Severity { return SeverityWarning }

func (unitInNaming) FixHint() string { return `use the "_in_<unit>" form` }

func (r unitInNaming) CheckProperty(ctx PropertyContext) []Finding {
	for _, unit := range unitSuffixes {
		if !strings.HasSuffix(ctx.Name, "_"+unit) {
			continue
		}
		// Already in the preferred form.
		if strings.HasSuffix(ctx.Name, "_in_"+unit) {
			return nil
		}

		prefix := strings.TrimSuffix(ctx.Name, "_"+unit)
		if prefix == "" {
			return nil
		}
		preferred := prefix + "_in_" + unit

		return []Finding{propertyFindingFix(r, ctx,
			fmt.Sprintf("property %q uses a bare unit suffix; prefer the %q naming convention (%q)", ctx.Path, "_in_"+unit, preferred),
			fmt.Sprintf("rename %q to %q", ctx.Name, preferred),
		)}
	}

	return nil
}
