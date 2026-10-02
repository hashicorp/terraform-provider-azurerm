// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package containers

import (
	"strings"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-05-01/agentpools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAgentPoolWindowsProfileImportPlan(t *testing.T) {
	for _, test := range []struct {
		name         string
		apiDisabled  bool
		configured   *bool
		emptyBlock   bool
		updateTags   bool
		replacement  bool
		noProfile    bool
		noFlag       bool
		optionalOnly bool
	}{
		{name: "omitted_enabled"},
		{name: "omitted_disabled", apiDisabled: true},
		{name: "explicit_enabled", configured: pointer.To(true)},
		{name: "explicit_disabled", apiDisabled: true, configured: pointer.To(false)},
		{name: "empty_block", emptyBlock: true},
		{name: "omitted_tag_update", updateTags: true},
		{name: "omitted_disabled_tag_update", apiDisabled: true, updateTags: true},
		{name: "enable_requires_replacement", apiDisabled: true, configured: pointer.To(true), replacement: true},
		{name: "disable_requires_replacement", configured: pointer.To(false), replacement: true},
		{name: "empty_block_disabled_requires_replacement", apiDisabled: true, emptyBlock: true, replacement: true},
		{name: "omitted_nil_profile", noProfile: true},
		{name: "omitted_nil_profile_tag_update", noProfile: true, updateTags: true},
		{name: "omitted_nil_flag", noFlag: true},
		{name: "omitted_nil_flag_tag_update", noFlag: true, updateTags: true},
		{name: "optional_only_omitted_enabled_requires_replacement", optionalOnly: true, replacement: true},
		{name: "optional_only_omitted_disabled_requires_replacement", optionalOnly: true, apiDisabled: true, replacement: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			resource := resourceKubernetesClusterNodePool()
			if test.optionalOnly {
				resource.Schema["windows_profile"].Computed = false
			}
			config := map[string]any{
				"name":                  "test",
				"kubernetes_cluster_id": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test/providers/Microsoft.ContainerService/managedClusters/test",
				"vm_size":               "Standard_D2s_v3",
				"os_type":               "Windows",
				"node_count":            1,
				"tags":                  map[string]any{"Environment": "dev"},
			}
			// Import has no prior Windows profile from which to infer configuration.
			data := schema.TestResourceDataRaw(t, resource.Schema, config)
			data.SetId(config["kubernetes_cluster_id"].(string) + "/agentPools/test")
			if err := data.Set("node_labels", map[string]any{}); err != nil {
				t.Fatal(err)
			}
			api := expandAgentPoolWindowsProfile([]any{map[string]any{"outbound_nat_enabled": !test.apiDisabled}})
			if api == nil || api.DisableOutboundNat == nil || *api.DisableOutboundNat != test.apiDisabled {
				t.Fatalf("Windows profile did not expand to the expected API flag: %#v", api)
			}
			if test.noProfile {
				api = nil
			} else if test.noFlag {
				api = &agentpools.AgentPoolWindowsProfile{}
			}
			profile := flattenAgentPoolWindowsProfile(api)
			if test.noProfile || test.noFlag {
				if len(profile) != 0 {
					t.Fatalf("absent API outbound NAT setting produced a profile: %#v", profile)
				}
			} else if len(profile) != 1 || profile[0].(map[string]any)["outbound_nat_enabled"] != !test.apiDisabled {
				t.Fatalf("API Windows profile was not preserved in imported state: %#v", profile)
			}
			if err := data.Set("windows_profile", profile); err != nil {
				t.Fatal(err)
			}
			state := data.State()
			if test.configured != nil {
				config["windows_profile"] = []any{map[string]any{"outbound_nat_enabled": *test.configured}}
			} else if test.emptyBlock {
				config["windows_profile"] = []any{map[string]any{}}
			}
			if test.updateTags {
				config["tags"] = map[string]any{"Environment": "prod"}
			}
			diff, err := resource.SimpleDiff(t.Context(), state, terraform.NewResourceConfigRaw(config), nil)
			if err != nil {
				t.Fatal(err)
			}
			if got := diff != nil && diff.RequiresNew(); got != test.replacement {
				t.Fatalf("replacement = %t, want %t; diff: %#v", got, test.replacement, diff)
			}
			if test.replacement {
				profileRequiresNew := false
				for name, attribute := range diff.Attributes {
					if strings.HasPrefix(name, "windows_profile.") && attribute.RequiresNew {
						profileRequiresNew = true
					}
				}
				if !profileRequiresNew {
					t.Fatalf("replacement was not caused by the Windows profile: %#v", diff)
				}
			}
			if !test.replacement && !test.updateTags && diff != nil && !diff.Empty() {
				t.Fatalf("imported Windows profile has a nonempty plan: %#v", diff)
			}
			if test.updateTags {
				if diff == nil || diff.Attributes["tags.Environment"] == nil || diff.Attributes["tags.Environment"].New != "prod" {
					t.Fatalf("expected the tag update in the plan: %#v", diff)
				}
				for name := range diff.Attributes {
					if !strings.HasPrefix(name, "tags.") {
						t.Fatalf("omitted Windows profile changed during a tag update: %s", name)
					}
				}
			}
		})
	}
}
