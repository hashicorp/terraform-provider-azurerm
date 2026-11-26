// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package digital_twins_time_series_database_connection

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/digitaltwins/2023-01-31/timeseriesdatabaseconnections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (m Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DigitalTwins.TimeSeriesDatabaseConnectionsClient

			var model TimeSeriesDatabaseConnectionModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			digitalTwinsId, err := timeseriesdatabaseconnections.ParseDigitalTwinsInstanceID(model.DigitalTwinsId)
			if err != nil {
				return err
			}

			id := timeseriesdatabaseconnections.NewTimeSeriesDatabaseConnectionID(digitalTwinsId.SubscriptionId, digitalTwinsId.ResourceGroupName, digitalTwinsId.DigitalTwinsInstanceName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if !response.WasNotFound(existing.HttpResponse) {
					if err != nil {
						return fmt.Errorf("retrieving %s: %+v", id, err)
					}
					return metadata.ResourceRequiresImport(m.ResourceType(), id)
				}
			}

			properties := timeseriesdatabaseconnections.AzureDataExplorerConnectionProperties{
				AdxDatabaseName:             model.KustoDatabaseName,
				AdxEndpointUri:              model.KustoClusterUri,
				AdxResourceId:               model.KustoClusterId,
				EventHubEndpointUri:         model.EventhubNamespaceEndpointUri,
				EventHubEntityPath:          model.EventhubName,
				EventHubNamespaceResourceId: model.EventhubNamespaceId,
			}

			if model.KustoTableName != "" {
				properties.AdxTableName = pointer.To(model.KustoTableName)
			}

			if model.EventhubConsumerGroupName != "" {
				properties.EventHubConsumerGroup = pointer.To(model.EventhubConsumerGroupName)
			}

			req := timeseriesdatabaseconnections.TimeSeriesDatabaseConnection{
				Properties: properties,
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, req, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
