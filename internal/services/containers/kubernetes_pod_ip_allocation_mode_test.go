// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package containers_test

import (
	"strings"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-05-01/managedclusters"
	"github.com/hashicorp/go-cty/cty"
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
		{name: "static_with_null_subnet", config: map[string]interface{}{"pod_ip_allocation_mode": static, "pod_subnet_id": nil}, wantError: true},
		{name: "empty_mode", config: map[string]interface{}{"pod_ip_allocation_mode": "", "pod_subnet_id": subnetId}, wantError: true},
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

func TestPodIPAllocationModeUnknownValidation(t *testing.T) {
	subnetId := commonids.NewSubnetID("00000000-0000-0000-0000-000000000000", "test", "test", "pods").ID()
	tests := []struct {
		name   string
		mode   cty.Value
		subnet cty.Value
	}{
		{name: "unknown_mode", mode: cty.UnknownVal(cty.String), subnet: cty.StringVal(subnetId)},
		{name: "unknown_subnet", mode: cty.StringVal("StaticBlock"), subnet: cty.UnknownVal(cty.String)},
		{name: "both_unknown", mode: cty.UnknownVal(cty.String), subnet: cty.UnknownVal(cty.String)},
	}
	for name, resource := range podIPAllocationModeTestResources() {
		t.Run(name, func(t *testing.T) {
			for _, test := range tests {
				t.Run(test.name, func(t *testing.T) {
					config := cty.ObjectVal(map[string]cty.Value{
						"pod_ip_allocation_mode": test.mode,
						"pod_subnet_id":          test.subnet,
					})
					if name == "default_node_pool" {
						config = cty.ObjectVal(map[string]cty.Value{"default_node_pool": cty.ListVal([]cty.Value{config})})
					}
					diagnostics := resource.Validate(terraform.NewResourceConfigShimmed(config, resource.CoreConfigSchema()))
					if diagnostics.HasError() {
						t.Fatalf("known-after-apply configuration must validate: %v", diagnostics)
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
		{name: "omitted_with_populated_api_value_replaces_resource", stateMode: static, wantReplacement: true},
		{name: "configured_value_unchanged", stateMode: static, configuredMode: static},
		{name: "configured_change_replaces_resource", stateMode: static, configuredMode: dynamic, wantReplacement: true},
	}

	resources := containers.Registration{}.SupportedResources()
	for _, resourceName := range []string{"azurerm_kubernetes_cluster_node_pool", "azurerm_kubernetes_cluster"} {
		resource := resources[resourceName]
		t.Run(resourceName, func(t *testing.T) {
			for _, test := range tests {
				t.Run(test.name, func(t *testing.T) {
					pool := map[string]interface{}{
						"name": "test", "vm_size": "Standard_D2s_v3", "node_count": 1,
						"pod_subnet_id": subnetId,
					}
					if test.stateMode != "" {
						pool["pod_ip_allocation_mode"] = test.stateMode
					}
					clusterId := commonids.NewKubernetesClusterID("00000000-0000-0000-0000-000000000000", "test", "test").ID()
					resourceId := clusterId + "/agentPools/test"
					config := pool
					prefix := ""
					if resourceName == "azurerm_kubernetes_cluster" {
						config = map[string]interface{}{
							"name": "test", "location": "eastus", "resource_group_name": "test", "dns_prefix": "test",
							"identity":          []interface{}{map[string]interface{}{"type": "SystemAssigned"}},
							"network_profile":   []interface{}{map[string]interface{}{"network_plugin": "azure"}},
							"default_node_pool": []interface{}{pool},
						}
						resourceId = clusterId
						prefix = "default_node_pool.0."
					} else {
						config["kubernetes_cluster_id"] = clusterId
					}
					createDiff, err := resource.SimpleDiff(t.Context(), nil, terraform.NewResourceConfigRaw(config), nil)
					if err != nil {
						t.Fatalf("building initial diff: %v", err)
					}
					// Existing state cannot contain plan-time unknowns; model computed collections as empty.
					for key, change := range createDiff.Attributes {
						if change.NewComputed {
							if strings.HasSuffix(key, ".#") || strings.HasSuffix(key, ".%") {
								change.New, change.NewComputed = "0", false
							} else {
								delete(createDiff.Attributes, key)
							}
						}
					}
					attributes, err := createDiff.Apply(nil, resource.CoreConfigSchema())
					if err != nil {
						t.Fatalf("building initial state: %v", err)
					}
					attributes[prefix+"pod_ip_allocation_mode"] = test.stateMode
					state := &terraform.InstanceState{ID: resourceId, Attributes: attributes}
					if test.configuredMode == "" {
						delete(pool, "pod_ip_allocation_mode")
					} else {
						pool["pod_ip_allocation_mode"] = test.configuredMode
					}
					diff, err := resource.SimpleDiff(t.Context(), state, terraform.NewResourceConfigRaw(config), nil)
					if err != nil {
						t.Fatalf("planning existing pool: %v", err)
					}
					if test.wantReplacement {
						if diff == nil || !diff.RequiresNew() {
							t.Fatalf("expected resource replacement for a changed or removed mode, got %#v", diff)
						}
						if change := diff.Attributes[prefix+"pod_ip_allocation_mode"]; change == nil || !change.RequiresNew || change.Old != test.stateMode || change.New != test.configuredMode {
							t.Fatalf("expected mode change %q -> %q to require replacement, got %#v", test.stateMode, test.configuredMode, change)
						}
					} else if diff != nil && !diff.Empty() {
						t.Fatalf("expected no changes to existing pool, got %#v", diff)
					}
				})
			}
		})
	}
}

func TestPodIPAllocationModeDefaultNodePoolExpansion(t *testing.T) {
	for _, mode := range []string{"", "DynamicIndividual", "StaticBlock"} {
		t.Run("mode_"+mode, func(t *testing.T) {
			config := map[string]interface{}{
				"name": "default", "vm_size": "Standard_D2s_v3", "node_count": 1,
				"pod_subnet_id": commonids.NewSubnetID("00000000-0000-0000-0000-000000000000", "test", "test", "pods").ID(),
			}
			if mode != "" {
				config["pod_ip_allocation_mode"] = mode
			}
			resource := &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{"default_node_pool": containers.SchemaDefaultNodePool()}}
			data := resource.Data(nil)
			if err := data.Set("default_node_pool", []interface{}{config}); err != nil {
				t.Fatal(err)
			}
			profiles, err := containers.ExpandDefaultNodePool(data)
			if err != nil {
				t.Fatal(err)
			}
			got := (*profiles)[0].PodIPAllocationMode
			if pointer.FromEnum(got) != mode || (got == nil) != (mode == "") {
				t.Fatalf("expanded mode = %v, want %q (nil when omitted)", got, mode)
			}
			flattened, err := containers.FlattenDefaultNodePool(profiles, data)
			if err != nil {
				t.Fatal(err)
			}
			if got := (*flattened)[0].(map[string]interface{})["pod_ip_allocation_mode"]; got != mode {
				t.Fatalf("flattened mode = %v, want %q", got, mode)
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
