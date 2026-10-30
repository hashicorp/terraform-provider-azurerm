// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package postgresql_flexible_server_virtual_endpoint

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2025-08-01/virtualendpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Postgres.VirtualEndpointClient

			id, err := commonids.ParseCompositeResourceID(metadata.ResourceData.Id(), &virtualendpoints.VirtualEndpointId{}, &virtualendpoints.VirtualEndpointId{})
			if err != nil {
				return err
			}

			// In case of a fail-over, we need to see if the endpoint lives under the source id or the replica id before deleting
			failOverHasOccurred := false
			virtualEndpointId := *id.First
			resp, err := client.Get(ctx, virtualEndpointId)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					virtualEndpointId = *id.Second
					// if the endpoint doesn't exist under the source server, look for it under the replica server
					resp, err = client.Get(ctx, virtualEndpointId)
					if err != nil {
						if response.WasNotFound(resp.HttpResponse) || response.WasBadRequest(resp.HttpResponse) { // Can return a 400 if attempting to query the replica for the virtual endpoint resource
							// the endpoint was not found under the source or the replica server so we can exit here
							return nil
						}
						return fmt.Errorf("retrieving %s: %+v", id, err)
					}
					failOverHasOccurred = true
				}
				// if we errored and didn't find the endpoint under the replica id, then we error here
				if !failOverHasOccurred {
					return fmt.Errorf("retrieving %s: %+v", id, err)
				}
			}

			locks.ByName(virtualEndpointId.FlexibleServerName, postgresqlFlexibleServerResourceName)
			defer locks.UnlockByName(virtualEndpointId.FlexibleServerName, postgresqlFlexibleServerResourceName)

			if err := client.DeleteThenPoll(ctx, virtualEndpointId); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
