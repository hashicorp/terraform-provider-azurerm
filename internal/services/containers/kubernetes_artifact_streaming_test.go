// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package containers_test

import (
	"context"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-04-01/managedclusters"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func TestArtifactStreamingSchema(t *testing.T) {
	schemas := map[string]*pluginsdk.Schema{
		"default_node_pool": containers.SchemaDefaultNodePool().Elem.(*pluginsdk.Resource).Schema["artifact_streaming_enabled"],
		"node_pool":         containers.Registration{}.SupportedResources()["azurerm_kubernetes_cluster_node_pool"].Schema["artifact_streaming_enabled"],
	}
	for name, schema := range schemas {
		t.Run(name, func(t *testing.T) {
			if schema == nil || schema.Type != pluginsdk.TypeBool || !schema.Optional || schema.ForceNew || schema.Computed {
				t.Fatal("artifact streaming must be an optional, updatable boolean on both node pool schemas")
			}
		})
	}
}

func TestArtifactStreamingDefaultNodePoolExpand(t *testing.T) {
	for _, enabled := range []bool{true, false} {
		data := schema.TestResourceDataRaw(t, map[string]*pluginsdk.Schema{
			"default_node_pool": containers.SchemaDefaultNodePool(),
		}, map[string]interface{}{
			"default_node_pool": []interface{}{map[string]interface{}{
				"name":                       "default",
				"node_count":                 1,
				"vm_size":                    "Standard_D2s_v3",
				"artifact_streaming_enabled": enabled,
			}},
		})
		data.MarkNewResource()
		profiles, err := containers.ExpandDefaultNodePool(data)
		if err != nil {
			t.Fatal(err)
		}
		profile := (*profiles)[0].ArtifactStreamingProfile
		if enabled {
			if profile == nil || !pointer.From(profile.Enabled) {
				t.Fatal("enabled artifact streaming must be sent during creation")
			}
		} else if profile != nil {
			t.Fatal("disabled artifact streaming must be omitted during creation")
		}
	}
}

func TestArtifactStreamingDefaultNodePoolRoundTrip(t *testing.T) {
	for name, profile := range map[string]*managedclusters.AgentPoolArtifactStreamingProfile{
		"omitted":  nil,
		"empty":    {},
		"enabled":  {Enabled: pointer.To(true)},
		"disabled": {Enabled: pointer.To(false)},
	} {
		t.Run(name, func(t *testing.T) {
			profiles := []managedclusters.ManagedClusterAgentPoolProfile{{
				Name:                     "default",
				Mode:                     pointer.To(managedclusters.AgentPoolModeSystem),
				ArtifactStreamingProfile: profile,
			}}
			pool := containers.ConvertDefaultNodePoolToAgentPool(&profiles)
			if profile == nil {
				if pool.Properties.ArtifactStreamingProfile != nil {
					t.Fatal("omitted artifact streaming profile must remain omitted")
				}
			} else if got := pool.Properties.ArtifactStreamingProfile; got == nil || got.Enabled != profile.Enabled {
				t.Fatal("conversion must preserve the artifact streaming value, including explicit false and nil")
			}
			data := schema.TestResourceDataRaw(t, map[string]*pluginsdk.Schema{
				"default_node_pool": containers.SchemaDefaultNodePool(),
			}, map[string]interface{}{
				"default_node_pool": []interface{}{map[string]interface{}{"name": "default"}},
			})
			flattened, err := containers.FlattenDefaultNodePool(&profiles, data)
			if err != nil {
				t.Fatal(err)
			}
			want := profile != nil && pointer.From(profile.Enabled)
			if got := (*flattened)[0].(map[string]interface{})["artifact_streaming_enabled"]; got != want {
				t.Fatalf("unexpected flattened value: got %v, want %t", got, want)
			}
		})
	}
}

func TestArtifactStreamingDefaultNodePoolDisable(t *testing.T) {
	var updated *managedclusters.AgentPoolArtifactStreamingProfile
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{"default_node_pool": containers.SchemaDefaultNodePool()},
		Update: func(data *schema.ResourceData, _ interface{}) error {
			profiles, err := containers.ExpandDefaultNodePool(data)
			if err == nil {
				updated = (*profiles)[0].ArtifactStreamingProfile
			}
			return err
		},
	}
	raw := map[string]interface{}{
		"default_node_pool": []interface{}{map[string]interface{}{
			"name":                       "default",
			"vm_size":                    "Standard_D2s_v3",
			"node_count":                 1,
			"artifact_streaming_enabled": true,
		}},
	}
	data := schema.TestResourceDataRaw(t, resource.Schema, raw)
	data.SetId("cluster")
	state := data.State()
	raw["default_node_pool"].([]interface{})[0].(map[string]interface{})["artifact_streaming_enabled"] = false
	diff, err := resource.Diff(context.Background(), state, terraform.NewResourceConfigRaw(raw), nil)
	if err != nil {
		t.Fatal(err)
	}
	if diff == nil || diff.RequiresNew() {
		t.Fatal("disabling artifact streaming must produce an in-place update")
	}
	if _, diags := resource.Apply(context.Background(), state, diff, nil); diags.HasError() {
		t.Fatal(diags)
	}
	if updated == nil || updated.Enabled == nil || *updated.Enabled {
		t.Fatal("disabling artifact streaming must send an explicit false")
	}
}
