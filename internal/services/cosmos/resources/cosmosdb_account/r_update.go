// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_account

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCosmosDbAccountUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cosmos.CosmosDBClient

	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := cosmosdb.ParseDatabaseAccountID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.DatabaseAccountsGet(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %s", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: model was nil", id)
	}

	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: properties was nil", id)
	}
	props := existing.Model.Properties

	configLocations, err := expandAzureRmCosmosDBAccountGeoLocations(d)
	if err != nil {
		return fmt.Errorf("expanding `geo_location`: %w", err)
	}

	// Normalize Locations...
	cosmosLocations := make([]cosmosdb.Location, 0)
	cosmosLocationsMap := map[string]cosmosdb.Location{}

	for _, l := range pointer.From(existing.Model.Properties.Locations) {
		loc := cosmosdb.Location{
			Id:               l.Id,
			LocationName:     l.LocationName,
			FailoverPriority: l.FailoverPriority,
			IsZoneRedundant:  l.IsZoneRedundant,
		}

		cosmosLocations = append(cosmosLocations, loc)
		cosmosLocationsMap[location.NormalizeNilable(loc.LocationName)] = loc
	}

	// backup must be updated independently
	var backup cosmosdb.BackupPolicy
	if existing.Model.Properties.BackupPolicy != nil {
		backup = existing.Model.Properties.BackupPolicy
		if d.HasChange("backup") {
			if v, ok := d.GetOk("backup"); ok {
				newBackup, err := expandCosmosdbAccountBackup(v.([]any), d.HasChange("backup.0.type"), pointer.FromEnum(existing.Model.Properties.CreateMode))
				if err != nil {
					return fmt.Errorf("expanding `backup`: %+v", err)
				}
				updateParameters := cosmosdb.DatabaseAccountUpdateParameters{
					Properties: &cosmosdb.DatabaseAccountUpdateProperties{
						BackupPolicy: newBackup,
					},
				}

				// Update Database 'backup'...
				if err := client.DatabaseAccountsUpdateThenPoll(ctx, *id, updateParameters); err != nil {
					return fmt.Errorf("updating CosmosDB Account %q (Resource Group %q): %+v", id.DatabaseAccountName, id.ResourceGroupName, err)
				}

				backup = newBackup
			} else if pointer.FromEnum(existing.Model.Properties.CreateMode) != "" {
				return fmt.Errorf("`create_mode` only works when `backup.type` is `Continuous`")
			}
		}
	}

	// NOTE: these fields are expanded directly into the
	// 'DatabaseAccountCreateUpdateParameters' below or
	// are included in the 'DatabaseAccountCreateUpdateParameters'
	// later, however we need to know if they changed or not...
	// lintignore:R019 // deliberate subset: geo_location, identity, default_identity and multi-write settings must be updated in separate atomic calls (see NOTE below)
	updateRequired := d.HasChanges("consistency_policy", "virtual_network_rule", "cors_rule", "access_key_metadata_writes_enabled",
		"network_acl_bypass_for_azure_services", "network_acl_bypass_ids", "analytical_storage",
		"capacity", "restore", "mongo_server_version",
		"public_network_access_enabled", "ip_range_filter", "offer_type", "is_virtual_network_filter_enabled",
		"tags", "automatic_failover_enabled", "analytical_storage_enabled",
		"local_authentication_enabled", "partition_merge_enabled", "minimal_tls_version", "burst_capacity_enabled")

	// Incident : #383341730
	// Azure Bug: #2209567 'Updating identities and default identity at the same time fails silently'
	//
	// The 'Identity' field should only ever be sent once to the endpoint, except for updates and removal. If the
	// 'Identity' field is included in the update call with the 'DefaultIdentity' it will silently fail
	// per the bug noted above (e.g. Azure Bug #2209567).
	//
	// In the update scenario where the end-user would like to update their 'Identity' and their 'DefaultIdentity'
	// fields at the same time both of these operations need to happen atomically in separate PUT/PATCH calls
	// to the service else you will hit the bug mentioned above. You need to update the 'Identity' field
	// first then update the 'DefaultIdentity' in totally different PUT/PATCH calls where you have to drop
	// the 'Identity' field on the floor when updating the 'DefaultIdentity' field.
	//
	// NOTE      : If the 'Identity' field has not changed in the resource, do not send it in the payload.
	//             this workaround can be removed once the service team fixes the above mentioned bug.
	//
	// ADDITIONAL: You cannot update properties and add/remove replication locations or update the enabling of
	//             multiple write locations at the same time. So you must update any changed properties
	//             first, then address the replication locations and/or updating/enabling of
	//             multiple write locations.

	account := cosmosdb.DatabaseAccountCreateUpdateParameters{
		Location: pointer.To(location.NormalizeNilable(existing.Model.Location)),
		Kind:     pointer.ToEnum[cosmosdb.DatabaseAccountKind](d.Get("kind").(string)),
		Properties: cosmosdb.DatabaseAccountCreateUpdateProperties{
			DatabaseAccountOfferType:           cosmosdb.DatabaseAccountOfferType(d.Get("offer_type").(string)),
			IPRules:                            common.CosmosDBIpRangeFilterToIpRules(*pluginsdk.ExpandStringSlice(d.Get("ip_range_filter").(*pluginsdk.Set).List())),
			IsVirtualNetworkFilterEnabled:      pointer.To(d.Get("is_virtual_network_filter_enabled").(bool)),
			EnableFreeTier:                     existing.Model.Properties.EnableFreeTier,
			EnableAutomaticFailover:            pointer.To(d.Get("automatic_failover_enabled").(bool)),
			MinimalTlsVersion:                  pointer.ToEnum[cosmosdb.MinimalTlsVersion](d.Get("minimal_tls_version").(string)),
			Capabilities:                       existing.Model.Properties.Capabilities,
			ConsistencyPolicy:                  expandAzureRmCosmosDBAccountConsistencyPolicy(d),
			KeyVaultKeyUri:                     existing.Model.Properties.KeyVaultKeyUri,
			Locations:                          cosmosLocations,
			VirtualNetworkRules:                expandAzureRmCosmosDBAccountVirtualNetworkRules(d),
			EnableMultipleWriteLocations:       props.EnableMultipleWriteLocations,
			PublicNetworkAccess:                expandCosmosdbAccountPublicNetworkAccess(d.Get("public_network_access_enabled").(bool)),
			EnableAnalyticalStorage:            pointer.To(d.Get("analytical_storage_enabled").(bool)),
			Cors:                               common.ExpandCosmosCorsRule(d.Get("cors_rule").([]any)),
			DisableKeyBasedMetadataWriteAccess: pointer.To(!d.Get("access_key_metadata_writes_enabled").(bool)),
			NetworkAclBypass:                   expandCosmosdbAccountNetworkBypass(d.Get("network_acl_bypass_for_azure_services").(bool)),
			NetworkAclBypassResourceIds:        pluginsdk.ExpandStringSlice(d.Get("network_acl_bypass_ids").([]any)),
			DisableLocalAuth:                   pointer.To(!d.Get("local_authentication_enabled").(bool)),
			BackupPolicy:                       backup,
			EnablePartitionMerge:               pointer.To(d.Get("partition_merge_enabled").(bool)),
			EnableBurstCapacity:                pointer.To(d.Get("burst_capacity_enabled").(bool)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	// 'default_identity_type' will always have a value since it now has a default value of "FirstPartyIdentity" per the API documentation.
	// I do not include 'DefaultIdentity' and 'Identity' in the 'accountProps' intentionally, these operations need to be
	// performed mutually exclusive from each other in an atomic fashion, else you will hit the service teams bug...
	updateDefaultIdentity := d.HasChange("default_identity_type")

	// adding 'DefaultIdentity' to avoid causing it to fallback
	// to "FirstPartyIdentity" on update(s), issue #22466
	if v, ok := d.GetOk("default_identity_type"); ok {
		account.Properties.DefaultIdentity = pointer.To(v.(string))
	}

	// we need the following in the accountProps even if they have not changed...
	if v, ok := d.GetOk("analytical_storage"); ok {
		account.Properties.AnalyticalStorageConfiguration = expandCosmosDBAccountAnalyticalStorageConfiguration(v.([]any))
	}

	if v, ok := d.GetOk("capacity"); ok {
		account.Properties.Capacity = expandCosmosDBAccountCapacity(v.([]any))
	}

	var createMode string
	if v, ok := d.GetOk("create_mode"); ok {
		createMode = v.(string)
		account.Properties.CreateMode = pointer.ToEnum[cosmosdb.CreateMode](createMode)
	}

	if v, ok := d.GetOk("restore"); ok {
		account.Properties.RestoreParameters = expandCosmosdbAccountRestoreParameters(v.([]any))
	}

	if !pluginsdk.IsExplicitlyNullInConfig(d, "mongo_server_version") {
		account.Properties.ApiProperties = &cosmosdb.ApiProperties{
			ServerVersion: pointer.ToEnum[cosmosdb.ServerVersion](d.Get("mongo_server_version").(string)),
		}
	}

	// Only do this update if a value has changed above...
	if updateRequired {
		if err = resourceCosmosDbAccountApiCreateOrUpdate(client, ctx, *id, account); err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}
	}

	// Update the following properties independently after the initial CreateOrUpdate...
	if d.HasChange("multiple_write_locations_enabled") {
		account.Properties.EnableMultipleWriteLocations = pointer.To(d.Get("multiple_write_locations_enabled").(bool))

		// Update the database...
		if err = resourceCosmosDbAccountApiCreateOrUpdate(client, ctx, *id, account); err != nil {
			return fmt.Errorf("updating %q EnableMultipleWriteLocations: %+v", id, err)
		}
	}

	// determine if any locations have been renamed/priority reordered and remove them
	updateLocations := false
	for _, configLoc := range configLocations {
		if cosmosLoc, ok := cosmosLocationsMap[pointer.From(configLoc.LocationName)]; ok {
			// is the location in the config also in the database with the same 'FailoverPriority'?
			if pointer.From(configLoc.FailoverPriority) != pointer.From(cosmosLoc.FailoverPriority) {
				// The Failover Priority has been changed in the config...
				if pointer.From(configLoc.FailoverPriority) == 0 {
					return fmt.Errorf("cannot change the failover priority of %q location %q to %d", id, pointer.From(configLoc.LocationName), pointer.From(configLoc.FailoverPriority))
				}

				// since the Locations FailoverPriority changed remove it from the map because
				// we have to update the Location in the database. The Locations
				// left in the map after this loop are the Locations that are
				// the same in the database and in the config file...
				delete(cosmosLocationsMap, pointer.From(configLoc.LocationName))
				updateLocations = true
			}
		}
	}

	if updateLocations {
		locationsUnchanged := make([]cosmosdb.Location, 0, len(cosmosLocationsMap))
		for _, value := range cosmosLocationsMap {
			locationsUnchanged = append(locationsUnchanged, value)
		}

		account.Properties.Locations = locationsUnchanged

		// Update the database...
		if err = resourceCosmosDbAccountApiCreateOrUpdate(client, ctx, *id, account); err != nil {
			return fmt.Errorf("removing %q renamed `locations`: %+v", id, err)
		}
	}

	if d.HasChanges("geo_location") {
		account.Properties.Locations = configLocations

		// Update the database locations...
		if err = resourceCosmosDbAccountApiCreateOrUpdate(client, ctx, *id, account); err != nil {
			return fmt.Errorf("updating %q `locations`: %+v", id, err)
		}
	}

	// Update Identity and Default Identity...
	identityChanged := false
	if d.HasChange("identity") {
		identityChanged = true

		// Looks like you have to always remove all the identities first before you can
		// reassign/modify them, else it will append any new/changed identities
		// resulting in a diff...

		// can't set this back to account, because that will hit the bug...
		identityVal := cosmosdb.DatabaseAccountUpdateParameters{
			Identity: pointer.To(identity.LegacySystemAndUserAssignedMap{
				Type: identity.TypeNone,
			}),
		}

		// Update the database 'Identity' to 'None'...
		if err = resourceCosmosDbAccountApiUpdate(client, ctx, *id, identityVal); err != nil {
			return fmt.Errorf("updating 'identity' %q: %+v", id, err)
		}

		expandedIdentity, err := identity.ExpandLegacySystemAndUserAssignedMap(d.Get("identity").([]any))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}

		// If the Identity was removed from the configuration file it will be set as type None
		// so we can skip setting the Identity if it is going to be set to None...
		if expandedIdentity.Type != identity.TypeNone {
			identityVal := cosmosdb.DatabaseAccountUpdateParameters{
				Identity: expandedIdentity,
			}

			// Update the database...
			if err = resourceCosmosDbAccountApiUpdate(client, ctx, *id, identityVal); err != nil {
				return fmt.Errorf("updating `identity` for %s: %+v", id, err)
			}
		}
	}

	// NOTE: updateDefaultIdentity now has a default value of 'FirstPartyIdentity'... This value now gets
	//       triggered if the default value does not match the value in Azure...
	//
	// NOTE: When you change the 'Identity', the 'DefaultIdentity' will be set to 'undefined', so if you change
	//       the identity you must also update the 'DefaultIdentity' as well...
	if updateDefaultIdentity || identityChanged {
		// This will now return the default of 'FirstPartyIdentity' if it
		// is not set in the config, which is correct.

		// PATCH instead of PUT...
		defaultIdentity := cosmosdb.DatabaseAccountUpdateParameters{
			Properties: &cosmosdb.DatabaseAccountUpdateProperties{
				DefaultIdentity: pointer.To(d.Get("default_identity_type").(string)),
			},
		}

		// Update the database...
		if err = resourceCosmosDbAccountApiUpdate(client, ctx, *id, defaultIdentity); err != nil {
			return fmt.Errorf("updating `default_identity_type` for %s: %+v", id, err)
		}
	}

	if existing.Model.Properties.Capabilities != nil {
		if d.HasChange("capabilities") {
			updateParameters := cosmosdb.DatabaseAccountUpdateParameters{
				Properties: &cosmosdb.DatabaseAccountUpdateProperties{
					Capabilities: expandAzureRmCosmosDBAccountCapabilities(d),
				},
			}

			if err := client.DatabaseAccountsUpdateThenPoll(ctx, *id, updateParameters); err != nil {
				return fmt.Errorf("updating `capabilities` for %s: %+v", id, err)
			}
		}
	}

	return resourceCosmosDbAccountRead(d, meta)
}

