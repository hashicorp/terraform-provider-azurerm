// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package rules

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/schema-api/providerjson"
)

func TestRedundantIsPrefix(t *testing.T) {
	cases := []struct {
		name      string
		propName  string
		typ       string
		wantHits  int
		wantInFix string
	}{
		{name: "bool with is_ prefix", propName: "is_storage_enabled", typ: TypeBool, wantHits: 1, wantInFix: "storage_enabled"},
		{name: "bool without is_ prefix", propName: "storage_enabled", typ: TypeBool, wantHits: 0},
		{name: "non-bool with is_ prefix", propName: "is_something", typ: TypeString, wantHits: 0},
	}

	rule := redundantIsPrefix{}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := rule.CheckProperty(PropertyContext{ResourceType: "azurerm_test", Kind: KindResource, Name: tc.propName, Path: tc.propName, Schema: providerjson.SchemaJSON{Type: tc.typ}})
			if len(got) != tc.wantHits {
				t.Fatalf("expected %d finding(s), got %d: %+v", tc.wantHits, len(got), got)
			}
			if tc.wantHits == 1 && !strings.Contains(got[0].FixSuggestion, tc.wantInFix) {
				t.Fatalf("expected fix suggestion to contain %q, got %q", tc.wantInFix, got[0].FixSuggestion)
			}
		})
	}
}
