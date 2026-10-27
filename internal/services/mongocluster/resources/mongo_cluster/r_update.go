// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mongo_cluster

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mongocluster/2026-06-01/mongoclusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MongoCluster.MongoClustersClient

			id, err := mongoclusters.ParseMongoClusterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state MongoClusterResourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}
			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", *id)
			}
			payload := existing.Model

			// Set SystemData to nil as the API returns `The property '#/systemData' of type null did not match the following type: object in schema 25debcc2-6915-5536-9566-a2ecd765b755"}}` error.
			// https://github.com/Azure/azure-rest-api-specs/issues/31377 has been filed to track it.
			payload.SystemData = nil

			// Set `identity` to `nil` to avoid schema validation errors returned by service API when `identity` type is `None`.
			if payload.Identity != nil && payload.Identity.Type == identity.TypeNone {
				// Setting the `identity` to nil is equivalent to setting the `identity` to `None`.
				payload.Identity = nil
			}

			// Clear `data_api_mode_enabled` for clusters created with `PointInTimeRestore` or `GeoReplica` mode to avoid schema validation errors
			// `data_api_mode_enabled` is only supported for `Default` create mode
			if state.CreateMode != string(mongoclusters.CreateModeDefault) {
				payload.Properties.DataApi = nil
			}

			cosmosDBNetworkBypassEnabledChanged := metadata.ResourceData.HasChange("cosmos_db_network_bypass_enabled")
			if cosmosDBNetworkBypassEnabledChanged {
				// Azure blocks public network access changes while bypass is enabled, so disable bypass before the main update if it's intended to be updated to `false`.
				if !state.CosmosDBNetworkBypassEnabled {
					if err := updateMongoClusterCosmosDBNetworkBypassEnabled(ctx, client, *id, false); err != nil {
						return fmt.Errorf("disabling `cosmos_db_network_bypass_enabled` for %s: %+v", *id, err)
					}
				}

				payload.Properties.NetworkBypassMode = pointer.To(mongoclusters.NetworkBypassModeNone)
			}

			// upgrades involving Free or M25(Burstable) compute tier require first upgrading the compute tier, after which other configurations can be updated.
			if metadata.ResourceData.HasChange("compute_tier") {
				payload.Properties.Compute = &mongoclusters.ComputeProperties{
					Tier: pointer.To(state.ComputeTier),
				}
				oldComputeTier, newComputeTier := metadata.ResourceData.GetChange("compute_tier")
				if (oldComputeTier == "Free" || oldComputeTier == "M25") && newComputeTier != "Free" && newComputeTier != "M25" {
					if err := client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
						return fmt.Errorf("updating %s: %+v", *id, err)
					}
				}
			}

			if metadata.ResourceData.HasChange("administrator_password") {
				payload.Properties.Administrator = &mongoclusters.AdministratorProperties{
					UserName: pointer.To(state.AdministratorUserName),
					Password: pointer.To(state.AdministratorPassword),
				}
			}

			if metadata.ResourceData.HasChange("high_availability_mode") {
				payload.Properties.HighAvailability = &mongoclusters.HighAvailabilityProperties{
					TargetMode: pointer.ToEnum[mongoclusters.HighAvailabilityMode](state.HighAvailabilityMode),
				}
			}

			if metadata.ResourceData.HasChange("public_network_access") {
				payload.Properties.PublicNetworkAccess = pointer.ToEnum[mongoclusters.PublicNetworkAccess](state.PublicNetworkAccess)
			}

			if metadata.ResourceData.HasChange("storage_size_in_gb") {
				payload.Properties.Storage = &mongoclusters.StorageProperties{
					SizeGb: pointer.To(state.StorageSizeInGb),
					Type:   pointer.ToEnum[mongoclusters.StorageType](state.StorageType),
				}
			}

			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = pointer.To(state.Tags)
			}

			if metadata.ResourceData.HasChange("version") {
				payload.Properties.ServerVersion = pointer.To(state.Version)
			}

			if metadata.ResourceData.HasChange("customer_managed_key") {
				payload.Properties.Encryption = expandMongoClusterCustomerManagedKey(state.CustomerManagedKey)
			}

			if metadata.ResourceData.HasChange("authentication_methods") {
				payload.Properties.AuthConfig = expandMongoClusterAuthConfig(state.AuthenticationMethods)
			}

			// Only allow enabling `data_api_mode_enabled`, disabling should trigger ForceNew (handled in CustomizeDiff)
			if metadata.ResourceData.HasChange("data_api_mode_enabled") {
				// CustomizeDiff ensures this change is always false -> true
				payload.Properties.DataApi = &mongoclusters.DataApiProperties{
					Mode: pointer.To(mongoclusters.DataApiModeEnabled),
				}
			}

			if metadata.ResourceData.HasChange("identity") {
				identityVal, err := identity.ExpandUserAssignedMapFromModel(state.Identity)
				if err != nil {
					return fmt.Errorf(`expanding "identity": %v`, err)
				}
				payload.Identity = identityVal
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			// Azure requires disabled public access and Entra-only authentication before bypass can be enabled, so we must ensure these conditions are met before attempting to enable it.
			if cosmosDBNetworkBypassEnabledChanged && state.CosmosDBNetworkBypassEnabled {
				if err := updateMongoClusterCosmosDBNetworkBypassEnabled(ctx, client, *id, true); err != nil {
					return fmt.Errorf("enabling `cosmos_db_network_bypass_enabled` for %s: %+v", *id, err)
				}
			}

			return nil
		},
	}
}
