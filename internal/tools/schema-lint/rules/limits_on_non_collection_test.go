// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package rules

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/providerschema"
)

func TestLimitsOnNonCollection(t *testing.T) {
	cases := []struct {
		name     string
		schema   providerschema.SchemaJSON
		wantHits int
	}{
		{name: "string with maxitems", schema: providerschema.SchemaJSON{Type: TypeString, MaxItems: 5}, wantHits: 1},
		{name: "string with minitems", schema: providerschema.SchemaJSON{Type: TypeString, MinItems: 1}, wantHits: 1},
		{name: "map with maxitems", schema: providerschema.SchemaJSON{Type: TypeMap, MaxItems: 3}, wantHits: 1},
		{name: "list with maxitems", schema: providerschema.SchemaJSON{Type: TypeList, MaxItems: 5, Elem: providerschema.SchemaJSON{Type: TypeString}}, wantHits: 0},
		{name: "set with minitems", schema: providerschema.SchemaJSON{Type: TypeSet, MinItems: 1, Elem: providerschema.SchemaJSON{Type: TypeString}}, wantHits: 0},
		{name: "string without limits", schema: providerschema.SchemaJSON{Type: TypeString}, wantHits: 0},
	}

	rule := limitsOnNonCollection{}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := rule.CheckProperty(PropertyContext{ResourceType: "azurerm_test", Kind: KindResource, Name: "field", Path: "field", Schema: tc.schema})
			if len(got) != tc.wantHits {
				t.Fatalf("expected %d finding(s), got %d: %+v", tc.wantHits, len(got), got)
			}
		})
	}
}
