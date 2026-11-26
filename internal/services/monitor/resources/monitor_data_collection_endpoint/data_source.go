// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_data_collection_endpoint

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2023-03-11/datacollectionendpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type DataSource struct{}

var _ sdk.DataSource = DataSource{}

func (d DataSource) ModelObject() any {
	return &DataCollectionEndpoint{}
}

func (d DataSource) ResourceType() string {
	return "azurerm_monitor_data_collection_endpoint"
}

func (d DataSource) Description() string {
	return "Use this data source to access information about an existing Data Collection Endpoint."
}

func (d DataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Data Collection Endpoint.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (d DataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"configuration_access_endpoint": {
			Description: "The endpoint used for accessing configuration, e.g., `https://mydce-abcd.eastus-1.control.monitor.azure.com`.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"immutable_id": {
			Description: "The immutable ID of the Data Collection Endpoint.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"location": commonschema.LocationComputed(),

		"logs_ingestion_endpoint": {
			Description: "The endpoint used for ingesting logs, e.g., `https://mydce-abcd.eastus-1.ingest.monitor.azure.com`.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"metrics_ingestion_endpoint": {
			Description: "The endpoint used for ingesting metrics, e.g., `https://mydce-abcd.eastus-1.metrics.ingest.monitor.azure.com`.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"public_network_access_enabled": {
			Description: "Whether network access from public internet to the Data Collection Endpoint are allowed.",
			Type:        pluginsdk.TypeBool,
			Computed:    true,
		},

		"description": {
			Description: "A description for the Data Collection Endpoint.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"kind": {
			Description: "The kind of the Data Collection Endpoint.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"tags": commonschema.TagsDataSource(),
	}
}

func (d DataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Monitor.DataCollectionEndpointsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state DataCollectionEndpoint
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := datacollectionendpoints.NewDataCollectionEndpointID(subscriptionId, state.ResourceGroupName, state.Name)
			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			var publicNetWorkAccessEnabled bool
			var description, kind, loc, configurationAccessEndpoint, logsIngestionEndpoint, metricsIngestionEndpoint, immutableId string
			var tag map[string]any
			if model := resp.Model; model != nil {
				kind = flattenDataCollectionEndpointKind(model.Kind)
				loc = location.Normalize(model.Location)
				tag = tags.Flatten(model.Tags)
				if prop := model.Properties; prop != nil {
					description = flattenDataCollectionEndpointDescription(prop.Description)
					if networkAcls := prop.NetworkAcls; networkAcls != nil {
						publicNetWorkAccessEnabled = flattenDataCollectionEndpointPublicNetworkAccess(networkAcls.PublicNetworkAccess)
					}

					if prop.ConfigurationAccess != nil && prop.ConfigurationAccess.Endpoint != nil {
						configurationAccessEndpoint = *prop.ConfigurationAccess.Endpoint
					}

					if prop.LogsIngestion != nil && prop.LogsIngestion.Endpoint != nil {
						logsIngestionEndpoint = *prop.LogsIngestion.Endpoint
					}

					if prop.MetricsIngestion != nil && prop.MetricsIngestion.Endpoint != nil {
						metricsIngestionEndpoint = *prop.MetricsIngestion.Endpoint
					}

					if prop.ImmutableId != nil {
						immutableId = *prop.ImmutableId
					}
				}
			}

			metadata.SetID(id)

			return metadata.Encode(&DataCollectionEndpoint{
				ConfigurationAccessEndpoint: configurationAccessEndpoint,
				Description:                 description,
				Kind:                        kind,
				ImmutableId:                 immutableId,
				Location:                    loc,
				LogsIngestionEndpoint:       logsIngestionEndpoint,
				MetricsIngestionEndpoint:    metricsIngestionEndpoint,
				Name:                        id.DataCollectionEndpointName,
				PublicNetworkAccessEnabled:  publicNetWorkAccessEnabled,
				ResourceGroupName:           id.ResourceGroupName,
				Tags:                        tag,
			})
		},
		Timeout: 5 * time.Minute,
	}
}
