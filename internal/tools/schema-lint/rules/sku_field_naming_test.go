// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package rules

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/schema-api/providerjson"
)

func TestSkuFieldNaming(t *testing.T) {
	str := providerjson.SchemaJSON{Type: TypeString, Optional: true}
	parent := &providerjson.SchemaJSON{Type: TypeList}

	single := map[string]providerjson.SchemaJSON{"sku_name": str, "location": str}
	multiple := map[string]providerjson.SchemaJSON{"sku_name": str, "sku_capacity": str}

	cases := []struct {
		name       string
		propName   string
		parent     *providerjson.SchemaJSON
		siblings   map[string]providerjson.SchemaJSON
		wantHits   int
		wantSubstr string
	}{
		{name: "single sku field is not flagged", propName: "sku_name", siblings: single, wantHits: 0},
		{name: "multiple sku fields prefer block", propName: "sku_name", siblings: multiple, wantHits: 1, wantSubstr: "`sku` block"},
		{name: "second sku field also prefers block", propName: "sku_capacity", siblings: multiple, wantHits: 1, wantSubstr: "`sku` block"},
		{name: "top-level sku (exact)", propName: "sku", siblings: map[string]providerjson.SchemaJSON{"sku": str}, wantHits: 0},
		{name: "top-level non-sku", propName: "name", siblings: map[string]providerjson.SchemaJSON{"name": str}, wantHits: 0},
		{name: "nested sku_name", propName: "sku_name", parent: parent, siblings: multiple, wantHits: 0},
	}

	rule := skuFieldNaming{}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := rule.CheckProperty(PropertyContext{ResourceType: "azurerm_test", Kind: KindResource, Name: tc.propName, Path: tc.propName, Schema: str, Parent: tc.parent, Siblings: tc.siblings})
			if len(got) != tc.wantHits {
				t.Fatalf("expected %d finding(s), got %d: %+v", tc.wantHits, len(got), got)
			}
			if tc.wantHits == 1 && tc.wantSubstr != "" && !strings.Contains(got[0].Message, tc.wantSubstr) {
				t.Fatalf("expected message to contain %q, got %q", tc.wantSubstr, got[0].Message)
			}
		})
	}
}
