// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package kubernetes_fleet_manager

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2024-04-01/fleets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerService.V20240401.Fleets
			schema := KubernetesFleetManagerResourceSchema{}

			id, err := fleets.ParseFleetID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if model := resp.Model; model != nil {
				schema.Name = id.FleetName
				schema.ResourceGroupName = id.ResourceGroupName
				r.mapFleetToKubernetesFleetManagerResourceSchema(*model, &schema)
			}

			return metadata.Encode(&schema)
		},
	}
}

type KubernetesFleetManagerResourceSchema struct {
	Location          string         `tfschema:"location"`
	Name              string         `tfschema:"name"`
	ResourceGroupName string         `tfschema:"resource_group_name"`
	Tags              map[string]any `tfschema:"tags"`
}

func (r Resource) mapKubernetesFleetManagerResourceSchemaToFleet(input KubernetesFleetManagerResourceSchema, output *fleets.Fleet) {
	output.Location = location.Normalize(input.Location)
	output.Tags = tags.Expand(input.Tags)

	if output.Properties == nil {
		output.Properties = &fleets.FleetProperties{}
	}
}

func (r Resource) mapFleetToKubernetesFleetManagerResourceSchema(input fleets.Fleet, output *KubernetesFleetManagerResourceSchema) {
	output.Location = location.Normalize(input.Location)
	output.Tags = tags.Flatten(input.Tags)

	if input.Properties == nil {
		input.Properties = &fleets.FleetProperties{}
	}
}
