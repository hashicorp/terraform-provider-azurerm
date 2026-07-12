// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package rules

import (
	"fmt"
	"strings"
)

var (
	_ PropertyRule = noAbbreviations{}
	_ Fixer        = noAbbreviations{}
)

// noAbbreviations flags property names that use an abbreviation instead of the
// full word(s), per the provider naming conventions (e.g. "vm" -> "virtual_machine",
// "sec" -> "seconds", "pct" -> "percentage").
type noAbbreviations struct{}

// abbreviations maps a name segment to its preferred expansion. Standard
// acronyms the provider keeps as-is (id, ip, os, sku, dns, url, cidr, ...) are
// intentionally excluded.
var abbreviations = map[string]string{
	"vm":            "virtual_machine",
	"vmss":          "virtual_machine_scale_set",
	"rg":            "resource_group",
	"vnet":          "virtual_network",
	"nsg":           "network_security_group",
	"nic":           "network_interface",
	"sec":           "seconds",
	"addr":          "address",
	"msg":           "message",
	"num":           "number",
	"cfg":           "configuration",
	"db":            "database",
	"dbs":           "databases",
	"mgmt":          "management",
	"min":           "minutes/minimum", // ambiguous, but "minutes" is more common
	"hr":            "hours",
	"max":           "maximum",
	"svc":           "service",
	"src":           "source",
	"dest":          "destination",
	"pwd":           "password",
	"passwd":        "password",
	"pct":           "percentage",
	"percent":       "percentage",
	"email_address": "email",
}

func (noAbbreviations) ID() string   { return "SL014" }
func (noAbbreviations) Name() string { return "no-abbreviations" }

func (noAbbreviations) Description() string {
	return "Property names should use full words rather than abbreviations (e.g. virtual_machine instead of vm)."
}

func (noAbbreviations) DefaultSeverity() Severity { return SeverityWarning }

func (noAbbreviations) FixHint() string { return "expand the abbreviation" }

func (r noAbbreviations) CheckProperty(ctx PropertyContext) []Finding {
	segments := strings.Split(ctx.Name, "_")
	found := false
	for i, seg := range segments {
		if full, ok := abbreviations[seg]; ok {
			segments[i] = full
			found = true
		}
	}
	if !found {
		return nil
	}

	preferred := strings.Join(segments, "_")

	return []Finding{propertyFindingFix(r, ctx,
		fmt.Sprintf("property %q uses an abbreviation; prefer full words (%q)", ctx.Path, preferred),
		fmt.Sprintf("rename %q to %q", ctx.Name, preferred),
	)}
}
