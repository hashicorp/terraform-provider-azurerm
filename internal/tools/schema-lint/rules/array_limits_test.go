// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package rules

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/providerschema"
)

func scalarList(optional bool, maxItems int) providerschema.SchemaJSON {
	return providerschema.SchemaJSON{Type: TypeList, Optional: optional, MaxItems: maxItems, Elem: providerschema.SchemaJSON{Type: TypeString}}
}

func TestArrayLimits(t *testing.T) {
	cases := []struct {
		name     string
		schema   providerschema.SchemaJSON
		wantHits int
	}{
		{name: "optional scalar list without max", schema: scalarList(true, 0), wantHits: 1},
		{name: "optional scalar list with max", schema: scalarList(true, 5), wantHits: 0},
		{name: "block list without max", schema: providerschema.SchemaJSON{Type: TypeList, Optional: true, Elem: &providerschema.ResourceJSON{Schema: map[string]providerschema.SchemaJSON{"a": {Type: TypeString}}}}, wantHits: 0},
		{name: "missing elem", schema: providerschema.SchemaJSON{Type: TypeList, Optional: true}, wantHits: 0},
		{name: "not a collection", schema: providerschema.SchemaJSON{Type: TypeString, Optional: true}, wantHits: 0},
		{name: "computed-only list", schema: providerschema.SchemaJSON{Type: TypeList, Computed: true, Elem: providerschema.SchemaJSON{Type: TypeString}}, wantHits: 0},
	}

	rule := arrayLimits{}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := rule.CheckProperty(PropertyContext{ResourceType: "azurerm_test", Kind: KindResource, Name: "items", Path: "items", Schema: tc.schema})
			if len(got) != tc.wantHits {
				t.Fatalf("expected %d finding(s), got %d: %+v", tc.wantHits, len(got), got)
			}
		})
	}
}
