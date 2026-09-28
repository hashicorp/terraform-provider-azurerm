// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package containers_test

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
)

func TestKubernetesClusterNodePoolWorkloadRuntimeConfigurations(t *testing.T) {
	data := acceptance.TestData{
		RandomInteger: 12345,
		Locations:     acceptance.Regions{Primary: "eastus"},
	}
	r := KubernetesClusterNodePoolResource{}
	tests := []struct {
		name   string
		config string
	}{
		{
			name:   "OCIContainer",
			config: r.workloadRuntime(data, "OCIContainer"),
		},
		{
			name:   "KataVmIsolation",
			config: r.workloadRuntimeKataVmIsolation(data),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config := []byte(test.config)
			file, diags := hclsyntax.ParseConfig(config, "config.tf", hcl.InitialPos)
			if diags.HasErrors() {
				t.Fatalf("parsing fixture: %s", diags)
			}
			resources, _, diags := file.Body.PartialContent(&hcl.BodySchema{
				Blocks: []hcl.BlockHeaderSchema{{Type: "resource", LabelNames: []string{"type", "name"}}},
			})
			if diags.HasErrors() {
				t.Fatalf("reading resources: %s", diags)
			}
			pools := 0
			for _, block := range resources.Blocks {
				if block.Labels[0] != "azurerm_kubernetes_cluster_node_pool" || block.Labels[1] != "test" {
					continue
				}
				pools++
				content, _, diags := block.Body.PartialContent(&hcl.BodySchema{
					Attributes: []hcl.AttributeSchema{{Name: "node_count", Required: true}},
				})
				if diags.HasErrors() {
					t.Fatalf("reading additional pool node_count: %s", diags)
				}
				if got := string(content.Attributes["node_count"].Expr.Range().SliceBytes(config)); got != "1" {
					t.Errorf("expected additional pool node_count to be literal 1, got %q", got)
				}
			}
			if pools != 1 {
				t.Errorf("expected one additional test pool, got %d", pools)
			}
		})
	}
}
