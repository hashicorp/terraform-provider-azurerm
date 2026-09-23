// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package containers

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-05-01/agentpools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-05-01/managedclusters"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func TestKubernetesNodePoolSecurity(t *testing.T) {
	cases := []struct {
		name       string
		config     map[string]interface{}
		secureBoot bool
		vtpm       bool
	}{
		{name: "omitted"},
		{name: "empty", config: map[string]interface{}{}},
		{name: "secure_boot", config: map[string]interface{}{"secure_boot_enabled": true}, secureBoot: true},
		{name: "vtpm", config: map[string]interface{}{"vtpm_enabled": true}, vtpm: true},
		{name: "enabled", config: map[string]interface{}{"secure_boot_enabled": true, "vtpm_enabled": true}, secureBoot: true, vtpm: true},
		{name: "disabled", config: map[string]interface{}{"secure_boot_enabled": false, "vtpm_enabled": false}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			pool := map[string]interface{}{"name": "default", "vm_size": "Standard_D2s_v3", "node_count": 1}
			config := map[string]interface{}{}
			var expectedAgentPool *agentpools.AgentPoolSecurityProfile
			var expectedDefaultPool *managedclusters.AgentPoolSecurityProfile
			expectedState := []interface{}{}
			if tc.config != nil {
				config["security"] = []interface{}{tc.config}
				pool["security"] = []interface{}{tc.config}
				expectedAgentPool = &agentpools.AgentPoolSecurityProfile{
					EnableSecureBoot: pointer.To(tc.secureBoot),
					EnableVTPM:       pointer.To(tc.vtpm),
				}
				expectedDefaultPool = &managedclusters.AgentPoolSecurityProfile{
					EnableSecureBoot: pointer.To(tc.secureBoot),
					EnableVTPM:       pointer.To(tc.vtpm),
				}
				expectedState = []interface{}{map[string]interface{}{"secure_boot_enabled": tc.secureBoot, "vtpm_enabled": tc.vtpm}}
			}

			nodePoolData := schema.TestResourceDataRaw(t, resourceKubernetesClusterNodePool().Schema, config)
			if actual := expandAgentPoolSecurityProfile(nodePoolData.Get("security").([]interface{})); !reflect.DeepEqual(actual, expectedAgentPool) {
				t.Fatalf("unexpected node pool security: %#v", actual)
			}
			if actual := flattenAgentPoolSecurityProfile(expectedAgentPool); !reflect.DeepEqual(actual, expectedState) {
				t.Fatalf("unexpected node pool security state: %#v", actual)
			}

			clusterData := schema.TestResourceDataRaw(t, resourceKubernetesCluster().Schema, map[string]interface{}{
				"default_node_pool": []interface{}{pool},
			})
			profiles, err := ExpandDefaultNodePool(clusterData)
			if err != nil {
				t.Fatal(err)
			}
			if actual := (*profiles)[0].SecurityProfile; !reflect.DeepEqual(actual, expectedDefaultPool) {
				t.Fatalf("unexpected default node pool security: %#v", actual)
			}
			if actual := ConvertDefaultNodePoolToAgentPool(profiles).Properties.SecurityProfile; !reflect.DeepEqual(actual, expectedAgentPool) {
				t.Fatalf("unexpected converted node pool security: %#v", actual)
			}
			flattened, err := FlattenDefaultNodePool(profiles, clusterData)
			if err != nil {
				t.Fatal(err)
			}
			if err := clusterData.Set("default_node_pool", flattened); err != nil {
				t.Fatal(err)
			}
			if actual := clusterData.Get("default_node_pool.0.security"); !reflect.DeepEqual(actual, expectedState) {
				t.Fatalf("unexpected default node pool security state: %#v", actual)
			}
		})
	}
}

func TestKubernetesNodePoolSecurityOmittedPreservesState(t *testing.T) {
	clusterSchema := resourceKubernetesCluster().Schema
	nodePoolSchema := resourceKubernetesClusterNodePool().Schema
	for _, tc := range []struct {
		name   string
		schema map[string]*pluginsdk.Schema
		key    string
		config map[string]interface{}
	}{
		{
			name:   "node_pool",
			schema: map[string]*pluginsdk.Schema{"security": nodePoolSchema["security"]},
			key:    "security",
			config: map[string]interface{}{},
		},
		{
			name:   "default_node_pool",
			schema: map[string]*pluginsdk.Schema{"default_node_pool": clusterSchema["default_node_pool"]},
			key:    "default_node_pool.0.security",
			config: map[string]interface{}{
				"default_node_pool": []interface{}{map[string]interface{}{"name": "default", "vm_size": "Standard_D2s_v3", "node_count": 1}},
			},
		},
	} {
		for _, state := range []struct {
			name    string
			enabled bool
		}{
			{name: "api_defaults"},
			{name: "enabled", enabled: true},
		} {
			t.Run(tc.name+"/"+state.name, func(t *testing.T) {
				resource := &pluginsdk.Resource{Schema: tc.schema}
				data := schema.TestResourceDataRaw(t, tc.schema, tc.config)
				security := flattenAgentPoolSecurityProfile(&agentpools.AgentPoolSecurityProfile{
					EnableSecureBoot: pointer.To(state.enabled),
					EnableVTPM:       pointer.To(state.enabled),
				})
				if tc.name == "node_pool" {
					if err := data.Set("security", security); err != nil {
						t.Fatal(err)
					}
				} else {
					pool := data.Get("default_node_pool").([]interface{})
					pool[0].(map[string]interface{})["security"] = security
					if err := data.Set("default_node_pool", pool); err != nil {
						t.Fatal(err)
					}
				}
				data.SetId("test")
				diff, err := resource.SimpleDiff(context.Background(), data.State(), terraform.NewResourceConfigRaw(tc.config), nil)
				if err != nil {
					t.Fatal(err)
				}
				if diff != nil {
					for key, change := range diff.Attributes {
						if strings.HasPrefix(key, tc.key+".") {
							t.Fatalf("omitted security must preserve API-reported state: %s: %#v", key, change)
						}
					}
				}
			})
		}
	}
}
