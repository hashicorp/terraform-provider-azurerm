// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package digitaltwins

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/digitaltwins/2023-01-31/timeseriesdatabaseconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2021-11-01/eventhubs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2023-05-02/clusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/digitaltwins/validate"
	eventhubValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/eventhub/validate"
	kustoValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type TimeSeriesDatabaseConnectionModel struct {
	Name                         string `tfschema:"name"`
	DigitalTwinsId               string `tfschema:"digital_twins_id"`
	EventhubConsumerGroupName    string `tfschema:"eventhub_consumer_group_name"`
	EventhubName                 string `tfschema:"eventhub_name"`
	EventhubNamespaceEndpointUri string `tfschema:"eventhub_namespace_endpoint_uri"`
	EventhubNamespaceId          string `tfschema:"eventhub_namespace_id"`
	KustoClusterId               string `tfschema:"kusto_cluster_id"`
	KustoClusterUri              string `tfschema:"kusto_cluster_uri"`
	KustoDatabaseName            string `tfschema:"kusto_database_name"`
	KustoTableName               string `tfschema:"kusto_table_name"`
}

type TimeSeriesDatabaseConnectionResource struct{}

func (m TimeSeriesDatabaseConnectionResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.DigitalTwinsTimeSeriesDatabaseConnectionName,
		},

		"digital_twins_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: timeseriesdatabaseconnections.ValidateDigitalTwinsInstanceID,
		},

		"eventhub_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: eventhubValidate.ValidateEventHubName(),
		},

		"eventhub_namespace_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: eventhubs.ValidateNamespaceID,
		},

		"eventhub_namespace_endpoint_uri": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"kusto_cluster_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: clusters.ValidateClusterID,
		},

		"kusto_cluster_uri": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"kusto_database_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: kustoValidate.DatabaseName,
		},

		"eventhub_consumer_group_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "$Default",
			ForceNew:     true,
			ValidateFunc: eventhubValidate.ValidateEventHubConsumerName(),
		},

		"kusto_table_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: kustoValidate.EntityName,
		},
	}
}

func (m TimeSeriesDatabaseConnectionResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (m TimeSeriesDatabaseConnectionResource) ModelObject() interface{} {
	return &TimeSeriesDatabaseConnectionModel{}
}

func (m TimeSeriesDatabaseConnectionResource) ResourceType() string {
	return "azurerm_digital_twins_time_series_database_connection"
}

func (m TimeSeriesDatabaseConnectionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return timeseriesdatabaseconnections.ValidateTimeSeriesDatabaseConnectionID
}

func (m TimeSeriesDatabaseConnectionResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			client := meta.Client.DigitalTwins.TimeSeriesDatabaseConnectionsClient

			var model TimeSeriesDatabaseConnectionModel
			if err := meta.Decode(&model); err != nil {
				return err
			}

			digitalTwinsId, err := timeseriesdatabaseconnections.ParseDigitalTwinsInstanceID(model.DigitalTwinsId)
			if err != nil {
				return err
			}

			id := timeseriesdatabaseconnections.NewTimeSeriesDatabaseConnectionID(digitalTwinsId.SubscriptionId, digitalTwinsId.ResourceGroupName, digitalTwinsId.DigitalTwinsInstanceName, model.Name)

			existing, err := client.Get(ctx, id)
			if !response.WasNotFound(existing.HttpResponse) {
				if err != nil {
					return fmt.Errorf("retrieving %s: %+v", id, err)
				}
				return meta.ResourceRequiresImport(m.ResourceType(), id)
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
				properties.AdxTableName = utils.String(model.KustoTableName)
			}

			if model.EventhubConsumerGroupName != "" {
				properties.EventHubConsumerGroup = utils.String(model.EventhubConsumerGroupName)
			}

			req := timeseriesdatabaseconnections.TimeSeriesDatabaseConnection{
				Properties: properties,
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, req); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			meta.SetID(id)
			return nil
		},
	}
}

func (m TimeSeriesDatabaseConnectionResource) Read() sdk.ResourceFunc {
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

				kustoClusterId, err := clusters.ParseClusterIDInsensitively(properties.AdxResourceId)
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

				kustoTableName := ""
				if properties.AdxTableName != nil {
					kustoTableName = *properties.AdxTableName
				}
				output.KustoTableName = kustoTableName
			}

			return meta.Encode(&output)
		},
	}
}

func (m TimeSeriesDatabaseConnectionResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			id, err := timeseriesdatabaseconnections.ParseTimeSeriesDatabaseConnectionID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			client := meta.Client.DigitalTwins.TimeSeriesDatabaseConnectionsClient
			if err = client.DeleteThenPoll(ctx, *id, timeseriesdatabaseconnections.DefaultDeleteOperationOptions()); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}
			return nil
		},
	}
}
