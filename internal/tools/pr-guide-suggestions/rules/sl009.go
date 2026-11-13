// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package rules

import (
	"fmt"
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

func checkSL009(res *schematree.Result, report ReportFunc) {
	for _, n := range res.All {
		avoidBareUnitSuffix(n, report)
		avoidPluralDataSizeSuffix(n, report)
		avoidSingularTimeSuffix(n, report)
	}
}

var suffixesToCheckForBareUnit = []string{
	"bytes",
	"kib", "mib", "gib", "tib",
	"kb", "mb", "gb", "tb", "pb",
	"kbps", "mbps", "gbps",
	"ms", "seconds", "secs", "minutes", "hours", "hrs", "days",
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

var dataSuffixesDenyList = []string{"kbs", "mbs", "gbs", "tbs", "pbs", "kibs", "mibs", "gibs", "tibs", "pibs"}

func avoidPluralDataSizeSuffix(n *schematree.Node, report ReportFunc) {
	for _, unit := range dataSuffixesDenyList {
		if !strings.HasSuffix(n.Name, "_"+unit) {
			continue
		}
		prefix := strings.TrimSuffix(n.Name, "_"+unit)
		if prefix == "" {
			break
		}
		prefix = strings.TrimSuffix(prefix, "_in")
		preferred := prefix + "_in_" + strings.TrimSuffix(unit, "s")
		report(n,
			fmt.Sprintf("property %q uses a plural data size suffix; prefer %q", n.Path, preferred),
			&Fix{Suggestion: fmt.Sprintf("rename %q to %q", n.Name, preferred), Rename: preferred},
		)
		break
	}
}

var timeSuffixesDenyList = []string{
	"nanosecond", "microsecond", "millisecond",
	"second", "sec", "minute", "min", "hour", "hr", "day", "week", "month", "year",
}

func avoidSingularTimeSuffix(n *schematree.Node, report ReportFunc) {
	for _, unit := range timeSuffixesDenyList {
		if !strings.HasSuffix(n.Name, "_"+unit) {
			continue
		}
		prefix := strings.TrimSuffix(n.Name, "_"+unit)
		if prefix == "" {
			break
		}
		prefix = strings.TrimSuffix(prefix, "_in")
		preferred := prefix + "_in_" + unit + "s"
		report(n,
			fmt.Sprintf("property %q uses a singular time suffix; prefer %q", n.Path, preferred),
			&Fix{Suggestion: fmt.Sprintf("rename %q to %q", n.Name, preferred), Rename: preferred},
		)
		break
	}
}
