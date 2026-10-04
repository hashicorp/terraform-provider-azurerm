// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package containers

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"
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
		config     map[string]any
		secureBoot bool
		vtpm       bool
	}{
		{name: "omitted"},
		{name: "empty", config: map[string]any{}},
		{name: "null_flags", config: map[string]any{"secure_boot_enabled": nil, "vtpm_enabled": nil}},
		{name: "secure_boot", config: map[string]any{"secure_boot_enabled": true}, secureBoot: true},
		{name: "vtpm", config: map[string]any{"vtpm_enabled": true}, vtpm: true},
		{name: "enabled", config: map[string]any{"secure_boot_enabled": true, "vtpm_enabled": true}, secureBoot: true, vtpm: true},
		{name: "disabled", config: map[string]any{"secure_boot_enabled": false, "vtpm_enabled": false}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			pool := map[string]any{"name": "default", "vm_size": "Standard_D2s_v3", "node_count": 1}
			config := map[string]any{}
			var expectedAgentPool *agentpools.AgentPoolSecurityProfile
			var expectedDefaultPool *managedclusters.AgentPoolSecurityProfile
			expectedState := []any{}
			if tc.config != nil {
				config["security"] = []any{tc.config}
				pool["security"] = []any{tc.config}
				expectedAgentPool = &agentpools.AgentPoolSecurityProfile{
					EnableSecureBoot: pointer.To(tc.secureBoot),
					EnableVTPM:       pointer.To(tc.vtpm),
				}
				expectedDefaultPool = &managedclusters.AgentPoolSecurityProfile{
					EnableSecureBoot: pointer.To(tc.secureBoot),
					EnableVTPM:       pointer.To(tc.vtpm),
				}
				expectedState = []any{map[string]any{"secure_boot_enabled": tc.secureBoot, "vtpm_enabled": tc.vtpm}}
			}

			nodePoolData := schema.TestResourceDataRaw(t, resourceKubernetesClusterNodePool().Schema, config)
			if actual := expandAgentPoolSecurityProfile(nodePoolData.Get("security").([]any)); !reflect.DeepEqual(actual, expectedAgentPool) {
				t.Fatalf("unexpected node pool security: %#v", actual)
			}
			if actual := flattenAgentPoolSecurityProfile(expectedAgentPool); !reflect.DeepEqual(actual, expectedState) {
				t.Fatalf("unexpected node pool security state: %#v", actual)
			}

			clusterData := schema.TestResourceDataRaw(t, resourceKubernetesCluster().Schema, map[string]any{
				"default_node_pool": []any{pool},
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
			for name, request := range map[string]any{
				"node_pool":      agentpools.ManagedClusterAgentPoolProfileProperties{SecurityProfile: expandAgentPoolSecurityProfile(nodePoolData.Get("security").([]any))},
				"default_pool":   (*profiles)[0],
				"converted_pool": ConvertDefaultNodePoolToAgentPool(profiles).Properties,
			} {
				body, err := json.Marshal(request)
				if err != nil {
					t.Fatal(err)
				}
				var payload map[string]json.RawMessage
				if err := json.Unmarshal(body, &payload); err != nil {
					t.Fatal(err)
				}
				if tc.config == nil {
					if _, ok := payload["securityProfile"]; ok {
						t.Fatalf("%s must omit securityProfile: %s", name, body)
					}
				} else {
					expected := fmt.Sprintf(`{"enableSecureBoot":%t,"enableVTPM":%t}`, tc.secureBoot, tc.vtpm)
					if actual := string(payload["securityProfile"]); actual != expected {
						t.Fatalf("%s securityProfile = %s, expected %s", name, actual, expected)
					}
				}
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

func TestKubernetesNodePoolSecurityAPI(t *testing.T) {
	for _, tc := range []struct {
		name    string
		profile *managedclusters.AgentPoolSecurityProfile
	}{
		{name: "nil_profile"},
		{name: "nil_flags", profile: &managedclusters.AgentPoolSecurityProfile{}},
		{name: "secure_boot_only", profile: &managedclusters.AgentPoolSecurityProfile{EnableSecureBoot: pointer.To(true)}},
		{name: "vtpm_only", profile: &managedclusters.AgentPoolSecurityProfile{EnableVTPM: pointer.To(true)}},
		{name: "ssh_access", profile: &managedclusters.AgentPoolSecurityProfile{SshAccess: pointer.To(managedclusters.AgentPoolSSHAccessDisabled)}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			converted := ConvertDefaultNodePoolToAgentPool(pointer.To([]managedclusters.ManagedClusterAgentPoolProfile{{SecurityProfile: tc.profile}})).Properties.SecurityProfile
			expected := []any{}
			if tc.profile == nil {
				if converted != nil {
					t.Fatalf("nil security profile must stay nil: %#v", converted)
				}
			} else {
				expected = []any{map[string]any{
					"secure_boot_enabled": pointer.From(tc.profile.EnableSecureBoot),
					"vtpm_enabled":        pointer.From(tc.profile.EnableVTPM),
				}}
				if converted == nil || !reflect.DeepEqual(converted.EnableSecureBoot, tc.profile.EnableSecureBoot) || !reflect.DeepEqual(converted.EnableVTPM, tc.profile.EnableVTPM) {
					t.Fatalf("conversion changed security flag pointers: %#v", converted)
				}
				if pointer.FromEnum(converted.SshAccess) != pointer.FromEnum(tc.profile.SshAccess) {
					t.Fatalf("conversion changed SSH access: %#v", converted)
				}
			}
			for _, actual := range []any{flattenAgentPoolSecurityProfile(converted), flattenManagedClusterAgentPoolSecurityProfile(tc.profile)} {
				if !reflect.DeepEqual(actual, expected) {
					t.Fatalf("unexpected security profile read: %#v, expected %#v", actual, expected)
				}
			}
		})
	}
	if expandAgentPoolSecurityProfile([]any{nil}) != nil || expandManagedClusterAgentPoolSecurityProfile([]any{nil}) != nil {
		t.Fatal("a nil block must expand to a nil profile")
	}
}

func TestKubernetesNodePoolSecurityDiff(t *testing.T) {
	for _, surface := range []struct {
		name   string
		schema map[string]*pluginsdk.Schema
		key    string
		config func(map[string]any) map[string]any
	}{
		{
			name:   "node_pool",
			schema: map[string]*pluginsdk.Schema{"security": resourceKubernetesClusterNodePool().Schema["security"]},
			key:    "security",
			config: func(security map[string]any) map[string]any {
				if security == nil {
					return map[string]any{}
				}
				return map[string]any{"security": []any{security}}
			},
		},
		{
			name:   "default_node_pool",
			schema: map[string]*pluginsdk.Schema{"default_node_pool": resourceKubernetesCluster().Schema["default_node_pool"]},
			key:    "default_node_pool.0.security",
			config: func(security map[string]any) map[string]any {
				pool := map[string]any{"name": "default", "vm_size": "Standard_D2s_v3", "node_count": 1}
				if security != nil {
					pool["security"] = []any{security}
				}
				return map[string]any{"default_node_pool": []any{pool}}
			},
		},
	} {
		for _, tc := range []struct {
			name       string
			security   map[string]any
			secureBoot bool
			vtpm       bool
		}{
			{name: "removed"},
			{name: "enabled", security: map[string]any{"secure_boot_enabled": true, "vtpm_enabled": true}, secureBoot: true, vtpm: true},
			{name: "disabled", security: map[string]any{"secure_boot_enabled": false, "vtpm_enabled": false}},
			{name: "empty", security: map[string]any{}},
			{name: "null_flags", security: map[string]any{"secure_boot_enabled": nil, "vtpm_enabled": nil}},
			{name: "secure_boot_only", security: map[string]any{"secure_boot_enabled": true}, secureBoot: true},
			{name: "vtpm_only", security: map[string]any{"vtpm_enabled": true}, vtpm: true},
		} {
			for _, enabled := range []bool{false, true} {
				t.Run(fmt.Sprintf("%s/initial_%t/%s", surface.name, enabled, tc.name), func(t *testing.T) {
					resource := &pluginsdk.Resource{Schema: surface.schema}
					initial := surface.config(map[string]any{"secure_boot_enabled": enabled, "vtpm_enabled": enabled})
					data := schema.TestResourceDataRaw(t, surface.schema, initial)
					data.SetId("test")
					diff, err := resource.SimpleDiff(context.Background(), data.State(), terraform.NewResourceConfigRaw(surface.config(tc.security)), nil)
					if err != nil {
						t.Fatal(err)
					}
					expected := map[string]*terraform.ResourceAttrDiff{}
					for flag, value := range map[string]bool{"secure_boot_enabled": tc.secureBoot, "vtpm_enabled": tc.vtpm} {
						if tc.security != nil && value != enabled {
							expected[surface.key+".0."+flag] = &terraform.ResourceAttrDiff{Old: fmt.Sprint(enabled), New: fmt.Sprint(value)}
						}
					}
					actual := map[string]*terraform.ResourceAttrDiff{}
					if diff != nil {
						if diff.Destroy || diff.RequiresNew() {
							t.Fatalf("security changes must not replace the Terraform resource: %#v", diff)
						}
						maps.Copy(actual, diff.Attributes)
					}
					if !reflect.DeepEqual(actual, expected) {
						t.Fatalf("unexpected diff: %#v, expected %#v", actual, expected)
					}
				})
			}
		}
	}
}

func TestKubernetesNodePoolSecurityOmittedPreservesState(t *testing.T) {
	clusterSchema := resourceKubernetesCluster().Schema
	nodePoolSchema := resourceKubernetesClusterNodePool().Schema
	for _, tc := range []struct {
		name   string
		schema map[string]*pluginsdk.Schema
		key    string
		config map[string]any
	}{
		{
			name:   "node_pool",
			schema: map[string]*pluginsdk.Schema{"security": nodePoolSchema["security"]},
			key:    "security",
			config: map[string]any{},
		},
		{
			name:   "default_node_pool",
			schema: map[string]*pluginsdk.Schema{"default_node_pool": clusterSchema["default_node_pool"]},
			key:    "default_node_pool.0.security",
			config: map[string]any{
				"default_node_pool": []any{map[string]any{"name": "default", "vm_size": "Standard_D2s_v3", "node_count": 1}},
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
					pool := data.Get("default_node_pool").([]any)
					pool[0].(map[string]any)["security"] = security
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
