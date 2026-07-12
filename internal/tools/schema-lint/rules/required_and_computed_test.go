// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package rules

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/schema-api/providerjson"
)

func TestRequiredAndComputed(t *testing.T) {
	cases := []struct {
		name     string
		schema   providerjson.SchemaJSON
		wantHits int
	}{
		{name: "required only", schema: providerjson.SchemaJSON{Type: TypeString, Required: true}, wantHits: 0},
		{name: "computed only", schema: providerjson.SchemaJSON{Type: TypeString, Computed: true}, wantHits: 0},
		{name: "optional and computed", schema: providerjson.SchemaJSON{Type: TypeString, Optional: true, Computed: true}, wantHits: 0},
		{name: "required and computed", schema: providerjson.SchemaJSON{Type: TypeString, Required: true, Computed: true}, wantHits: 1},
	}

	rule := requiredAndComputed{}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := rule.CheckProperty(PropertyContext{ResourceType: "azurerm_test", Kind: KindResource, Name: "name", Path: "name", Schema: tc.schema})
			if len(got) != tc.wantHits {
				t.Fatalf("expected %d finding(s), got %d: %+v", tc.wantHits, len(got), got)
			}
		})
	}
}