func resourceCosmosDbAccountApiUpdate(client *cosmosdb.CosmosDBClient, ctx context.Context, id cosmosdb.DatabaseAccountId, account cosmosdb.DatabaseAccountUpdateParameters) error {
	if err := client.DatabaseAccountsUpdateThenPoll(ctx, id, account); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   []string{"Creating", "Updating", "Deleting", "Initializing", "Dequeued", "Enqueued"},
		Target:                    []string{"Succeeded"},
		MinTimeout:                15 * time.Second,
		ContinuousTargetOccurence: 2,
		Refresh: func() (any, string, error) {
			resp, err2 := client.DatabaseAccountsGet(ctx, id)
			if err2 != nil || resp.HttpResponse == nil || resp.HttpResponse.StatusCode == http.StatusNotFound {
				return nil, "", fmt.Errorf("retrieving %s: %+v", id, err2)
			}
			status := "Succeeded"

			return resp, status, nil
		},
	}

	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("internal-error: context had no deadline")
	}
	stateConf.Timeout = time.Until(deadline)

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to update: %+v", id, err)
	}

	return nil
}

func resourceCosmosDbAccountApiCreateOrUpdate(client *cosmosdb.CosmosDBClient, ctx context.Context, id cosmosdb.DatabaseAccountId, account cosmosdb.DatabaseAccountCreateUpdateParameters) error {
	if err := client.DatabaseAccountsCreateOrUpdateThenPoll(ctx, id, account); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	// if a replication location is added or removed it can take some time to provision
	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{"Creating", "Updating", "Deleting", "Initializing", "Dequeued", "Enqueued"},
		Target:     []string{"Succeeded"},
		MinTimeout: 15 * time.Second,
		Delay:      30 * time.Second, // required because it takes some time before the 'creating' location shows up
		Refresh: func() (any, string, error) {
			resp, err2 := client.DatabaseAccountsGet(ctx, id)
			if err2 != nil || resp.HttpResponse == nil || resp.HttpResponse.StatusCode == http.StatusNotFound {
				return nil, "", fmt.Errorf("retrieving %s: %+v", id, err2)
			}
			status := "Succeeded"
			if props := resp.Model.Properties; props != nil {
				var locations []cosmosdb.Location

				if props.ReadLocations != nil {
					locations = append(locations, *props.ReadLocations...)
				}
				if props.WriteLocations != nil {
					locations = append(locations, *props.WriteLocations...)
				}

				for _, l := range locations {
					if status = *l.ProvisioningState; status == "Creating" || status == "Updating" || status == "Deleting" {
						break // return the first non successful status.
					}
				}

				for _, desiredLocation := range account.Properties.Locations {
					for index, l := range locations {
						if location.Normalize(*desiredLocation.LocationName) == location.Normalize(*l.LocationName) {
							break
						}

						if (index + 1) == len(locations) {
							return resp, "Updating", nil
						}
					}
				}
			}

			return resp, status, nil
		},
	}

	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("internal-error: context had no deadline")
	}
	stateConf.Timeout = time.Until(deadline)

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to provision: %+v", id, err)
	}

	return nil
}
