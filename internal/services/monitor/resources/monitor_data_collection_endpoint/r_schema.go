// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_data_collection_endpoint

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2023-03-11/datacollectionendpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/preflight"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name which should be used for this Data Collection Endpoint.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"public_network_access_enabled": {
			Description: "Whether network access from public internet to the Data Collection Endpoint are allowed.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"description": {
			Description: "A description for the Data Collection Endpoint.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},

		"kind": {
			Description: "The kind of the Data Collection Endpoint.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			ValidateFunc: validation.StringInSlice(
				datacollectionendpoints.PossibleValuesForKnownDataCollectionEndpointResourceKind(), false,
			),
		},

		"tags": commonschema.Tags(),
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
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
	}
}

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			if metadata.ResourceDiff == nil {
				return nil
			}

			if metadata.Client.Features.EnhancedValidation.PreflightEnabled {
				// Only perform preflight validation if there are changes. This avoids validation failures and
				// additional API calls for resources that are unchanged between plan invocations.
				if len(metadata.ResourceDiff.GetChangedKeysPrefix("")) > 0 || metadata.ResourceDiff.Id() == "" {
					for _, key := range []string{"name", "resource_group_name", "location"} {
						if !metadata.ResourceDiff.NewValueKnown(key) {
							metadata.Logger.Infof("skipping preflight validation for `azurerm_monitor_data_collection_endpoint`: `%s` is not known until apply", key)
							return nil
						}
					}

					var model DataCollectionEndpoint
					if err := metadata.DecodeDiff(&model); err != nil {
						return err
					}

					req := expandCreateForDataCollectionEndpoint(model)

					id := datacollectionendpoints.NewDataCollectionEndpointID(metadata.Client.Account.SubscriptionId, model.ResourceGroupName, model.Name)

					preflightValidate, err := preflight.NewValidationRequest(pointer.To(model.Location), pointer.To(id), "2023-03-11", req)
					if err != nil {
						return fmt.Errorf("constructing preflight validation request: %w", err)
					}

					if err = preflightValidate.ValidateResource(ctx, metadata); err != nil {
						return err
					}
				}
			}

			return nil
		},
	}
}
