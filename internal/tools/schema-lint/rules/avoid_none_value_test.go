// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package rules

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/providerschema"
)

func TestAvoidNoneValue(t *testing.T) {
	// enumAccept mimics a finite-set validator (like StringInSlice) that accepts
	// exactly the given values.
	enumAccept := func(allowed ...string) func(string) bool {
		set := make(map[string]bool, len(allowed))
		for _, a := range allowed {
			set[a] = true
		}
		return func(v string) bool { return set[v] }
	}
	// acceptAll mimics a free-form validator (like StringIsNotEmpty).
	acceptAll := func(string) bool { return true }

	cases := []struct {
		name     string
		schema   providerschema.SchemaJSON
		wantHits int
	}{
		{name: "optional enum accepting None", schema: providerschema.SchemaJSON{Type: TypeString, Optional: true, AcceptsValue: enumAccept("None", "Auto")}, wantHits: 1},
		{name: "required enum accepting Off and Default", schema: providerschema.SchemaJSON{Type: TypeString, Required: true, AcceptsValue: enumAccept("Off", "On", "Default")}, wantHits: 1},
		{name: "enum without special values", schema: providerschema.SchemaJSON{Type: TypeString, Optional: true, AcceptsValue: enumAccept("Standard", "Premium")}, wantHits: 0},
		{name: "free-form validator is not an enum", schema: providerschema.SchemaJSON{Type: TypeString, Optional: true, AcceptsValue: acceptAll}, wantHits: 0},
		{name: "no validator", schema: providerschema.SchemaJSON{Type: TypeString, Optional: true}, wantHits: 0},
		{name: "computed-only accepting None", schema: providerschema.SchemaJSON{Type: TypeString, Computed: true, AcceptsValue: enumAccept("None")}, wantHits: 0},
	}

	rule := avoidNoneValue{}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := rule.CheckProperty(PropertyContext{ResourceType: "azurerm_test", Kind: KindResource, Name: "mode", Path: "mode", Schema: tc.schema})
			if len(got) != tc.wantHits {
				t.Fatalf("expected %d finding(s), got %d: %+v", tc.wantHits, len(got), got)
			}
		})
	}
}
