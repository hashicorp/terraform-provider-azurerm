// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mongo_cluster

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mongocluster/2026-06-01/mongoclusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MongoCluster.MongoClustersClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state MongoClusterResourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := mongoclusters.NewMongoClusterID(subscriptionId, state.ResourceGroupName, state.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
					}
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			parameter := mongoclusters.MongoCluster{
				Location: location.Normalize(state.Location),
				Properties: &mongoclusters.MongoClusterProperties{
					AuthConfig:        expandMongoClusterAuthConfig(state.AuthenticationMethods),
					Encryption:        expandMongoClusterCustomerManagedKey(state.CustomerManagedKey),
					RestoreParameters: expandMongoClusterRestore(state.Restore),
				},
			}

			identityVal, err := identity.ExpandUserAssignedMapFromModel(state.Identity)
			if err != nil {
				return fmt.Errorf(`expanding "identity": %v`, err)
			}
			// Per the current service API design, they don’t allow setting `userAssignedIdentities` to `nil` in the request payload when the `type` of `identity` is `nil`; otherwise, the API would return an error.
			// Therefore, we have to use the customized function instead of the common one, since the common function always sets `userAssignedIdentities` to `nil` in the request payload.
			// Service team confirmed that it will be more flexible, and we will allow `userAssignedIdentities = nil` in the future. Tracking issue: https://github.com/Azure/azure-rest-api-specs/issues/38575
			if identityVal != nil && identityVal.Type == identity.TypeNone {
				identityVal = nil
			}
			parameter.Identity = identityVal

			if state.AdministratorUserName != "" {
				parameter.Properties.Administrator = &mongoclusters.AdministratorProperties{
					UserName: pointer.To(state.AdministratorUserName),
					Password: pointer.To(state.AdministratorPassword),
				}
			}

			if state.CreateMode != "" {
				parameter.Properties.CreateMode = pointer.ToEnum[mongoclusters.CreateMode](state.CreateMode)
			}

			parameter.Properties.PreviewFeatures = expandPreviewFeatures(state.PreviewFeatures)

			if state.ShardCount != 0 {
				parameter.Properties.Sharding = &mongoclusters.ShardingProperties{
					ShardCount: pointer.To(state.ShardCount),
				}
			}

			if state.CreateMode == string(mongoclusters.CreateModeGeoReplica) {
				if state.SourceServerId == "" {
					return fmt.Errorf("`source_server_id` is required when `create_mode` is `GeoReplica`")
				}

				parameter.Properties.ReplicaParameters = &mongoclusters.MongoClusterReplicaParameters{
					SourceLocation:   state.SourceLocation,
					SourceResourceId: state.SourceServerId,
				}
			}

			if state.ComputeTier != "" {
				parameter.Properties.Compute = &mongoclusters.ComputeProperties{
					Tier: pointer.To(state.ComputeTier),
				}
			}

			if state.HighAvailabilityMode != "" {
				parameter.Properties.HighAvailability = &mongoclusters.HighAvailabilityProperties{
					TargetMode: pointer.ToEnum[mongoclusters.HighAvailabilityMode](state.HighAvailabilityMode),
				}
			}

			parameter.Properties.PublicNetworkAccess = pointer.ToEnum[mongoclusters.PublicNetworkAccess](state.PublicNetworkAccess)

			if state.StorageSizeInGb != 0 {
				parameter.Properties.Storage = &mongoclusters.StorageProperties{
					SizeGb: pointer.To(state.StorageSizeInGb),
					Type:   pointer.ToEnum[mongoclusters.StorageType](state.StorageType),
				}
			}

			if state.Tags != nil {
				parameter.Tags = pointer.To(state.Tags)
			}

			if state.Version != "" {
				parameter.Properties.ServerVersion = pointer.To(state.Version)
			}

			// Always set `data_api_mode_enabled` to `false` during creation for `Default` mode
			// Update it separately if user wants it enabled
			if state.CreateMode == string(mongoclusters.CreateModeDefault) {
				parameter.Properties.DataApi = &mongoclusters.DataApiProperties{
					Mode: pointer.To(mongoclusters.DataApiModeDisabled),
				}
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, parameter, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
			metadata.SetID(id)

			// `networkBypassMode` can only be set by a separate PATCH request after the cluster is created, so enable it after the cluster exists.
			if state.CosmosDBNetworkBypassEnabled {
				if err := updateMongoClusterCosmosDBNetworkBypassEnabled(ctx, client, id, true); err != nil {
					return fmt.Errorf("enabling `cosmos_db_network_bypass_enabled` for %s: %+v", id, err)
				}
			}

			// `data_api_mode_enabled` can only be enabled after the resource is created
			if state.CreateMode == string(mongoclusters.CreateModeDefault) && state.DataApiModeEnabled {
				updatePayload := mongoclusters.MongoClusterUpdate{
					Properties: &mongoclusters.MongoClusterUpdateProperties{
						DataApi: &mongoclusters.DataApiProperties{
							Mode: pointer.To(mongoclusters.DataApiModeEnabled),
						},
					},
				}

				if err := client.UpdateThenPoll(ctx, id, updatePayload); err != nil {
					return fmt.Errorf("updating `data_api_mode_enabled`: %+v", err)
				}
			}

			return nil
		},
	}
}

func expandPreviewFeatures(input []string) *[]mongoclusters.PreviewFeature {
	if len(input) == 0 {
		return nil
	}

	result := make([]mongoclusters.PreviewFeature, 0)

	for _, v := range input {
		if v != "" {
			result = append(result, mongoclusters.PreviewFeature(v))
		}
	}

	return &result
}

func expandMongoClusterCustomerManagedKey(input []CustomerManagedKey) *mongoclusters.EncryptionProperties {
	if len(input) == 0 {
		return nil
	}

	cmk := input[0]

	return &mongoclusters.EncryptionProperties{
		CustomerManagedKeyEncryption: &mongoclusters.CustomerManagedKeyEncryptionProperties{
			KeyEncryptionKeyURL: pointer.To(cmk.KeyVaultKeyId),
			KeyEncryptionKeyIdentity: &mongoclusters.KeyEncryptionKeyIdentity{
				IdentityType:                   pointer.To(mongoclusters.KeyEncryptionKeyIdentityTypeUserAssignedIdentity),
				UserAssignedIdentityResourceId: pointer.To(cmk.UserAssignedIdentityId),
			},
		},
	}
}

func expandMongoClusterRestore(input []Restore) *mongoclusters.MongoClusterRestoreParameters {
	if len(input) == 0 {
		return nil
	}

	restoreParams := input[0]

	return &mongoclusters.MongoClusterRestoreParameters{
		PointInTimeUTC:   pointer.To(restoreParams.PointInTimeUtc),
		SourceResourceId: pointer.To(restoreParams.SourceId),
	}
}

func expandMongoClusterAuthConfig(input []string) *mongoclusters.AuthConfigProperties {
	if len(input) == 0 {
		return nil
	}

	allowedModes := make([]mongoclusters.AuthenticationMode, 0)
	for _, v := range input {
		allowedModes = append(allowedModes, mongoclusters.AuthenticationMode(v))
	}

	return &mongoclusters.AuthConfigProperties{
		AllowedModes: pointer.To(allowedModes),
	}
}
