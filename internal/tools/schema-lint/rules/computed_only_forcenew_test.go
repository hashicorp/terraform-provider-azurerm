// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package rules

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/providerschema"
)

func TestComputedOnlyForceNew(t *testing.T) {
	cases := []struct {
		name     string
		schema   providerschema.SchemaJSON
		wantHits int
	}{
		{name: "computed only without forcenew", schema: providerschema.SchemaJSON{Type: TypeString, Computed: true}, wantHits: 0},
		{name: "optional with forcenew", schema: providerschema.SchemaJSON{Type: TypeString, Optional: true, ForceNew: true}, wantHits: 0},
		{name: "optional computed with forcenew", schema: providerschema.SchemaJSON{Type: TypeString, Optional: true, Computed: true, ForceNew: true}, wantHits: 0},
		{name: "computed only with forcenew", schema: providerschema.SchemaJSON{Type: TypeString, Computed: true, ForceNew: true}, wantHits: 1},
	}

	rule := computedOnlyForceNew{}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := rule.CheckProperty(PropertyContext{ResourceType: "azurerm_test", Kind: KindResource, Name: "name", Path: "name", Schema: tc.schema})
			if len(got) != tc.wantHits {
				t.Fatalf("expected %d finding(s), got %d: %+v", tc.wantHits, len(got), got)
			}
		})
	}
}
