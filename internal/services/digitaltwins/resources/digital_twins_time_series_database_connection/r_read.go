// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package digital_twins_time_series_database_connection

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/digitaltwins/2023-01-31/timeseriesdatabaseconnections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (m Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			id, err := timeseriesdatabaseconnections.ParseTimeSeriesDatabaseConnectionID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			client := meta.Client.DigitalTwins.TimeSeriesDatabaseConnectionsClient
			result, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(result.HttpResponse) {
					return meta.MarkAsGone(id)
				}
				return err
			}

			if result.Model == nil {
				return fmt.Errorf("retrieving %s got nil model", id)
			}

			var output TimeSeriesDatabaseConnectionModel
			output.Name = id.TimeSeriesDatabaseConnectionName
			output.DigitalTwinsId = timeseriesdatabaseconnections.NewDigitalTwinsInstanceID(id.SubscriptionId, id.ResourceGroupName, id.DigitalTwinsInstanceName).ID()

			if properties, ok := result.Model.Properties.(timeseriesdatabaseconnections.AzureDataExplorerConnectionProperties); ok {
				output.EventhubName = properties.EventHubEntityPath
				output.EventhubNamespaceEndpointUri = properties.EventHubEndpointUri
				output.EventhubNamespaceId = properties.EventHubNamespaceResourceId

				kustoClusterId, err := commonids.ParseKustoClusterIDInsensitively(properties.AdxResourceId)
				if err != nil {
					return fmt.Errorf("parsing `kusto_cluster_uri`: %+v", err)
				}
				output.KustoClusterId = kustoClusterId.ID()

				output.KustoClusterUri = properties.AdxEndpointUri
				output.KustoDatabaseName = properties.AdxDatabaseName

				eventhubConsumerGroupName := "$Default"
				if properties.EventHubConsumerGroup != nil {
					eventhubConsumerGroupName = *properties.EventHubConsumerGroup
				}
				output.EventhubConsumerGroupName = eventhubConsumerGroupName

				output.KustoTableName = pointer.From(properties.AdxTableName)
			}

			return meta.Encode(&output)
		},
	}
}
