// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package postgresql_flexible_server_virtual_endpoint

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2025-08-01/servers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2025-08-01/virtualendpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var virtualEndpoint PostgresqlFlexibleServerVirtualEndpointModel
			client := metadata.Client.Postgres.VirtualEndpointClient

			if err := metadata.Decode(&virtualEndpoint); err != nil {
				return err
			}

			id, err := commonids.ParseCompositeResourceID(metadata.ResourceData.Id(), &virtualendpoints.VirtualEndpointId{}, &virtualendpoints.VirtualEndpointId{})
			if err != nil {
				return err
			}

			// attempt to retrieve the endpoint and see if a fail-over has occurred, if so error as we shouldn't update to a different replica server with the `source_server_id` and the `replica_server_id` being swapped
			virtualEndpointId := *id.First
			resp, err := client.Get(ctx, virtualEndpointId)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					virtualEndpointId = virtualendpoints.NewVirtualEndpointID(id.Second.SubscriptionId, id.Second.ResourceGroupName, id.Second.FlexibleServerName, id.Second.VirtualEndpointName)
					// if the endpoint doesn't exist under the source server, look for it under the replica server
					if _, err = client.Get(ctx, virtualEndpointId); err != nil {
						return fmt.Errorf("retrieving %s: %+v", virtualEndpointId, err)
					}
					return fmt.Errorf("a fail-over has occurred and the `source_server_id` in the config is no longer the SourceServerId for the virtual endpoint. If you wish to change the `replica_server_id`, remove this resource from state and reimport it back in with the `replica_server_id` and `source_server_id` swapped")
				}
				return fmt.Errorf("retrieving %s: %+v", virtualEndpointId, err)
			}

			replicaServerId, err := servers.ParseFlexibleServerID(virtualEndpoint.ReplicaServerId)
			if err != nil {
				return err
			}

			locks.ByName(id.First.FlexibleServerName, postgresqlFlexibleServerResourceName)
			defer locks.UnlockByName(id.First.FlexibleServerName, postgresqlFlexibleServerResourceName)

			if replicaServerId.FlexibleServerName != id.First.FlexibleServerName {
				locks.ByName(replicaServerId.FlexibleServerName, postgresqlFlexibleServerResourceName)
				defer locks.UnlockByName(replicaServerId.FlexibleServerName, postgresqlFlexibleServerResourceName)
			}

			endpointId := virtualendpoints.NewVirtualEndpointID(id.First.SubscriptionId, id.First.ResourceGroupName, id.First.FlexibleServerName, virtualEndpoint.Name)
			if err := client.UpdateThenPoll(ctx, endpointId, virtualendpoints.VirtualEndpointResourceForPatch{
				Properties: &virtualendpoints.VirtualEndpointResourceProperties{
					EndpointType: pointer.ToEnum[virtualendpoints.VirtualEndpointType](virtualEndpoint.Type),
					Members:      pointer.To([]string{replicaServerId.FlexibleServerName}),
				},
			}); err != nil {
				return fmt.Errorf("updating %q: %+v", id, err)
			}

			// the id has changed and needs to be updated
			replicaEndpointId := virtualendpoints.NewVirtualEndpointID(replicaServerId.SubscriptionId, replicaServerId.ResourceGroupName, replicaServerId.FlexibleServerName, virtualEndpoint.Name)
			metadata.SetID(commonids.NewCompositeResourceID(&virtualEndpointId, &replicaEndpointId))

			return nil
		},
	}
}
