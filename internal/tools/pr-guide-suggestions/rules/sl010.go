// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package rules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/pr-guide-suggestions/schematree"
)

var sl010 = &Rule{
	ID:       "SL010",
	Name:     "no-abbreviations",
	Severity: Warning,
	Fixable:  true,
	Check:    checkSL010,
}

var abbreviations = map[string]string{
	"vm":            "virtual_machine",
	"vmss":          "virtual_machine_scale_set",
	"rg":            "resource_group",
	"vnet":          "virtual_network",
	"nsg":           "network_security_group",
	"nic":           "network_interface",
	"fqdn":          "fully_qualified_domain_name",
	"rt":            "route_table",
	"lb":            "load_balancer",
	"waf":           "web_application_firewall",
	"sec":           "seconds",
	"addr":          "address",
	"msg":           "message",
	"num":           "number",
	"cfg":           "configuration",
	"db":            "database",
	"dbs":           "databases",
	"mgmt":          "management",
	"min":           "minimum",
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

func checkSL010(res *schematree.Result, report ReportFunc) {
	for _, n := range res.All {
		segments := strings.Split(n.Name, "_")
		found := false
		for i, seg := range segments {
			if full, ok := abbreviations[seg]; ok {
				segments[i] = full
				found = true
			}
		}
		if !found {
			continue
		}

		preferred := strings.Join(segments, "_")
		report(n,
			fmt.Sprintf("property %q uses an abbreviation; prefer full words (%q)", n.Path, preferred),
			&Fix{Suggestion: fmt.Sprintf("rename %q to %q", n.Name, preferred), Rename: preferred},
		)
	}
}
