// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package mongocluster_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mongocluster"
)

// These tests exercise the `storage_type` schema behaviour offline, using the plugin
// SDK's diff machinery directly. They require neither Azure credentials nor TF_ACC.
//
// Background: `storage_type` previously carried a static `Default` of "PremiumSSD",
// which meant the provider transmitted an explicit storage type even when the
// practitioner never configured one, preventing the service from applying its own
// contextual default. Removing that `Default` is only safe if `Computed` is added at
// the same time -- `storage_type` is `ForceNew`, so a configuration that omits it must
// not be seen as "changed to empty", which would plan a destroy-and-recreate of an
// existing cluster.
//
// TestAccMongoCluster_storageTypeOmittedNoReplacementOnUpgrade is therefore the
// regression guard for the single highest-risk scenario: upgrading the provider with
// no configuration change at all must not replace existing clusters.

// priorStateForExistingCluster models the state an existing cluster would have on disk,
// as written by a provider version that applied the static "PremiumSSD" default.
func priorStateForExistingCluster(storageType string) *terraform.InstanceState {
	attributes := map[string]string{
		"id":                     "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.DocumentDB/mongoClusters/example",
		"name":                   "example",
		"resource_group_name":    "example",
		"location":               "westeurope",
		"compute_tier":           "M30",
		"high_availability_mode": "Disabled",
		"shard_count":            "1",
		"storage_size_in_gb":     "64",
		"version":                "7.0",
		"create_mode":            "Default",
		"public_network_access":  "Enabled",
	}
	if storageType != "" {
		attributes["storage_type"] = storageType
	}

	return &terraform.InstanceState{
		ID:         attributes["id"],
		Attributes: attributes,
	}
}

// configForExistingCluster mirrors the practitioner's HCL. `storage_type` is only
// present when explicitlySetStorageType is non-empty, which is what lets these tests
// distinguish "omitted in configuration" from "explicitly configured".
func configForExistingCluster(explicitlySetStorageType string) *terraform.ResourceConfig {
	raw := map[string]interface{}{
		"name":                   "example",
		"resource_group_name":    "example",
		"location":               "westeurope",
		"compute_tier":           "M30",
		"high_availability_mode": "Disabled",
		"shard_count":            1,
		"storage_size_in_gb":     64,
		"version":                "7.0",
		"create_mode":            "Default",
		"public_network_access":  "Enabled",
	}
	if explicitlySetStorageType != "" {
		raw["storage_type"] = explicitlySetStorageType
	}

	return terraform.NewResourceConfigRaw(raw)
}

func diffForMongoCluster(t *testing.T, state *terraform.InstanceState, config *terraform.ResourceConfig) *terraform.InstanceDiff {
	t.Helper()

	// Diff against the resource's schema only. The typed wrapper's CustomizeDiff
	// requires a fully constructed *clients.Client, which would turn this into an
	// integration test; the behaviour under test here (`storage_type` being
	// Optional+Computed+ForceNew with no static default) is decided entirely by the
	// schema, so a bare schema.Resource isolates it precisely.
	wrapped := sdk.WrappedResource(mongocluster.MongoClusterResource{})
	resource := &schema.Resource{Schema: wrapped.Schema}

	diff, err := resource.Diff(context.Background(), state, config, nil)
	if err != nil {
		t.Fatalf("generating diff: %+v", err)
	}

	return diff
}

// assertNoStorageTypeReplacement fails the test when the generated diff would destroy
// and recreate the cluster because of `storage_type`.
func assertNoStorageTypeReplacement(t *testing.T, diff *terraform.InstanceDiff) {
	t.Helper()

	if diff == nil {
		return
	}

	if diff.RequiresNew() {
		t.Fatalf("expected no replacement, but the diff requires the resource to be recreated: %#v", diff.Attributes)
	}

	if attr, ok := diff.Attributes["storage_type"]; ok && attr.RequiresNew {
		t.Fatalf("expected no replacement, but `storage_type` is marked ForceNew (old: %q, new: %q)", attr.Old, attr.New)
	}
}

// Scenario A1 (release blocking): an existing cluster created while the static default
// was in place holds "PremiumSSD" in state, and the configuration never mentioned
// `storage_type`. Upgrading the provider must not plan a replacement.
func TestMongoClusterStorageType_omittedInConfigDoesNotReplaceExistingCluster(t *testing.T) {
	state := priorStateForExistingCluster("PremiumSSD")
	config := configForExistingCluster("")

	diff := diffForMongoCluster(t, state, config)

	assertNoStorageTypeReplacement(t, diff)

	if diff != nil {
		if attr, ok := diff.Attributes["storage_type"]; ok && attr.New != "" && attr.New != "PremiumSSD" {
			t.Fatalf("expected `storage_type` to stay PremiumSSD, got %q", attr.New)
		}
	}
}

