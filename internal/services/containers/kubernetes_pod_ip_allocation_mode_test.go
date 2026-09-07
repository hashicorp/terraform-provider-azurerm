// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package containers_test

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-10-01/managedclusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func TestPodIPAllocationModeSchema(t *testing.T) {
	schemas := map[string]*pluginsdk.Schema{
		"default_node_pool": containers.SchemaDefaultNodePool().Elem.(*pluginsdk.Resource).Schema["pod_ip_allocation_mode"],
		"node_pool":         containers.Registration{}.SupportedResources()["azurerm_kubernetes_cluster_node_pool"].Schema["pod_ip_allocation_mode"],
	}
	for name, schema := range schemas {
		t.Run(name, func(t *testing.T) {
			if !schema.Optional || !schema.Computed || !schema.ForceNew || schema.Default != nil || len(schema.RequiredWith) != 0 {
				t.Fatal("pod_ip_allocation_mode must preserve an omitted API-computed value without a Terraform default")
			}
			for _, mode := range []string{"DynamicIndividual", "StaticBlock"} {
				if _, errors := schema.ValidateFunc(mode, "pod_ip_allocation_mode"); len(errors) != 0 {
					t.Fatalf("expected %s to be accepted: %v", mode, errors)
				}
			}
		})
	}
}

func TestPodIPAllocationModeDefaultNodePoolConversion(t *testing.T) {
	for _, mode := range []string{"", "DynamicIndividual", "StaticBlock"} {
		t.Run(mode, func(t *testing.T) {
			profile := managedclusters.ManagedClusterAgentPoolProfile{Name: "default"}
			if mode != "" {
				profile.PodIPAllocationMode = pointer.ToEnum[managedclusters.PodIPAllocationMode](mode)
			}
			result := containers.ConvertDefaultNodePoolToAgentPool(&[]managedclusters.ManagedClusterAgentPoolProfile{profile})
			if mode == "" {
				if result.Properties.PodIPAllocationMode != nil {
					t.Fatal("omitted pod IP allocation mode must remain omitted")
				}
			} else if got := pointer.FromEnum(result.Properties.PodIPAllocationMode); got != mode {
				t.Fatalf("pod IP allocation mode lost during node pool conversion: got %q, want %q", got, mode)
			}
		})
	}
}
