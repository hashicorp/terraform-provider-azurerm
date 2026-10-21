// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package kubernetes_cluster

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2025-11-01/registries"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-05-01/managedclusters"
	"github.com/hashicorp/go-cty/cty"
)

func TestBootstrapProfileContainerRegistryIDConfigured(t *testing.T) {
	registryID := registries.NewRegistryID("00000000-0000-0000-0000-000000000000", "MC_cluster", "managed").ID()
	profileType := cty.Object(map[string]cty.Type{"container_registry_id": cty.String})
	for _, test := range []struct {
		name            string
		config          cty.Value
		stateRegistryID string
		want            bool
	}{
		{name: "omitted profile", config: cty.EmptyObjectVal},
		{name: "omitted registry ID", config: cty.ObjectVal(map[string]cty.Value{"bootstrap_profile": cty.ListVal([]cty.Value{cty.ObjectVal(map[string]cty.Value{"container_registry_id": cty.NullVal(cty.String)})})})},
		{name: "configured registry ID", config: cty.ObjectVal(map[string]cty.Value{"bootstrap_profile": cty.ListVal([]cty.Value{cty.ObjectVal(map[string]cty.Value{"container_registry_id": cty.StringVal(registryID)})})}), want: true},
		{name: "null config with configured state", config: cty.NullVal(cty.EmptyObject), stateRegistryID: registryID, want: true},
		{name: "null config without configured state", config: cty.NullVal(cty.EmptyObject)},
		{name: "omitted registry ID ignores state", config: cty.ObjectVal(map[string]cty.Value{"bootstrap_profile": cty.ListVal([]cty.Value{cty.ObjectVal(map[string]cty.Value{"container_registry_id": cty.NullVal(cty.String)})})}), stateRegistryID: registryID},
		{name: "unknown profile element uses state", config: cty.ObjectVal(map[string]cty.Value{"bootstrap_profile": cty.ListVal([]cty.Value{cty.UnknownVal(profileType)})}), stateRegistryID: registryID, want: true},
		{name: "unknown profile element without state", config: cty.ObjectVal(map[string]cty.Value{"bootstrap_profile": cty.ListVal([]cty.Value{cty.UnknownVal(profileType)})})},
		{name: "null profile element uses state", config: cty.ObjectVal(map[string]cty.Value{"bootstrap_profile": cty.ListVal([]cty.Value{cty.NullVal(profileType)})}), stateRegistryID: registryID, want: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			if got := bootstrapProfileContainerRegistryIDConfigured(test.config, test.stateRegistryID); got != test.want {
				t.Fatalf("configured registry ID = %t, want %t", got, test.want)
			}
		})
	}
}

func TestFlattenBootstrapProfileContainerRegistryID(t *testing.T) {
	managedRegistryID := registries.NewRegistryID("00000000-0000-0000-0000-000000000000", "MC_cluster", "managed").ID()
	userRegistryID := registries.NewRegistryID("00000000-0000-0000-0000-000000000000", "user", "custom").ID()

	for _, test := range []struct {
		name                 string
		registryID           string
		nodeSubscriptionID   string
		nodeResourceGroup    string
		configuredRegistryID bool
		want                 string
	}{
		{name: "managed registry in resource", registryID: managedRegistryID, nodeSubscriptionID: "00000000-0000-0000-0000-000000000000", nodeResourceGroup: "MC_cluster"},
		{name: "user registry in resource", registryID: userRegistryID, nodeSubscriptionID: "00000000-0000-0000-0000-000000000000", nodeResourceGroup: "MC_cluster", want: userRegistryID},
		{name: "configured registry in node resource group", registryID: managedRegistryID, nodeSubscriptionID: "00000000-0000-0000-0000-000000000000", nodeResourceGroup: "MC_cluster", configuredRegistryID: true, want: managedRegistryID},
		{name: "other subscription in resource", registryID: managedRegistryID, nodeSubscriptionID: "11111111-1111-1111-1111-111111111111", nodeResourceGroup: "MC_cluster", want: managedRegistryID},
		{name: "managed registry in data source", registryID: managedRegistryID, want: managedRegistryID},
	} {
		t.Run(test.name, func(t *testing.T) {
			profile := &managedclusters.ManagedClusterBootstrapProfile{
				ArtifactSource:      pointer.To(managedclusters.ArtifactSourceCache),
				ContainerRegistryId: pointer.To(test.registryID),
			}
			flattened, err := flattenBootstrapProfile(profile, test.nodeSubscriptionID, test.nodeResourceGroup, test.configuredRegistryID)
			if err != nil {
				t.Fatal(err)
			}
			got := flattened[0].(map[string]any)["container_registry_id"].(string)
			if got != test.want {
				t.Fatalf("container_registry_id = %q, want %q", got, test.want)
			}
		})
	}
}
