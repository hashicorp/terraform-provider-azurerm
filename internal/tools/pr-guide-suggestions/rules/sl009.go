// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package rules

import (
	"fmt"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/pr-guide-suggestions/schematree"
)

var sl009 = &Rule{
	ID:       "SL009",
	Name:     "unit-in-naming",
	Severity: Warning,
	Fixable:  true,
	Check:    checkSL009,
}

var dataSizeSuffixes = []string{
	"byte", "kilobyte", "megabyte", "gigabyte", "terabyte", "petabyte",
	"kib", "mib", "gib", "tib", "pib",
	"b", "kb", "mb", "gb", "tb", "pb",
}

var timeDurationSuffixes = []string{
	"nanoseconds", "microseconds", "milliseconds",
	"seconds", "minutes", "hours", "days", "weeks", "months", "years",
	"ns", "ms", "secs", "mins", "hrs",
}

var suffixesToCheckForBareUnit = slices.Concat(
	dataSizeSuffixes,
	timeDurationSuffixes,
	[]string{"kbps", "mbps", "gbps"},
)

func checkSL009(res *schematree.Result, report ReportFunc) {
	for _, n := range res.All {
		avoidBareUnitSuffix(n, report)
		avoidPluralDataSizeSuffix(n, report)
		avoidSingularTimeSuffix(n, report)
	}
}

func avoidBareUnitSuffix(n *schematree.Node, report ReportFunc) {
	for _, unit := range suffixesToCheckForBareUnit {
		if !strings.HasSuffix(n.Name, "_"+unit) {
			continue
		}
		if strings.HasSuffix(n.Name, "_in_"+unit) {
			break
		}
		prefix := strings.TrimSuffix(n.Name, "_"+unit)
		if prefix == "" {
			break
		}
		preferred := prefix + "_in_" + unit
		report(n,
			fmt.Sprintf("property %q uses a bare unit suffix; prefer the %q naming convention (%q)", n.Path, "_in_"+unit, preferred),
			&Fix{Suggestion: fmt.Sprintf("rename %q to %q", n.Name, preferred), Rename: preferred},
		)
		break
	}
}

func avoidPluralDataSizeSuffix(n *schematree.Node, report ReportFunc) {
	for _, unit := range dataSizeSuffixes {
		pluralSuffix := "_" + unit + "s"
		if !strings.HasSuffix(n.Name, pluralSuffix) {
			continue
		}
		prefix := strings.TrimSuffix(n.Name, pluralSuffix)
		if prefix == "" {
			break
		}
		prefix = strings.TrimSuffix(prefix, "_in")
		preferred := prefix + "_in_" + unit
		report(n,
			fmt.Sprintf("property %q uses a plural data size suffix; prefer %q", n.Path, preferred),
			&Fix{Suggestion: fmt.Sprintf("rename %q to %q", n.Name, preferred), Rename: preferred},
		)
		break
	}
}

func avoidSingularTimeSuffix(n *schematree.Node, report ReportFunc) {
	for _, unit := range timeDurationSuffixes {
		singularSuffix := "_" + strings.TrimSuffix(unit, "s")
		if !strings.HasSuffix(n.Name, singularSuffix) {
			continue
		}
		prefix := strings.TrimSuffix(n.Name, singularSuffix)
		if prefix == "" {
			break
		}
		prefix = strings.TrimSuffix(prefix, "_in")
		preferred := prefix + "_in_" + unit
		report(n,
			fmt.Sprintf("property %q uses a singular time suffix; prefer %q", n.Path, preferred),
			&Fix{Suggestion: fmt.Sprintf("rename %q to %q", n.Name, preferred), Rename: preferred},
		)
		break
	}
}
