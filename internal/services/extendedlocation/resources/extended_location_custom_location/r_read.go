// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package extended_location_custom_location

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/extendedlocation/2021-08-15/customlocations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ExtendedLocation.CustomLocationsClient
			id, err := customlocations.ParseCustomLocationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if model := resp.Model; model != nil {
				props := model.Properties

				state := ExtendedLocationCustomLocationResourceModel{
					Name:                id.CustomLocationName,
					ResourceGroupName:   id.ResourceGroupName,
					Location:            location.Normalize(model.Location),
					ClusterExtensionIds: pointer.From(props.ClusterExtensionIds),
					DisplayName:         pointer.From(props.DisplayName),
					HostResourceId:      pointer.From(props.HostResourceId),
					HostType:            pointer.FromEnum(props.HostType),
					Namespace:           pointer.From(props.Namespace),
				}

				// API always returns an empty `authentication` block even it's not specified. Tracing the bug: https://github.com/Azure/azure-rest-api-specs/issues/30101
				if props.Authentication != nil && props.Authentication.Type != nil && props.Authentication.Value != nil {
					state.Authentication = []AuthModel{
						{
							Type:  pointer.From(props.Authentication.Type),
							Value: pointer.From(props.Authentication.Value),
						},
					}
				}

				return metadata.Encode(&state)
			}
			return nil
		},
	}
}
