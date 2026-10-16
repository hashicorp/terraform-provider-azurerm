// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_registry

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2025-11-01/registries"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2025-11-01/replications"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandNetworkRuleSet(profiles []any) *registries.NetworkRuleSet {
	if len(profiles) == 0 {
		return nil
	}

	profile := profiles[0].(map[string]any)

	ipRuleConfigs := profile["ip_rule"].(*pluginsdk.Set).List()
	ipRules := make([]registries.IPRule, 0)
	for _, ipRuleInterface := range ipRuleConfigs {
		config := ipRuleInterface.(map[string]any)
		newIpRule := registries.IPRule{
			Action: pointer.ToEnum[registries.Action](config["action"].(string)),
			Value:  config["ip_range"].(string),
		}
		ipRules = append(ipRules, newIpRule)
	}

	return &registries.NetworkRuleSet{
		DefaultAction: registries.DefaultAction(profile["default_action"].(string)),
		IPRules:       &ipRules,
	}
}

func expandQuarantinePolicy(enabled bool) *registries.QuarantinePolicy {
	quarantinePolicy := registries.QuarantinePolicy{
		Status: pointer.To(registries.PolicyStatusDisabled),
	}

	if enabled {
		quarantinePolicy.Status = pointer.To(registries.PolicyStatusEnabled)
	}

	return &quarantinePolicy
}

func expandExportPolicy(enabled bool) *registries.ExportPolicy {
	exportPolicy := registries.ExportPolicy{
		Status: pointer.To(registries.ExportPolicyStatusDisabled),
	}

	if enabled {
		exportPolicy.Status = pointer.To(registries.ExportPolicyStatusEnabled)
	}

	return &exportPolicy
}

func expandAadAuthAsArmPolicy(enabled bool) *registries.AzureADAuthenticationAsArmPolicy {
	policy := registries.AzureADAuthenticationAsArmPolicy{
		Status: pointer.To(registries.AzureADAuthenticationAsArmPolicyStatusDisabled),
	}

	if enabled {
		policy.Status = pointer.To(registries.AzureADAuthenticationAsArmPolicyStatusEnabled)
	}

	return &policy
}

func expandReplications(p []any) []replications.Replication {
	reps := make([]replications.Replication, 0)
	if p == nil {
		return reps
	}
	for _, v := range p {
		value := v.(map[string]any)
		location := location.Normalize(value["location"].(string))
		tags := tags.Expand(value["tags"].(map[string]any))
		zoneRedundancy := replications.ZoneRedundancyDisabled
		if value["zone_redundancy_enabled"].(bool) {
			zoneRedundancy = replications.ZoneRedundancyEnabled
		}
		regionEndpointEnabled := value["global_endpoint_routing_enabled"].(bool)
		reps = append(reps, replications.Replication{
			Location: location,
			Name:     &location,
			Tags:     tags,
			Properties: &replications.ReplicationProperties{
				ZoneRedundancy:        &zoneRedundancy,
				RegionEndpointEnabled: pointer.To(regionEndpointEnabled),
			},
		})
	}
	return reps
}

func expandEncryption(input []any) *registries.EncryptionProperty {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]any)
	return &registries.EncryptionProperty{
		KeyVaultProperties: &registries.KeyVaultProperties{
			Identity:      pointer.To(v["identity_client_id"].(string)),
			KeyIdentifier: pointer.To(v["key_vault_key_id"].(string)),
		},
		Status: pointer.To(registries.EncryptionStatusEnabled),
	}
}

func applyGeoReplicationLocations(ctx context.Context, meta any, registryId registries.RegistryId, oldGeoReplications []replications.Replication, newGeoReplications []replications.Replication) error {
	replicationClient := meta.(*clients.Client).Containers.ContainerRegistryClient.Replications

	oldReplications := map[string]replications.Replication{}
	for _, replication := range oldGeoReplications {
		loc := location.Normalize(replication.Location)
		oldReplications[loc] = replication
	}

	newReplications := map[string]replications.Replication{}
	for _, replication := range newGeoReplications {
		loc := location.Normalize(replication.Location)
		newReplications[loc] = replication
	}

	// Delete replications that only appear in the old locations.
	for loc := range oldReplications {
		if _, ok := newReplications[loc]; ok {
			continue
		}
		id := replications.NewReplicationID(registryId.SubscriptionId, registryId.ResourceGroupName, registryId.RegistryName, loc)
		if err := replicationClient.DeleteThenPoll(ctx, id); err != nil {
			return fmt.Errorf("deleting %s: %+v", id, err)
		}
	}

	// Create replications that only exists in the new locations.
	for loc, repl := range newReplications {
		if _, ok := oldReplications[loc]; ok {
			continue
		}
		id := replications.NewReplicationID(registryId.SubscriptionId, registryId.ResourceGroupName, registryId.RegistryName, loc)
		if err := replicationClient.CreateThenPoll(ctx, id, repl); err != nil {
			return fmt.Errorf("creating %s: %+v", id, err)
		}
	}

	// Update (potentially replace) replications that exists at both side.
	for loc, newRepl := range newReplications {
		oldRepl, ok := oldReplications[loc]
		if !ok {
			continue
		}
		// Compare old and new replication parameters to see whether it has updated.
		// If there no update, then skip it. Otherwise, need to check whether the update
		// can happen in place, or need a recreation.

		var (
			needUpdate  bool
			needReplace bool
		)
		// Since the replications here are all derived from expand function, where we guaranteed
		// each properties are non-nil. Whilst we are still doing nil check here in case.
		if oprop, nprop := oldRepl.Properties, newRepl.Properties; oprop != nil && nprop != nil {
			// zoneRedundancy can't be updated in place
			if ov, nv := oprop.ZoneRedundancy, nprop.ZoneRedundancy; ov != nil && nv != nil && *ov != *nv {
				needUpdate = true
				needReplace = true
			}
			if ov, nv := oprop.RegionEndpointEnabled, nprop.RegionEndpointEnabled; ov != nil && nv != nil && *ov != *nv {
				needUpdate = true
			}
		}
		otag, ntag := *oldRepl.Tags, *newRepl.Tags
		if len(otag) != len(ntag) {
			needUpdate = true
		} else {
			for k, ov := range otag {
				nv, ok := ntag[k]
				if !ok {
					needUpdate = true
					break
				}
				if ov != nv {
					needUpdate = true
					break
				}
			}
		}

		if !needUpdate {
			continue
		}

		if needReplace {
			id := replications.NewReplicationID(registryId.SubscriptionId, registryId.ResourceGroupName, registryId.RegistryName, loc)
			if err := replicationClient.DeleteThenPoll(ctx, id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			// Following can be removed once https://github.com/Azure/azure-rest-api-specs/issues/18934 is resolved. Otherwise, the create right after delete will always fail.
			poller := custompollers.NewEventualConsistencyPoller(5, func(pollerCtx context.Context) (*http.Response, error) {
				resp, err := replicationClient.Get(pollerCtx, id)
				return resp.HttpResponse, err
			}, &custompollers.EventualConsistencyPollerOptions{
				Interval:         5 * time.Second,
				TargetStatusCode: pointer.To(http.StatusNotFound),
			})
			if err := poller.PollUntilDone(ctx); err != nil {
				return fmt.Errorf("waiting for deletion of %s: %+v", id, err)
			}
		}

		id := replications.NewReplicationID(registryId.SubscriptionId, registryId.ResourceGroupName, registryId.RegistryName, loc)
		if err := replicationClient.CreateThenPoll(ctx, id, newRepl); err != nil {
			return fmt.Errorf("creating/updating %s: %+v", id, err)
		}
	}

	return nil
}
