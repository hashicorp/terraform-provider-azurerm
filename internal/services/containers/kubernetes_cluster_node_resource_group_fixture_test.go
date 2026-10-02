// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package containers_test

import (
	"bytes"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/zclconf/go-cty/cty"
)

func TestKubernetesClusterNodeResourceGroupRemovalFixture(t *testing.T) {
	data := acceptance.TestData{
		RandomInteger: 12345,
		Locations:     acceptance.Regions{Primary: "eastus"},
	}
	r := KubernetesClusterResource{}
	for _, test := range []struct {
		name  string
		level string
	}{
		{name: "configured", level: "Unrestricted"},
		{name: "omitted"},
	} {
		t.Run(test.name, func(t *testing.T) {
			original := []byte(r.nodeResourceGroupRestrictionLevel(data, test.level))
			removal := []byte(r.nodeResourceGroupRestrictionLevelRemoval(data, test.level))
			before, diagnostics := hclsyntax.ParseConfig(original, "original.tf", hcl.InitialPos)
			if diagnostics.HasErrors() {
				t.Fatal(diagnostics)
			}
			after, diagnostics := hclsyntax.ParseConfig(removal, "removal.tf", hcl.InitialPos)
			if diagnostics.HasErrors() {
				t.Fatal(diagnostics)
			}
			beforeBlocks := before.Body.(*hclsyntax.Body).Blocks
			afterBlocks := after.Body.(*hclsyntax.Body).Blocks
			if len(beforeBlocks) != 3 || len(afterBlocks) != 3 {
				t.Fatalf("expected one provider and two resources, got %d and %d blocks", len(beforeBlocks), len(afterBlocks))
			}
			for i, block := range afterBlocks {
				if i != 0 {
					if block.Type != "resource" || !bytes.Equal(beforeBlocks[i].Range().SliceBytes(original), block.Range().SliceBytes(removal)) {
						t.Fatalf("resource block %d changed", i)
					}
					continue
				}
				if block.Type != "provider" || len(block.Labels) != 1 || block.Labels[0] != "azurerm" || len(block.Body.Attributes) != 0 || len(block.Body.Blocks) != 1 {
					t.Fatal("expected only the azurerm provider features block")
				}
				features := block.Body.Blocks[0]
				originalFeatures := beforeBlocks[i].Body.Blocks
				if len(originalFeatures) != 1 || len(originalFeatures[0].Body.Attributes) != 0 || len(originalFeatures[0].Body.Blocks) != 0 {
					t.Fatal("the ordinary restriction fixture must retain default features")
				}
				if features.Type != "features" || len(features.Body.Attributes) != 1 || len(features.Body.Blocks) != 0 {
					t.Fatal("expected only the early-ID persistence feature; no RG deletion safety override")
				}
				attribute, ok := features.Body.Attributes["persist_id_on_create_before_polling_for_completion"]
				if !ok {
					t.Fatal("missing early-ID persistence feature")
				}
				value, diagnostics := attribute.Expr.Value(nil)
				if diagnostics.HasErrors() || !value.RawEquals(cty.True) {
					t.Fatal("early-ID persistence must be explicitly true")
				}
			}
		})
	}
}
