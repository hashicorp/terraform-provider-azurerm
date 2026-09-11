// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package containers_test

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-10-01/managedclusters"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func TestPodIPAllocationModeValidation(t *testing.T) {
	resources := podIPAllocationModeTestResources()
	subnetId := commonids.NewSubnetID("00000000-0000-0000-0000-000000000000", "test", "test", "pods").ID()
	dynamic := string(managedclusters.PodIPAllocationModeDynamicIndividual)
	static := string(managedclusters.PodIPAllocationModeStaticBlock)
	tests := []struct {
		name      string
		config    map[string]interface{}
		wantError bool
	}{
		{name: "omitted", config: map[string]interface{}{}},
		{name: "null_mode", config: map[string]interface{}{"pod_ip_allocation_mode": nil}},
		{name: "subnet_without_mode", config: map[string]interface{}{"pod_subnet_id": subnetId}},
		{name: "subnet_with_null_mode", config: map[string]interface{}{"pod_subnet_id": subnetId, "pod_ip_allocation_mode": nil}},
		{name: "dynamic_with_subnet", config: map[string]interface{}{"pod_ip_allocation_mode": dynamic, "pod_subnet_id": subnetId}},
		{name: "static_with_subnet", config: map[string]interface{}{"pod_ip_allocation_mode": static, "pod_subnet_id": subnetId}},
		{name: "dynamic_without_subnet", config: map[string]interface{}{"pod_ip_allocation_mode": dynamic}, wantError: true},
		{name: "static_without_subnet", config: map[string]interface{}{"pod_ip_allocation_mode": static}, wantError: true},
		{name: "invalid_mode", config: map[string]interface{}{"pod_ip_allocation_mode": "invalid", "pod_subnet_id": subnetId}, wantError: true},
	}
	for name, resource := range resources {
		t.Run(name, func(t *testing.T) {
			for _, test := range tests {
				t.Run(test.name, func(t *testing.T) {
					config := test.config
					if name == "default_node_pool" {
						config = map[string]interface{}{"default_node_pool": []interface{}{config}}
					}
					diagnostics := resource.Validate(terraform.NewResourceConfigRaw(config))
					if diagnostics.HasError() != test.wantError {
						t.Fatalf("validation error = %t, want %t: %v", diagnostics.HasError(), test.wantError, diagnostics)
					}
				})
			}
		})
	}
}

func TestPodIPAllocationModeExistingState(t *testing.T) {
	subnetId := commonids.NewSubnetID("00000000-0000-0000-0000-000000000000", "test", "test", "pods").ID()
	static := string(managedclusters.PodIPAllocationModeStaticBlock)
	dynamic := string(managedclusters.PodIPAllocationModeDynamicIndividual)
	tests := []struct {
		name            string
		stateMode       string
		configuredMode  string
		wantReplacement bool
	}{
		{name: "omitted_with_empty_api_value"},
		{name: "api_value_preserved_when_config_omits", stateMode: static},
		{name: "configured_value_unchanged", stateMode: static, configuredMode: static},
		{name: "configured_change_replaces_pool", stateMode: static, configuredMode: dynamic, wantReplacement: true},
	}

	for name, resource := range podIPAllocationModeTestResources() {
		t.Run(name, func(t *testing.T) {
			for _, test := range tests {
				t.Run(test.name, func(t *testing.T) {
					config := map[string]interface{}{"pod_subnet_id": subnetId}
					if test.configuredMode != "" {
						config["pod_ip_allocation_mode"] = test.configuredMode
					}
					attributes := map[string]string{"id": "test"}
					prefix := ""
					if name == "default_node_pool" {
						config = map[string]interface{}{"default_node_pool": []interface{}{config}}
						attributes["default_node_pool.#"] = "1"
						prefix = "default_node_pool.0."
					}
					attributes[prefix+"pod_subnet_id"] = subnetId
					attributes[prefix+"pod_ip_allocation_mode"] = test.stateMode
					state := &terraform.InstanceState{ID: "test", Attributes: attributes}
					diff, err := resource.SimpleDiff(t.Context(), state, terraform.NewResourceConfigRaw(config), nil)
					if err != nil {
						t.Fatalf("planning existing pool: %v", err)
					}
					if test.wantReplacement {
						if diff == nil || !diff.RequiresNew() {
							t.Fatalf("expected replacement for a configured mode change, got %#v", diff)
						}
					} else if diff != nil && !diff.Empty() {
						t.Fatalf("expected no changes to existing pool, got %#v", diff)
					}
				})
			}
		})
	}
}

func TestPodIPAllocationModeDefaultNodePoolConversion(t *testing.T) {
	tests := []struct {
		name string
		mode *managedclusters.PodIPAllocationMode
	}{
		{name: "omitted"},
		{name: "dynamic", mode: pointer.To(managedclusters.PodIPAllocationModeDynamicIndividual)},
		{name: "static", mode: pointer.To(managedclusters.PodIPAllocationModeStaticBlock)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			profile := managedclusters.ManagedClusterAgentPoolProfile{Name: "default", PodIPAllocationMode: test.mode}
			result := containers.ConvertDefaultNodePoolToAgentPool(&[]managedclusters.ManagedClusterAgentPoolProfile{profile})
			if test.mode == nil {
				if result.Properties.PodIPAllocationMode != nil {
					t.Fatal("omitted pod IP allocation mode must remain omitted")
				}
			} else if got, want := pointer.FromEnum(result.Properties.PodIPAllocationMode), pointer.FromEnum(test.mode); got != want {
				t.Fatalf("pod IP allocation mode lost during node pool conversion: got %q, want %q", got, want)
			}
		})
	}
}

func podIPAllocationModeTestResources() map[string]*pluginsdk.Resource {
	nodePool := containers.Registration{}.SupportedResources()["azurerm_kubernetes_cluster_node_pool"]
	defaultPool := *containers.SchemaDefaultNodePool()
	defaultFields := defaultPool.Elem.(*pluginsdk.Resource).Schema
	defaultPool.Elem = &pluginsdk.Resource{
		Schema: map[string]*pluginsdk.Schema{
			"pod_ip_allocation_mode": defaultFields["pod_ip_allocation_mode"],
			"pod_subnet_id":          defaultFields["pod_subnet_id"],
		},
	}
	return map[string]*pluginsdk.Resource{
		"node_pool": {
			Schema: map[string]*pluginsdk.Schema{
				"pod_ip_allocation_mode": nodePool.Schema["pod_ip_allocation_mode"],
				"pod_subnet_id":          nodePool.Schema["pod_subnet_id"],
			},
		},
		"default_node_pool": {
			Schema: map[string]*pluginsdk.Schema{"default_node_pool": &defaultPool},
		},
	}
}
