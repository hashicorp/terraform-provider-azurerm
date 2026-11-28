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

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cosmos.ClustersClient

			id, err := clusters.ParseServerGroupsv2ID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			state := CosmosDbPostgreSQLClusterModel{
				Name:              id.ServerGroupsv2Name,
				ResourceGroupName: id.ResourceGroupName,
				Location:          location.Normalize(model.Location),
			}

			if props := model.Properties; props != nil {
				state.AdministratorLoginPassword = metadata.ResourceData.Get("administrator_login_password").(string)
				state.SourceResourceId = metadata.ResourceData.Get("source_resource_id").(string)
				state.SourceLocation = metadata.ResourceData.Get("source_location").(string)
				state.PointInTimeInUTC = metadata.ResourceData.Get("point_in_time_in_utc").(string)
				state.CoordinatorPublicIPAccessEnabled = pointer.From(props.CoordinatorEnablePublicIPAccess)
				state.CoordinatorServerEdition = pointer.From(props.CoordinatorServerEdition)
				state.CoordinatorStorageQuotaInMb = pointer.From(props.CoordinatorStorageQuotaInMb)
				state.CoordinatorVCoreCount = pointer.From(props.CoordinatorVCores)
				state.ServerNames = flattenServerNames(props.ServerNames)
				state.HaEnabled = pointer.From(props.EnableHa)
				state.NodeCount = pointer.From(props.NodeCount)
				state.NodePublicIPAccessEnabled = pointer.From(props.NodeEnablePublicIPAccess)
				state.NodeServerEdition = pointer.From(props.NodeServerEdition)
				state.NodeStorageQuotaInMb = pointer.From(props.NodeStorageQuotaInMb)
				state.NodeVCores = pointer.From(props.NodeVCores)
				state.ShardsOnCoordinatorEnabled = pointer.From(props.EnableShardsOnCoordinator)
				state.CitusVersion = pointer.From(props.CitusVersion)
				state.PreferredPrimaryZone = pointer.From(props.PreferredPrimaryZone)
				state.SqlVersion = pointer.From(props.PostgresqlVersion)
				state.EarliestRestoreTime = pointer.From(props.EarliestRestoreTime)

				if v := props.MaintenanceWindow; v != nil {
					state.MaintenanceWindow = flattenMaintenanceWindow(v)
				}
			}

			if model.Tags != nil {
				state.Tags = *model.Tags
			}

			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
				return err
			}

			return metadata.Encode(&state)
		},
	}
}

func flattenMaintenanceWindow(input *clusters.MaintenanceWindow) []MaintenanceWindow {
	if input == nil || input.CustomWindow == nil || *input.CustomWindow == "Disabled" {
		return []MaintenanceWindow{}
	}

	return []MaintenanceWindow{
		{
			DayOfWeek:   pointer.From(input.DayOfWeek),
			StartHour:   pointer.From(input.StartHour),
			StartMinute: pointer.From(input.StartMinute),
		},
	}
}

func flattenServerNames(input *[]clusters.ServerNameItem) []ServerNameItem {
	if input == nil {
		return []ServerNameItem{}
	}

	output := make([]ServerNameItem, 0, len(*input))
	for _, v := range *input {
		output = append(output, ServerNameItem{
			FullyQualifiedDomainName: *v.FullyQualifiedDomainName,
			Name:                     *v.Name,
		})
	}

	return output
}

type ServerNameItem struct {
	Name                     string `tfschema:"name"`
	FullyQualifiedDomainName string `tfschema:"fqdn"`
}

type MaintenanceWindow struct {
	DayOfWeek   int64 `tfschema:"day_of_week"`
	StartHour   int64 `tfschema:"start_hour"`
	StartMinute int64 `tfschema:"start_minute"`
}
