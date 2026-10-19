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

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Dashboard.ManagedPrivateEndpointsClient
			id, err := managedprivateendpointmodels.ParseManagedPrivateEndpointID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.ManagedPrivateEndpointsDeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
