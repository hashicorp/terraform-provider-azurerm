// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dev_center_environment_type

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/environmenttypes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.V20250201.EnvironmentTypes

			id, err := environmenttypes.ParseDevCenterEnvironmentTypeID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.EnvironmentTypesGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := DevCenterEnvironmentTypeResourceModel{
				Name:        id.EnvironmentTypeName,
				DevCenterId: environmenttypes.NewDevCenterID(id.SubscriptionId, id.ResourceGroupName, id.DevCenterName).ID(),
			}

			if model := resp.Model; model != nil {
				state.Tags = pointer.From(model.Tags)
			}

			return metadata.Encode(&state)
		},
	}
}
