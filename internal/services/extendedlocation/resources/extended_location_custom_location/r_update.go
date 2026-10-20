// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package extended_location_custom_location

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/extendedlocation/2021-08-15/customlocations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ExtendedLocation.CustomLocationsClient
			id, err := customlocations.ParseCustomLocationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state ExtendedLocationCustomLocationResourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return err
			}

			model := existing.Model

			if model.Properties == nil {
				return fmt.Errorf("retrieving properties for %s for update: %+v", *id, err)
			}
			d := metadata.ResourceData

			if d.HasChanges("authentication") {
				if len(state.Authentication) > 0 {
					auth := state.Authentication[0]
					model.Properties.Authentication = &customlocations.CustomLocationPropertiesAuthentication{
						Type:  pointer.To(auth.Type),
						Value: pointer.To(auth.Value),
					}
				}
			} else {
				model.Properties.Authentication = nil
			}

			if d.HasChange("display_name") {
				model.Properties.DisplayName = pointer.To(state.DisplayName)
			}

			if d.HasChange("cluster_extension_ids") {
				model.Properties.ClusterExtensionIds = pointer.To(state.ClusterExtensionIds)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *model); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}
			return nil
		},
	}
}
