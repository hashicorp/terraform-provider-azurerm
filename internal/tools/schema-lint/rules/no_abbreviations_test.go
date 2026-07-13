// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package rules

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/providerschema"
)

func TestNoAbbreviations(t *testing.T) {
	cases := []struct {
		name      string
		propName  string
		wantHits  int
		wantInFix string
	}{
		{name: "vm", propName: "vm_size", wantHits: 1, wantInFix: "virtual_machine_size"},
		{name: "sec", propName: "timeout_sec", wantHits: 1, wantInFix: "timeout_seconds"},
		{name: "rg name", propName: "rg_name", wantHits: 1, wantInFix: "resource_group_name"},
		{name: "pct", propName: "sampling_pct", wantHits: 1, wantInFix: "sampling_percentage"},
		{name: "percent", propName: "cpu_percent", wantHits: 1, wantInFix: "cpu_percentage"},
		{name: "already percentage", propName: "sampling_percentage", wantHits: 0},
		{name: "no abbreviation", propName: "virtual_machine_size", wantHits: 0},
		{name: "standard acronym id", propName: "subnet_id", wantHits: 0},
		{name: "standard acronym ip", propName: "ip_address", wantHits: 0},
	}

	rule := noAbbreviations{}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := rule.CheckProperty(PropertyContext{ResourceType: "azurerm_test", Kind: KindResource, Name: tc.propName, Path: tc.propName, Schema: providerschema.SchemaJSON{Type: TypeString}})
			if len(got) != tc.wantHits {
				t.Fatalf("expected %d finding(s), got %d: %+v", tc.wantHits, len(got), got)
			}
			if tc.wantHits == 1 && !strings.Contains(got[0].FixSuggestion, tc.wantInFix) {
				t.Fatalf("expected fix suggestion to contain %q, got %q", tc.wantInFix, got[0].FixSuggestion)
			}
		})
	}
}
