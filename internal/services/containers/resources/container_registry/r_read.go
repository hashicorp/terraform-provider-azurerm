// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_registry

import (
	"fmt"
	"log"
	"sort"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2025-11-01/registries"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2025-11-01/replications"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceContainerRegistryRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Containers.ContainerRegistryClient.Registries
	replicationClient := meta.(*clients.Client).Containers.ContainerRegistryClient.Replications
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := registries.ParseRegistryID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.RegistryName)
	d.Set("resource_group_name", id.ResourceGroupName)

	// this must be set to filter out the georeplication for the container registry's current location
	loc := ""

	if model := resp.Model; model != nil {
		loc = location.Normalize(model.Location)
		d.Set("location", loc)

		identity, err := identity.FlattenSystemAndUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		if err := d.Set("identity", identity); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		d.Set("sku", pointer.FromEnum(model.Sku.Tier))

		if props := model.Properties; props != nil {
			d.Set("admin_enabled", props.AdminUserEnabled)
			d.Set("login_server", props.LoginServer)
			d.Set("public_network_access_enabled", *props.PublicNetworkAccess == registries.PublicNetworkAccessEnabled)

			if err := d.Set("network_rule_set", flattenNetworkRuleSet(props.NetworkRuleSet)); err != nil {
				return fmt.Errorf("setting `network_rule_set`: %+v", err)
			}

			if err := d.Set("encryption", flattenEncryption(props.Encryption)); err != nil {
				return fmt.Errorf("setting `encryption`: %+v", err)
			}
			d.Set("zone_redundancy_enabled", *props.ZoneRedundancy == registries.ZoneRedundancyEnabled)
			d.Set("anonymous_pull_enabled", props.AnonymousPullEnabled)
			d.Set("data_endpoint_enabled", props.DataEndpointEnabled)
			d.Set("data_endpoint_host_names", props.DataEndpointHostNames)
			d.Set("network_rule_bypass_option", pointer.FromEnum(props.NetworkRuleBypassOptions))
			d.Set("role_assignment_mode", pointer.FromEnum(props.RoleAssignmentMode))
			d.Set("network_rule_bypass_for_tasks_enabled", pointer.From(props.NetworkRuleBypassAllowedForTasks))

			if policies := props.Policies; policies != nil {
				var retentionInDays int64
				if policies.RetentionPolicy != nil && policies.RetentionPolicy.Status != nil && *policies.RetentionPolicy.Status == registries.PolicyStatusEnabled {
					retentionInDays = pointer.From(policies.RetentionPolicy.Days)
				}
				d.Set("retention_policy_in_days", retentionInDays)

				d.Set("quarantine_policy_enabled", flattenQuarantinePolicy(props.Policies))
				d.Set("export_policy_enabled", flattenExportPolicy(props.Policies))
				d.Set("azuread_authentication_as_arm_policy_enabled", flattenAadAuthAsArmPolicy(props.Policies))
			}

			if *props.AdminUserEnabled {
				credsResp, errList := client.ListCredentials(ctx, *id)
				if errList != nil {
					return fmt.Errorf("retrieving credentials for %s: %s", *id, errList)
				}

				if credsModel := credsResp.Model; credsModel != nil {
					d.Set("admin_username", credsModel.Username)
					for _, v := range *credsModel.Passwords {
						d.Set("admin_password", v.Value)
						break
					}
				}
			} else {
				d.Set("admin_username", "")
				d.Set("admin_password", "")
			}
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return fmt.Errorf("flattening `tags`: %+v", err)
		}
	}

	rId, err := replications.ParseRegistryID(id.ID())
	if err != nil {
		return err
	}
	replicationsResp, err := replicationClient.List(ctx, *rId)
	if err != nil {
		return fmt.Errorf("retrieving replications for %s: %s", *id, err)
	}

	geoReplications := make([]any, 0)
	if replicationsModel := replicationsResp.Model; replicationsModel != nil {
		for _, value := range *replicationsModel {
			valueLocation := location.Normalize(value.Location)
			if valueLocation != loc {
				replication := make(map[string]any)
				replication["location"] = valueLocation
				replication["tags"] = tags.Flatten(value.Tags)
				replication["zone_redundancy_enabled"] = *value.Properties.ZoneRedundancy == replications.ZoneRedundancyEnabled
				replication["global_endpoint_routing_enabled"] = value.Properties.RegionEndpointEnabled != nil && *value.Properties.RegionEndpointEnabled
				geoReplications = append(geoReplications, replication)
			}
		}
	}

	// The order of the georeplications returned from the list API is not consistent. We simply order it alphabetically to be consistent.
	sort.Slice(geoReplications, func(i, j int) bool {
		return geoReplications[i].(map[string]any)["location"].(string) < geoReplications[j].(map[string]any)["location"].(string)
	})

	d.Set("georeplications", geoReplications)

	return nil
}

func applyContainerRegistrySku(d *pluginsdk.ResourceData, meta any, sku string, id registries.RegistryId) error {
	client := meta.(*clients.Client).Containers.ContainerRegistryClient.Registries
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	parameters := registries.RegistryUpdateParameters{
		Sku: &registries.Sku{
			Name: registries.SkuName(sku),
			Tier: pointer.ToEnum[registries.SkuTier](sku),
		},
	}

	if err := client.UpdateThenPoll(ctx, id, parameters); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return nil
}
