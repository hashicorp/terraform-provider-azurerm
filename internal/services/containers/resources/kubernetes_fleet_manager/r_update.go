// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package kubernetes_fleet_manager

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2024-04-01/fleets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerService.V20240401.Fleets

			id, err := fleets.ParseFleetID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config KubernetesFleetManagerResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving existing %s: %+v", *id, err)
			}
			if existing.Model == nil {
				return fmt.Errorf("retrieving existing %s: properties was nil", *id)
			}
			payload := *existing.Model

			r.mapKubernetesFleetManagerResourceSchemaToFleet(config, &payload)

			if err := client.CreateOrUpdateThenPoll(ctx, *id, payload, fleets.DefaultCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
