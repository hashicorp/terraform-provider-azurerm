// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package monitor

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2023-03-11/datacollectionendpoints"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type DataCollectionEndpointDataSource struct{}

var _ sdk.DataSource = DataCollectionEndpointDataSource{}

func (d DataCollectionEndpointDataSource) ModelObject() interface{} {
	return &DataCollectionEndpoint{}
}

func (d DataCollectionEndpointDataSource) ResourceType() string {
	return "azurerm_monitor_data_collection_endpoint"
}

func (d DataCollectionEndpointDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (d DataCollectionEndpointDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"configuration_access_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"immutable_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"location": commonschema.LocationComputed(),

		"logs_ingestion_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"metrics_ingestion_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"kind": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tags": commonschema.TagsDataSource(),
	}
}

func (d DataCollectionEndpointDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Monitor.DataCollectionEndpointsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state DataCollectionEndpoint
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := datacollectionendpoints.NewDataCollectionEndpointID(subscriptionId, state.ResourceGroupName, state.Name)
			metadata.Logger.Infof("retrieving %s", id)
			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			var publicNetWorkAccessEnabled bool
			var description, kind, location, configurationAccessEndpoint, logsIngestionEndpoint, metricsIngestionEndpoint, immutableId string
			var tag map[string]interface{}
			if model := resp.Model; model != nil {
				kind = flattenDataCollectionEndpointKind(model.Kind)
				location = azure.NormalizeLocation(model.Location)
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
				Location:                    location,
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
