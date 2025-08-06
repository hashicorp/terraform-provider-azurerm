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
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

var (
	_ sdk.Resource           = DataCollectionEndpointResource{}
	_ sdk.ResourceWithUpdate = DataCollectionEndpointResource{}
)

type DataCollectionEndpointResource struct{}

type DataCollectionEndpoint struct {
	ConfigurationAccessEndpoint string                 `tfschema:"configuration_access_endpoint"`
	Description                 string                 `tfschema:"description"`
	ImmutableId                 string                 `tfschema:"immutable_id"`
	Kind                        string                 `tfschema:"kind"`
	Name                        string                 `tfschema:"name"`
	Location                    string                 `tfschema:"location"`
	LogsIngestionEndpoint       string                 `tfschema:"logs_ingestion_endpoint"`
	MetricsIngestionEndpoint    string                 `tfschema:"metrics_ingestion_endpoint"`
	PublicNetworkAccessEnabled  bool                   `tfschema:"public_network_access_enabled"`
	ResourceGroupName           string                 `tfschema:"resource_group_name"`
	Tags                        map[string]interface{} `tfschema:"tags"`
}

func (r DataCollectionEndpointResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"kind": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice(
				datacollectionendpoints.PossibleValuesForKnownDataCollectionEndpointResourceKind(), false),
		},

		"tags": commonschema.Tags(),
	}
}

func (r DataCollectionEndpointResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"configuration_access_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"immutable_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"logs_ingestion_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"metrics_ingestion_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r DataCollectionEndpointResource) ResourceType() string {
	return "azurerm_monitor_data_collection_endpoint"
}

func (r DataCollectionEndpointResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return datacollectionendpoints.ValidateDataCollectionEndpointID
}

func (r DataCollectionEndpointResource) ModelObject() interface{} {
	return &DataCollectionEndpoint{}
}

func (r DataCollectionEndpointResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			metadata.Logger.Info("Decoding state..")
			var state DataCollectionEndpoint
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			client := metadata.Client.Monitor.DataCollectionEndpointsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := datacollectionendpoints.NewDataCollectionEndpointID(subscriptionId, state.ResourceGroupName, state.Name)
			metadata.Logger.Infof("creating %s", id)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			input := datacollectionendpoints.DataCollectionEndpointResource{
				Kind:     expandDataCollectionEndpointKind(state.Kind),
				Location: azure.NormalizeLocation(state.Location),
				Name:     utils.String(state.Name),
				Properties: &datacollectionendpoints.DataCollectionEndpoint{
					Description: utils.String(state.Description),
					NetworkAcls: &datacollectionendpoints.NetworkRuleSet{
						PublicNetworkAccess: expandDataCollectionEndpointPublicNetworkAccess(state.PublicNetworkAccessEnabled),
					},
				},
				Tags: tags.Expand(state.Tags),
			}

			if _, err := client.Create(ctx, id, input); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r DataCollectionEndpointResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Monitor.DataCollectionEndpointsClient
			id, err := datacollectionendpoints.ParseDataCollectionEndpointID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("retrieving %s", *id)
			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					metadata.Logger.Infof("%s was not found - removing from state!", *id)
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
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

func (r DataCollectionEndpointResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := datacollectionendpoints.ParseDataCollectionEndpointID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("updating %s..", *id)
			client := metadata.Client.Monitor.DataCollectionEndpointsClient
			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if resp.Model == nil {
				return fmt.Errorf("unexpected null model of %s", *id)
			}
			existing := resp.Model
			if existing.Properties == nil {
				return fmt.Errorf("unexpected null properties of %s", *id)
			}

			var state DataCollectionEndpoint
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			if metadata.ResourceData.HasChange("description") {
				existing.Properties.Description = utils.String(state.Description)
			}

			if metadata.ResourceData.HasChange("kind") {
				existing.Kind = expandDataCollectionEndpointKind(state.Kind)
			}

			if metadata.ResourceData.HasChange("public_network_access_enabled") {
				existing.Properties.NetworkAcls = &datacollectionendpoints.NetworkRuleSet{
					PublicNetworkAccess: expandDataCollectionEndpointPublicNetworkAccess(state.PublicNetworkAccessEnabled),
				}
			}

			if metadata.ResourceData.HasChange("tags") {
				existing.Tags = tags.Expand(state.Tags)
			}

			if _, err := client.Create(ctx, *id, *existing); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r DataCollectionEndpointResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Monitor.DataCollectionEndpointsClient
			id, err := datacollectionendpoints.ParseDataCollectionEndpointID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("deleting %s..", *id)
			resp, err := client.Delete(ctx, *id)
			if err != nil && !response.WasNotFound(resp.HttpResponse) {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func expandDataCollectionEndpointKind(input string) *datacollectionendpoints.KnownDataCollectionEndpointResourceKind {
	if input == "" {
		return nil
	}

	result := datacollectionendpoints.KnownDataCollectionEndpointResourceKind(input)
	return &result
}

func expandDataCollectionEndpointPublicNetworkAccess(input bool) *datacollectionendpoints.KnownPublicNetworkAccessOptions {
	var result datacollectionendpoints.KnownPublicNetworkAccessOptions
	if input {
		result = datacollectionendpoints.KnownPublicNetworkAccessOptionsEnabled
	} else {
		result = datacollectionendpoints.KnownPublicNetworkAccessOptionsDisabled
	}
	return &result
}

func flattenDataCollectionEndpointKind(input *datacollectionendpoints.KnownDataCollectionEndpointResourceKind) string {
	if input == nil {
		return ""
	}

	return string(*input)
}

func flattenDataCollectionEndpointDescription(input *string) string {
	if input == nil {
		return ""
	}

	return *input
}

func flattenDataCollectionEndpointPublicNetworkAccess(input *datacollectionendpoints.KnownPublicNetworkAccessOptions) bool {
	if input == nil {
		return false
	}
	var result bool
	if *input == datacollectionendpoints.KnownPublicNetworkAccessOptionsEnabled {
		result = true
	} else if *input == datacollectionendpoints.KnownPublicNetworkAccessOptionsDisabled {
		result = false
	}
	return result
}
