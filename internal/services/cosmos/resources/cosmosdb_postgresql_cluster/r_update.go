// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_postgresql_cluster

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresqlhsc/2022-11-08/clusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 3 * time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cosmos.ClustersClient

			id, err := clusters.ParseServerGroupsv2ID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model CosmosDbPostgreSQLClusterModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			parameters := clusters.ClusterForUpdate{
				Properties: &clusters.ClusterPropertiesForUpdate{},
			}

			if metadata.ResourceData.HasChange("administrator_login_password") {
				if model.SourceResourceId == "" && model.AdministratorLoginPassword == "" {
					return fmt.Errorf("`administrator_login_password` is required when `source_resource_id` isn't set")
				}

				parameters.Properties.AdministratorLoginPassword = &model.AdministratorLoginPassword
			}

			if metadata.ResourceData.HasChange("citus_version") {
				parameters.Properties.CitusVersion = &model.CitusVersion
			}

			if metadata.ResourceData.HasChange("coordinator_public_ip_access_enabled") {
				parameters.Properties.CoordinatorEnablePublicIPAccess = &model.CoordinatorPublicIPAccessEnabled
			}

			if metadata.ResourceData.HasChange("coordinator_server_edition") {
				parameters.Properties.CoordinatorServerEdition = &model.CoordinatorServerEdition
			}

			if metadata.ResourceData.HasChange("coordinator_storage_quota_in_mb") {
				if model.SourceResourceId == "" && model.CoordinatorStorageQuotaInMb == 0 {
					return fmt.Errorf("`coordinator_storage_quota_in_mb` is required when `source_resource_id` isn't set")
				}

				parameters.Properties.CoordinatorStorageQuotaInMb = &model.CoordinatorStorageQuotaInMb
			}

			if metadata.ResourceData.HasChange("coordinator_vcore_count") {
				if model.SourceResourceId == "" && model.CoordinatorVCoreCount == 0 {
					return fmt.Errorf("`coordinator_vcore_count` is required when `source_resource_id` isn't set")
				}

				parameters.Properties.CoordinatorVCores = &model.CoordinatorVCoreCount
			}

			if metadata.ResourceData.HasChange("ha_enabled") {
				parameters.Properties.EnableHa = &model.HaEnabled
			}

			if metadata.ResourceData.HasChange("maintenance_window") {
				parameters.Properties.MaintenanceWindow = expandMaintenanceWindow(model.MaintenanceWindow)
			}

			if metadata.ResourceData.HasChange("node_count") {
				parameters.Properties.NodeCount = &model.NodeCount
			}

			if metadata.ResourceData.HasChange("node_public_ip_access_enabled") {
				parameters.Properties.NodeEnablePublicIPAccess = &model.NodePublicIPAccessEnabled
			}

			if metadata.ResourceData.HasChange("node_server_edition") {
				parameters.Properties.NodeServerEdition = &model.NodeServerEdition
			}

			if metadata.ResourceData.HasChange("node_storage_quota_in_mb") {
				parameters.Properties.NodeStorageQuotaInMb = pointer.To(model.NodeStorageQuotaInMb)
			}

			if metadata.ResourceData.HasChange("node_vcores") {
				parameters.Properties.NodeVCores = pointer.To(model.NodeVCores)
			}

			if metadata.ResourceData.HasChange("preferred_primary_zone") {
				parameters.Properties.PreferredPrimaryZone = &model.PreferredPrimaryZone
			}

			if metadata.ResourceData.HasChange("shards_on_coordinator_enabled") {
				parameters.Properties.EnableShardsOnCoordinator = &model.ShardsOnCoordinatorEnabled
			}

			if metadata.ResourceData.HasChange("sql_version") {
				parameters.Properties.PostgresqlVersion = &model.SqlVersion
			}

			if metadata.ResourceData.HasChange("tags") {
				parameters.Tags = &model.Tags
			}

			if err := client.UpdateThenPoll(ctx, *id, parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
