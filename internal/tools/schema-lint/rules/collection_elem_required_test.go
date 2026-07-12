// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package rules

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/schema-api/providerjson"
)

func TestCollectionElemRequired(t *testing.T) {
	cases := []struct {
		name     string
		schema   providerjson.SchemaJSON
		wantHits int
	}{
		{name: "string is not a collection", schema: providerjson.SchemaJSON{Type: TypeString}, wantHits: 0},
		{name: "list with elem", schema: providerjson.SchemaJSON{Type: TypeList, Elem: providerjson.SchemaJSON{Type: TypeString}}, wantHits: 0},
		{name: "set with block elem", schema: providerjson.SchemaJSON{Type: TypeSet, Elem: &providerjson.ResourceJSON{}}, wantHits: 0},
		{name: "list without elem", schema: providerjson.SchemaJSON{Type: TypeList}, wantHits: 1},
		{name: "set without elem", schema: providerjson.SchemaJSON{Type: TypeSet}, wantHits: 1},
	}

	rule := collectionElemRequired{}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := rule.CheckProperty(PropertyContext{ResourceType: "azurerm_test", Kind: KindResource, Name: "items", Path: "items", Schema: tc.schema})
			if len(got) != tc.wantHits {
				t.Fatalf("expected %d finding(s), got %d: %+v", tc.wantHits, len(got), got)
			}
		})
	}
}
