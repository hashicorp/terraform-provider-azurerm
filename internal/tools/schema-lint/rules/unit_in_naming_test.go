// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package rules

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/providerschema"
)

func TestUnitInNaming(t *testing.T) {
	cases := []struct {
		name      string
		propName  string
		wantHits  int
		wantInFix string
	}{
		{name: "size_mb", propName: "size_mb", wantHits: 1, wantInFix: "size_in_mb"},
		{name: "time_sec", propName: "time_sec", wantHits: 1, wantInFix: "time_in_sec"},
		{name: "retention_days", propName: "retention_days", wantHits: 1, wantInFix: "retention_in_days"},
		{name: "disk_size_gb", propName: "disk_size_gb", wantHits: 1, wantInFix: "disk_size_in_gb"},
		{name: "throughput_mbps", propName: "throughput_mbps", wantHits: 1, wantInFix: "throughput_in_mbps"},
		{name: "already size_in_mb", propName: "size_in_mb", wantHits: 0},
		{name: "already duration_in_seconds", propName: "duration_in_seconds", wantHits: 0},
		{name: "no unit suffix", propName: "name", wantHits: 0},
		{name: "max_items is not a unit", propName: "max_items", wantHits: 0},
		{name: "public_holidays underscore boundary", propName: "public_holidays", wantHits: 0},
	}

	rule := unitInNaming{}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := rule.CheckProperty(PropertyContext{ResourceType: "azurerm_test", Kind: KindResource, Name: tc.propName, Path: tc.propName, Schema: providerschema.SchemaJSON{Type: TypeInt}})
			if len(got) != tc.wantHits {
				t.Fatalf("expected %d finding(s), got %d: %+v", tc.wantHits, len(got), got)
			}
			if tc.wantHits == 1 && !strings.Contains(got[0].FixSuggestion, tc.wantInFix) {
				t.Fatalf("expected fix suggestion to contain %q, got %q", tc.wantInFix, got[0].FixSuggestion)
			}
		})
	}
}
