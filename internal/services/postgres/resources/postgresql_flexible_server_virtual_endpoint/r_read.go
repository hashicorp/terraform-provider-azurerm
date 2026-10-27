// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package postgresql_flexible_server_virtual_endpoint

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2025-08-01/servers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2025-08-01/virtualendpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Postgres.VirtualEndpointClient
			flexibleServerClient := metadata.Client.Postgres.FlexibleServersClient

			state := PostgresqlFlexibleServerVirtualEndpointModel{}

			id, err := commonids.ParseCompositeResourceID(metadata.ResourceData.Id(), &virtualendpoints.VirtualEndpointId{}, &virtualendpoints.VirtualEndpointId{})
			if err != nil {
				return err
			}

			// In case of a fail-over, we need to see if the endpoint lives under the source id or the replica id
			failOverHasOccurred := false
			virtualEndpointId := *id.First
			resp, err := client.Get(ctx, virtualEndpointId)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) || response.WasBadRequest(resp.HttpResponse) { // Can return a 400 if attempting to query the replica for the virtual endpoint resource
					virtualEndpointId = *id.Second
					// if the endpoint doesn't exist under the source server, look for it under the replica server
					resp, err = client.Get(ctx, virtualEndpointId)
					if err != nil {
						if response.WasNotFound(resp.HttpResponse) || response.WasBadRequest(resp.HttpResponse) { // Can return a 400 if attempting to query the replica for the virtual endpoint resource
							// the endpoint was not found under the source or the replica server so it can safely be removed from state
							log.Printf("[INFO] %s does not exist - removing from state", metadata.ResourceData.Id())
							return metadata.MarkAsGone(id)
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

			state.Name = virtualEndpointId.VirtualEndpointName

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.Type = pointer.FromEnum(props.EndpointType)

					if props.Members == nil || len(*props.Members) == 0 {
						// if members list is nil or empty, this is an endpoint that was previously deleted
						log.Printf("[INFO] Postgresql Flexible Server Endpoint %q was previously deleted - removing from state", id.First.ID())
						return metadata.MarkAsGone(id)
					}

					state.SourceServerId = servers.NewFlexibleServerID(id.First.SubscriptionId, id.First.ResourceGroupName, (*props.Members)[0]).ID()
					// Model.Properties.Members can contain 1 member which means source and replica are identical, or it can contain
					// 2 members when source and replica are different => [source_server_id, replication_server_name]
					replicaServerId := state.SourceServerId

					if len(*props.Members) == 2 {
						replicaServer, err := lookupFlexibleServerByName(ctx, flexibleServerClient, virtualEndpointId, (*props.Members)[1], state.SourceServerId)
						if err != nil {
							return err
						}

						if replicaServer != nil {
							replicaId, err := servers.ParseFlexibleServerID(*replicaServer.Id)
							if err != nil {
								return err
							}

							replicaServerId = replicaId.ID()
						} else {
							// if the replica server is not found, it may exists in different subscription, read from the metadata
							replicaServerId = metadata.ResourceData.Get("replica_server_id").(string)
						}
					}

					state.ReplicaServerId = replicaServerId
				}
			}

			// if a fail-over has occurred, the source/replica ids have swapped so we'll have to swap them back in Terraform to prevent a diff
			if failOverHasOccurred {
				state.SourceServerId, state.ReplicaServerId = state.ReplicaServerId, state.SourceServerId
			}

			return metadata.Encode(&state)
		},
	}
}

// The flexible endpoint API does not store the location/rg information on replicas it only stores the name.
// This lookup is safe because replicas for a given source server are *not* allowed to have identical names
func lookupFlexibleServerByName(ctx context.Context, flexibleServerClient *servers.ServersClient, virtualEndpointId virtualendpoints.VirtualEndpointId, replicaServerName string, sourceServerId string) (*servers.Server, error) {
	postgresServers, err := flexibleServerClient.ListBySubscriptionCompleteMatchingPredicate(ctx, commonids.NewSubscriptionID(virtualEndpointId.SubscriptionId), servers.ServerOperationPredicate{
		Name: &replicaServerName,
	})
	if err != nil {
		return nil, err
	}

	// loop to find the replica server associated with this flexible endpoint
	for _, server := range postgresServers.Items {
		if server.Properties != nil && pointer.From(server.Properties.SourceServerResourceId) == sourceServerId {
			return &server, nil
		}
	}

	return nil, nil
}
