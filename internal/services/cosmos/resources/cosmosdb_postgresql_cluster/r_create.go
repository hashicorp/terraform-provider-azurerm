// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_postgresql_cluster

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresqlhsc/2022-11-08/clusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 3 * time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model CosmosDbPostgreSQLClusterModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Cosmos.ClustersClient
			subscriptionId := metadata.Client.Account.SubscriptionId
			id := clusters.NewServerGroupsv2ID(subscriptionId, model.ResourceGroupName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			parameters := &clusters.Cluster{
				Location: location.Normalize(model.Location),
				Properties: &clusters.ClusterProperties{
					AdministratorLoginPassword:      &model.AdministratorLoginPassword,
					CoordinatorEnablePublicIPAccess: &model.CoordinatorPublicIPAccessEnabled,
					CoordinatorServerEdition:        &model.CoordinatorServerEdition,
					CoordinatorStorageQuotaInMb:     &model.CoordinatorStorageQuotaInMb,
					CoordinatorVCores:               &model.CoordinatorVCoreCount,
					EnableHa:                        &model.HaEnabled,
					NodeCount:                       &model.NodeCount,
					NodeEnablePublicIPAccess:        &model.NodePublicIPAccessEnabled,
					NodeServerEdition:               &model.NodeServerEdition,
				},
			}

			if v := model.CitusVersion; v != "" {
				parameters.Properties.CitusVersion = &model.CitusVersion
			}

			if v := model.MaintenanceWindow; v != nil {
				parameters.Properties.MaintenanceWindow = expandMaintenanceWindow(v)
			}

			if v := model.NodeStorageQuotaInMb; v != 0 {
				parameters.Properties.NodeStorageQuotaInMb = pointer.To(model.NodeStorageQuotaInMb)
			}

			if v := model.NodeVCores; v != 0 {
				parameters.Properties.NodeVCores = pointer.To(model.NodeVCores)
			}

			if v := model.PointInTimeInUTC; v != "" {
				parameters.Properties.PointInTimeUTC = &model.PointInTimeInUTC
			}

			if v := model.SqlVersion; v != "" {
				parameters.Properties.PostgresqlVersion = &model.SqlVersion
			}

			if v := model.PreferredPrimaryZone; v != "" {
				parameters.Properties.PreferredPrimaryZone = &model.PreferredPrimaryZone
			}

			if v := model.SourceLocation; v != "" {
				parameters.Properties.SourceLocation = &model.SourceLocation
			}

			switch {
			case model.SourceResourceId != "":
				parameters.Properties.SourceResourceId = &model.SourceResourceId
			case model.AdministratorLoginPassword == "":
				return fmt.Errorf("`administrator_login_password` is required when `source_resource_id` isn't set")
			case model.CoordinatorStorageQuotaInMb == 0:
				return fmt.Errorf("`coordinator_storage_quota_in_mb` is required when `source_resource_id` isn't set")
			case model.CoordinatorVCoreCount == 0:
				return fmt.Errorf("`coordinator_vcore_count` is required when `source_resource_id` isn't set")
			}

			// If `shards_on_coordinator_enabled` isn't set, API would set it to `true` when `node_count` is `0`.
			// If `shards_on_coordinator_enabled` isn't set, API would set it to `false` when `node_count` is greater than or equal to `2`.
			// As `shards_on_coordinator_enabled` is `bool` and it's always set to `false` as zero value when it isn't set, so we cannot use `model.ShardsOnCoordinatorEnabled` to check if this property is set in tf config.
			//nolint:staticcheck,tfproviderlint // XR001: GetOkExists is needed here, see above
			if v, ok := metadata.ResourceData.GetOkExists("shards_on_coordinator_enabled"); ok {
				parameters.Properties.EnableShardsOnCoordinator = pointer.To(v.(bool))
			}

			if v := model.Tags; v != nil {
				parameters.Tags = &v
			}

			if err := client.CreateCallbackThenPoll(ctx, id, *parameters, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id)
		},
	}
}

func expandMaintenanceWindow(input []MaintenanceWindow) *clusters.MaintenanceWindow {
	if len(input) == 0 {
		return &clusters.MaintenanceWindow{
			CustomWindow: pointer.To("Disabled"),
		}
	}

	v := input[0]

	maintenanceWindow := clusters.MaintenanceWindow{
		CustomWindow: pointer.To("Enabled"),
		StartHour:    pointer.To(v.StartHour),
		StartMinute:  pointer.To(v.StartMinute),
		DayOfWeek:    pointer.To(v.DayOfWeek),
	}

	return &maintenanceWindow
}
