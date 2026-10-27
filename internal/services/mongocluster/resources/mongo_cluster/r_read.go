// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mongo_cluster

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mongocluster/2026-06-01/mongoclusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MongoCluster.MongoClustersClient

			id, err := mongoclusters.ParseMongoClusterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := MongoClusterResourceModel{
				Name:              id.MongoClusterName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)

				identity, err := identity.FlattenUserAssignedMapToModel(model.Identity)
				if err != nil {
					return fmt.Errorf("flattening `identity`: %+v", err)
				}
				state.Identity = pointer.From(identity)

				if props := model.Properties; props != nil {
					// API doesn't return the value of administrator_password
					state.AdministratorPassword = metadata.ResourceData.Get("administrator_password").(string)

					// API doesn't return the value of create_mode, https://github.com/Azure/azure-rest-api-specs/issues/31266 has been filed to track it.
					state.CreateMode = metadata.ResourceData.Get("create_mode").(string)

					state.CustomerManagedKey = flattenMongoClusterCustomerManagedKey(props.Encryption)

					if v := props.Administrator; v != nil {
						state.AdministratorUserName = pointer.From(v.UserName)
					}

					if v := props.Replica; v != nil {
						// API doesn't return the value of source_location, https://github.com/Azure/azure-rest-api-specs/issues/31266 has been filed to track it.
						state.SourceLocation = metadata.ResourceData.Get("source_location").(string)
						if v.SourceResourceId != nil {
							id, err := mongoclusters.ParseMongoClusterID(pointer.From(v.SourceResourceId))
							if err != nil {
								return err
							}
							state.SourceServerId = id.ID()
						}
					}

					state.Restore = flattenMongoClusterRestore(props.RestoreParameters)

					if v := props.Sharding; v != nil {
						state.ShardCount = pointer.From(v.ShardCount)
					}
					if v := props.Compute; v != nil {
						state.ComputeTier = pointer.From(v.Tier)
					}

					if v := props.HighAvailability; v != nil {
						state.HighAvailabilityMode = pointer.FromEnum(v.TargetMode)
					}
					state.CosmosDBNetworkBypassEnabled = pointer.From(props.NetworkBypassMode) == mongoclusters.NetworkBypassModeAzureCosmosDB
					state.PublicNetworkAccess = pointer.FromEnum(props.PublicNetworkAccess)

					if v := props.Storage; v != nil {
						state.StorageSizeInGb = pointer.From(v.SizeGb)
						state.StorageType = pointer.FromEnum(v.Type)
					}
					if v := props.PreviewFeatures; v != nil {
						state.PreviewFeatures = flattenMongoClusterPreviewFeatures(v)
					}
					state.Version = pointer.From(props.ServerVersion)

					state.AuthenticationMethods = flattenMongoClusterAuthConfig(props.AuthConfig)

					if v := props.DataApi; v != nil {
						state.DataApiModeEnabled = pointer.From(v.Mode) == mongoclusters.DataApiModeEnabled
					}
				}

				state.Tags = pointer.From(model.Tags)
			}

			csResp, err := client.ListConnectionStrings(ctx, *id)
			if err != nil {
				return fmt.Errorf("listing connection strings for %s: %+v", *id, err)
			}
			if model := csResp.Model; model != nil {
				state.ConnectionStrings = flattenMongoClusterConnectionStrings(model.ConnectionStrings, state.AdministratorUserName, state.AdministratorPassword)
			}

			return metadata.Encode(&state)
		},
	}
}

type MongoClusterConnectionString struct {
	Value       string `tfschema:"value"`
	Description string `tfschema:"description"`
	Name        string `tfschema:"name"`
}

type CustomerManagedKey struct {
	KeyVaultKeyId          string `tfschema:"key_vault_key_id"`
	UserAssignedIdentityId string `tfschema:"user_assigned_identity_id"`
}

type Restore struct {
	PointInTimeUtc string `tfschema:"point_in_time_utc"`
	SourceId       string `tfschema:"source_id"`
}

func validateMongoClusterCosmosDBNetworkBypassEnabled(metadata sdk.ResourceMetaData, state MongoClusterResourceModel) error {
	rawConfig := metadata.ResourceDiff.GetRawConfig()
	if !rawConfig.IsKnown() || rawConfig.IsNull() {
		return nil
	}

	rawConfigMap := rawConfig.AsValueMap()
	cosmosDBNetworkBypassEnabled := rawConfigMap["cosmos_db_network_bypass_enabled"]
	if !cosmosDBNetworkBypassEnabled.IsKnown() || cosmosDBNetworkBypassEnabled.IsNull() || !state.CosmosDBNetworkBypassEnabled {
		return nil
	}

	if !rawConfigMap["public_network_access"].IsKnown() || !rawConfigMap["authentication_methods"].IsKnown() {
		return nil
	}

	if state.PublicNetworkAccess != string(mongoclusters.PublicNetworkAccessDisabled) {
		return errors.New("`public_network_access` must be `Disabled` when `cosmos_db_network_bypass_enabled` is `true`")
	}

	if len(state.AuthenticationMethods) != 1 || state.AuthenticationMethods[0] != string(mongoclusters.AuthenticationModeMicrosoftEntraID) {
		return errors.New("`authentication_methods` must contain only `MicrosoftEntraID` when `cosmos_db_network_bypass_enabled` is `true`")
	}

	return nil
}

// updateMongoClusterCosmosDBNetworkBypassEnabled updates the `networkBypassMode` of a MongoDB cluster to enable or disable the Cosmos DB network bypass feature. This field can only be modified via a separate PATCH request after the cluster is created.
func updateMongoClusterCosmosDBNetworkBypassEnabled(ctx context.Context, client *mongoclusters.MongoClustersClient, id mongoclusters.MongoClusterId, enabled bool) error {
	mode := mongoclusters.NetworkBypassModeNone
	if enabled {
		mode = mongoclusters.NetworkBypassModeAzureCosmosDB
	}

	payload := mongoclusters.MongoClusterUpdate{
		Properties: &mongoclusters.MongoClusterUpdateProperties{
			NetworkBypassMode: pointer.To(mode),
		},
	}

	return client.UpdateThenPoll(ctx, id, payload)
}

func isNativeAuthRequired(metadata sdk.ResourceMetaData) bool {
	authMethodsRaw := metadata.ResourceDiff.GetRawConfig().AsValueMap()["authentication_methods"]
	if !authMethodsRaw.IsKnown() {
		return false
	}
	if authMethodsRaw.IsNull() {
		return true
	}

	authMethods := authMethodsRaw.AsValueSet().Values()
	if len(authMethods) == 0 {
		return true
	}

	for _, v := range authMethods {
		if !v.IsKnown() {
			continue
		}
		if v.AsString() == string(mongoclusters.AuthenticationModeNativeAuth) {
			return true
		}
	}

	return false
}
