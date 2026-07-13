// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package rules

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/providerschema"
)

func TestIDReferenceValidation(t *testing.T) {
	acceptAll := func(string) bool { return true }                                         // StringIsNotEmpty-like
	uuidLike := func(v string) bool { return v == "550e8400-e29b-41d4-a716-446655440000" } // rejects garbage

	cases := []struct {
		name     string
		propName string
		schema   providerschema.SchemaJSON
		wantHits int
	}{
		{name: "id without validation", propName: "subnet_id", schema: providerschema.SchemaJSON{Type: TypeString, Required: true}, wantHits: 1},
		{name: "id with weak validation", propName: "subnet_id", schema: providerschema.SchemaJSON{Type: TypeString, Required: true, AcceptsValue: acceptAll}, wantHits: 1},
		{name: "id with strong validation", propName: "subnet_id", schema: providerschema.SchemaJSON{Type: TypeString, Required: true, AcceptsValue: uuidLike}, wantHits: 0},
		{name: "computed-only id", propName: "subnet_id", schema: providerschema.SchemaJSON{Type: TypeString, Computed: true, AcceptsValue: acceptAll}, wantHits: 0},
		{name: "plural ids not matched", propName: "subnet_ids", schema: providerschema.SchemaJSON{Type: TypeString, Required: true}, wantHits: 0},
		{name: "non-id property", propName: "location", schema: providerschema.SchemaJSON{Type: TypeString, Required: true}, wantHits: 0},
		{name: "resource id itself", propName: "id", schema: providerschema.SchemaJSON{Type: TypeString, Computed: true}, wantHits: 0},
	}

	rule := idReferenceValidation{}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := rule.CheckProperty(PropertyContext{ResourceType: "azurerm_test", Kind: KindResource, Name: tc.propName, Path: tc.propName, Schema: tc.schema})
			if len(got) != tc.wantHits {
				t.Fatalf("expected %d finding(s), got %d: %+v", tc.wantHits, len(got), got)
			}
		})
	}
}
