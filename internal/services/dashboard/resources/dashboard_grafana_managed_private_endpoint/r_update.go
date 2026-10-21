// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dashboard_grafana_managed_private_endpoint

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/dashboard/2025-08-01/managedprivateendpointmodels"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Dashboard.ManagedPrivateEndpointsClient

			id, err := managedprivateendpointmodels.ParseManagedPrivateEndpointID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var mpe ManagedPrivateEndpointModel
			if err := metadata.Decode(&mpe); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.ManagedPrivateEndpointsGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			model := resp.Model
			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}

			if metadata.ResourceData.HasChange("tags") {
				model.Tags = &mpe.Tags
			}

			if metadata.ResourceData.HasChange("private_link_service_url") {
				model.Properties.PrivateLinkServiceURL = &mpe.PrivateLinkServiceURL
			}

			if err := client.ManagedPrivateEndpointsCreateThenPoll(ctx, *id, *model); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
