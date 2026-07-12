// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package rules

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/schema-api/providerjson"
)

func scalarList(optional bool, maxItems int) providerjson.SchemaJSON {
	return providerjson.SchemaJSON{Type: TypeList, Optional: optional, MaxItems: maxItems, Elem: providerjson.SchemaJSON{Type: TypeString}}
}

func TestArrayLimits(t *testing.T) {
	cases := []struct {
		name     string
		schema   providerjson.SchemaJSON
		wantHits int
	}{
		{name: "optional scalar list without max", schema: scalarList(true, 0), wantHits: 1},
		{name: "optional scalar list with max", schema: scalarList(true, 5), wantHits: 0},
		{name: "block list without max", schema: providerjson.SchemaJSON{Type: TypeList, Optional: true, Elem: &providerjson.ResourceJSON{Schema: map[string]providerjson.SchemaJSON{"a": {Type: TypeString}}}}, wantHits: 0},
		{name: "missing elem", schema: providerjson.SchemaJSON{Type: TypeList, Optional: true}, wantHits: 0},
		{name: "not a collection", schema: providerjson.SchemaJSON{Type: TypeString, Optional: true}, wantHits: 0},
		{name: "computed-only list", schema: providerjson.SchemaJSON{Type: TypeList, Computed: true, Elem: providerjson.SchemaJSON{Type: TypeString}}, wantHits: 0},
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
