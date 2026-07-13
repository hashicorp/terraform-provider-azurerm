// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package rules

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/providerschema"
)

func TestRedundantSuffix(t *testing.T) {
	cases := []struct {
		name      string
		propName  string
		wantHits  int
		wantInFix string
	}{
		{name: "properties suffix", propName: "firewall_properties", wantHits: 1, wantInFix: `to "firewall"`},
		{name: "config suffix", propName: "site_config", wantHits: 1, wantInFix: `to "site"`},
		{name: "profile suffix", propName: "os_profile", wantHits: 1, wantInFix: `to "os"`},
		{name: "configuration is not flagged", propName: "network_configuration", wantHits: 0},
		{name: "no suffix", propName: "firewall", wantHits: 0},
	}

	rule := redundantSuffix{}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := rule.CheckProperty(PropertyContext{ResourceType: "azurerm_test", Kind: KindResource, Name: tc.propName, Path: tc.propName, Schema: providerschema.SchemaJSON{Type: TypeList}})
			if len(got) != tc.wantHits {
				t.Fatalf("expected %d finding(s), got %d: %+v", tc.wantHits, len(got), got)
			}
			if tc.wantHits == 1 && !strings.Contains(got[0].FixSuggestion, tc.wantInFix) {
				t.Fatalf("expected fix suggestion to contain %q, got %q", tc.wantInFix, got[0].FixSuggestion)
			}
		})
	}
}