// Scenario A5: the same guarantee must hold for a cluster that was explicitly created
// as PremiumSSDv2. Adding `Computed` must not drag it back to a default.
func TestMongoClusterStorageType_omittedInConfigDoesNotReplaceExistingSSDv2Cluster(t *testing.T) {
	state := priorStateForExistingCluster("PremiumSSDv2")
	config := configForExistingCluster("")

	diff := diffForMongoCluster(t, state, config)

	assertNoStorageTypeReplacement(t, diff)

	if diff != nil {
		if attr, ok := diff.Attributes["storage_type"]; ok && attr.New != "" && attr.New != "PremiumSSDv2" {
			t.Fatalf("expected `storage_type` to stay PremiumSSDv2, got %q", attr.New)
		}
	}
}

// Scenario A2: a configuration that explicitly pins PremiumSSD must remain stable --
// no diff and no replacement -- across the provider upgrade.
func TestMongoClusterStorageType_explicitlyConfiguredValueIsStable(t *testing.T) {
	state := priorStateForExistingCluster("PremiumSSD")
	config := configForExistingCluster("PremiumSSD")

	diff := diffForMongoCluster(t, state, config)

	assertNoStorageTypeReplacement(t, diff)
}

// Scenario C2: removing an explicitly configured `storage_type` from the configuration
// must be a no-op rather than a replacement. `Computed` should retain the value that is
// already in state.
func TestMongoClusterStorageType_removingExplicitValueFromConfigDoesNotReplace(t *testing.T) {
	state := priorStateForExistingCluster("PremiumSSD")
	config := configForExistingCluster("")

	diff := diffForMongoCluster(t, state, config)

	assertNoStorageTypeReplacement(t, diff)
}

// Scenario C1: changing `storage_type` explicitly is a genuine storage migration, which
// the service cannot perform in place. That must still be surfaced as a replacement --
// this asserts the safety property that only an explicit change is destructive.
func TestMongoClusterStorageType_explicitChangeStillForcesReplacement(t *testing.T) {
	state := priorStateForExistingCluster("PremiumSSD")
	config := configForExistingCluster("PremiumSSDv2")

	diff := diffForMongoCluster(t, state, config)

	if diff == nil {
		t.Fatal("expected a diff when `storage_type` is explicitly changed, got none")
	}

	attr, ok := diff.Attributes["storage_type"]
	if !ok {
		t.Fatalf("expected `storage_type` in the diff, got: %#v", diff.Attributes)
	}

	if !attr.RequiresNew {
		t.Fatalf("expected an explicit `storage_type` change to force replacement, got: %#v", attr)
	}
}

// A fresh resource with no `storage_type` in configuration must leave the value unknown
// so the service can decide it, rather than baking a client-side default into the plan.
func TestMongoClusterStorageType_newResourceOmittedTypeIsNotDefaultedClientSide(t *testing.T) {
	config := configForExistingCluster("")

	diff := diffForMongoCluster(t, nil, config)

	if diff == nil {
		t.Fatal("expected a diff for a new resource, got none")
	}

	attr, ok := diff.Attributes["storage_type"]
	if !ok {
		// Not present in the diff at all is also acceptable: nothing is being asserted
		// about the storage type, so the service remains free to choose it.
		return
	}

	if attr.New == "PremiumSSD" && !attr.NewComputed {
		t.Fatal("expected `storage_type` to be service-determined for a new resource, but the provider planned a client-side PremiumSSD default")
	}
}

// An in-place update must not disturb the storage type, whichever type the cluster has.
//
// Note what this does *not* cover: a replacement. When a diff requires a new resource,
// helper/schema discards the state and re-runs the diff and CustomizeDiff from scratch,
// and the prior value is unreachable at that point — both the state and the raw state are
// gone. So the recreated cluster is planned with `storage_type` unknown and takes whatever
// the service chooses. That is verified harmless today (live run R1: the cluster came back
// PremiumSSD), but it means the storage type cannot be carried across a replacement from
// within the provider. See docs in the Marlin repo for the Phase 2 implications.
func TestMongoClusterStorageType_inPlaceUpdateDoesNotDisturbExistingType(t *testing.T) {
	for _, storageType := range []string{"PremiumSSD", "PremiumSSDv2"} {
		t.Run(storageType, func(t *testing.T) {
			state := priorStateForExistingCluster(storageType)
			config := configForExistingCluster("")

			diff := diffForMongoCluster(t, state, config)

			assertNoStorageTypeReplacement(t, diff)

			if diff != nil {
				if attr, ok := diff.Attributes["storage_type"]; ok {
					if attr.New != "" && attr.New != storageType {
						t.Fatalf("expected `storage_type` to stay %q, got %q", storageType, attr.New)
					}
					if attr.NewRemoved {
						t.Fatalf("expected `storage_type` to be retained, but it was marked removed")
					}
				}
			}
		})
	}
}
