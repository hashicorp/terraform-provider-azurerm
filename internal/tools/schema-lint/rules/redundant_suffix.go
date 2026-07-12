// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package rules

import (
	"fmt"
	"strings"
)

var (
	_ PropertyRule = redundantSuffix{}
	_ Fixer        = redundantSuffix{}
)

// redundantSuffix flags names ending in a redundant grouping word that adds no
// informational value (e.g. firewall_properties -> firewall, site_config ->
// site, os_profile -> os).
type redundantSuffix struct{}

// redundantSuffixes are the grouping-word suffixes that should be dropped.
var redundantSuffixes = []string{"_properties", "_config", "_profile"}

func (redundantSuffix) ID() string   { return "SL016" }
func (redundantSuffix) Name() string { return "redundant-suffix" }

func (redundantSuffix) Description() string {
	return `Names should not include a redundant grouping-word suffix (e.g. "_properties", "_config", "_profile").`
}

func (redundantSuffix) DefaultSeverity() Severity { return SeverityWarning }

func (redundantSuffix) FixHint() string { return "drop the redundant suffix" }

func (r redundantSuffix) CheckProperty(ctx PropertyContext) []Finding {
	for _, suffix := range redundantSuffixes {
		if !strings.HasSuffix(ctx.Name, suffix) {
			continue
		}
		preferred := strings.TrimSuffix(ctx.Name, suffix)
		if preferred == "" {
			return nil
		}

		return []Finding{propertyFindingFix(r, ctx,
			fmt.Sprintf("property %q has a redundant %q suffix (%q)", ctx.Path, suffix, preferred),
			fmt.Sprintf("rename %q to %q", ctx.Name, preferred),
		)}
	}

	return nil
}
