// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package rules

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/schema-api/providerjson"
)

func TestValidationRequired(t *testing.T) {
	cases := []struct {
		name     string
		schema   providerjson.SchemaJSON
		wantHits int
	}{
		{name: "optional string without validation", schema: providerjson.SchemaJSON{Type: TypeString, Optional: true}, wantHits: 1},
		{name: "required int without validation", schema: providerjson.SchemaJSON{Type: TypeInt, Required: true}, wantHits: 1},
		{name: "optional float without validation", schema: providerjson.SchemaJSON{Type: TypeFloat, Optional: true}, wantHits: 1},
		{name: "string with validation", schema: providerjson.SchemaJSON{Type: TypeString, Optional: true, HasValidateFunc: true}, wantHits: 0},
		{name: "computed-only string", schema: providerjson.SchemaJSON{Type: TypeString, Computed: true}, wantHits: 0},
		{name: "bool needs no validation", schema: providerjson.SchemaJSON{Type: TypeBool, Optional: true}, wantHits: 0},
		{name: "list is not a scalar", schema: providerjson.SchemaJSON{Type: TypeList, Optional: true}, wantHits: 0},
	}

	rule := validationRequired{}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := rule.CheckProperty(PropertyContext{ResourceType: "azurerm_test", Kind: KindResource, Name: "field", Path: "field", Schema: tc.schema})
			if len(got) != tc.wantHits {
				t.Fatalf("expected %d finding(s), got %d: %+v", tc.wantHits, len(got), got)
			}
		})
	}
}
