// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_data_collection_endpoint

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2023-03-11/datacollectionendpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Monitor.DataCollectionEndpointsClient
			id, err := datacollectionendpoints.ParseDataCollectionEndpointID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					metadata.Logger.Infof("%s was not found - removing from state!", *id)
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
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

type DataCollectionEndpoint struct {
	ConfigurationAccessEndpoint string         `tfschema:"configuration_access_endpoint"`
	Description                 string         `tfschema:"description"`
	ImmutableId                 string         `tfschema:"immutable_id"`
	Kind                        string         `tfschema:"kind"`
	Name                        string         `tfschema:"name"`
	Location                    string         `tfschema:"location"`
	LogsIngestionEndpoint       string         `tfschema:"logs_ingestion_endpoint"`
	MetricsIngestionEndpoint    string         `tfschema:"metrics_ingestion_endpoint"`
	PublicNetworkAccessEnabled  bool           `tfschema:"public_network_access_enabled"`
	ResourceGroupName           string         `tfschema:"resource_group_name"`
	Tags                        map[string]any `tfschema:"tags"`
}
