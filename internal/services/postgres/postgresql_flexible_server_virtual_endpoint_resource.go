// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2024-08-01/servers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2024-08-01/virtualendpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type PostgresqlFlexibleServerVirtualEndpointResource struct{}

type PostgresqlFlexibleServerVirtualEndpointModel struct {
	Name            string `tfschema:"name"`
	SourceServerId  string `tfschema:"source_server_id"`
	ReplicaServerId string `tfschema:"replica_server_id"`
	Type            string `tfschema:"type"`
}

var _ sdk.ResourceWithUpdate = PostgresqlFlexibleServerVirtualEndpointResource{}

func (r PostgresqlFlexibleServerVirtualEndpointResource) ModelObject() interface{} {
	return &PostgresqlFlexibleServerVirtualEndpointModel{}
}

func (r PostgresqlFlexibleServerVirtualEndpointResource) ResourceType() string {
	return "azurerm_postgresql_flexible_server_virtual_endpoint"
}

func (r PostgresqlFlexibleServerVirtualEndpointResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return virtualendpoints.ValidateVirtualEndpointID
}

func (r PostgresqlFlexibleServerVirtualEndpointResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r PostgresqlFlexibleServerVirtualEndpointResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Description:  "The name of the Virtual Endpoint",
			ForceNew:     true,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"source_server_id": {
			Type:         pluginsdk.TypeString,
			Description:  "The Resource ID of the *Source* Postgres Flexible Server this should be associated with",
			ForceNew:     true,
			Required:     true,
			ValidateFunc: servers.ValidateFlexibleServerID,
		},
		"replica_server_id": {
			Type:         pluginsdk.TypeString,
			Description:  "The Resource ID of the *Replica* Postgres Flexible Server this should be associated with",
			Required:     true,
			ValidateFunc: servers.ValidateFlexibleServerID,
		},
		"type": {
			Type:         pluginsdk.TypeString,
			Description:  "The type of Virtual Endpoint",
			ForceNew:     true,
			Required:     true,
			ValidateFunc: validation.StringInSlice(virtualendpoints.PossibleValuesForVirtualEndpointType(), true),
		},
	}
}

func (r PostgresqlFlexibleServerVirtualEndpointResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var virtualEndpoint PostgresqlFlexibleServerVirtualEndpointModel

			if err := metadata.Decode(&virtualEndpoint); err != nil {
				return err
			}

			client := metadata.Client.Postgres.VirtualEndpointClient

			sourceServerId, err := servers.ParseFlexibleServerID(virtualEndpoint.SourceServerId)
			if err != nil {
				return err
			}

			replicaServerId, err := servers.ParseFlexibleServerID(virtualEndpoint.ReplicaServerId)
			if err != nil {
				return err
			}

			id := virtualendpoints.NewVirtualEndpointID(sourceServerId.SubscriptionId, sourceServerId.ResourceGroupName, sourceServerId.FlexibleServerName, virtualEndpoint.Name)

			locks.ByName(id.FlexibleServerName, postgresqlFlexibleServerResourceName)
			defer locks.UnlockByName(id.FlexibleServerName, postgresqlFlexibleServerResourceName)

			if replicaServerId.FlexibleServerName != id.FlexibleServerName {
				locks.ByName(replicaServerId.FlexibleServerName, postgresqlFlexibleServerResourceName)
				defer locks.UnlockByName(replicaServerId.FlexibleServerName, postgresqlFlexibleServerResourceName)
			}

			// This API can be a bit flaky if the same named resource is created/destroyed quickly
			// usually waiting a minute or two before redeploying is enough to resolve the conflict
			if err = client.CreateThenPoll(ctx, id, virtualendpoints.VirtualEndpointResource{
				Name: &virtualEndpoint.Name,
				Properties: &virtualendpoints.VirtualEndpointResourceProperties{
					EndpointType: pointer.To(virtualendpoints.VirtualEndpointType(virtualEndpoint.Type)),
					Members:      &[]string{replicaServerId.FlexibleServerName},
				},
			}); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r PostgresqlFlexibleServerVirtualEndpointResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Postgres.VirtualEndpointClient
			flexibleServerClient := metadata.Client.Postgres.FlexibleServersClient

			state := PostgresqlFlexibleServerVirtualEndpointModel{}

			id, err := virtualendpoints.ParseVirtualEndpointID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					// Check if a failover occurred by looking for the endpoint with the current name
					// but on what was previously the replica server
					olderState := PostgresqlFlexibleServerVirtualEndpointModel{}
					if err := metadata.Decode(&olderState); err == nil && olderState.ReplicaServerId != "" {
						// Try to find endpoint on what used to be the replica
						replicaServerId, parseErr := servers.ParseFlexibleServerID(olderState.ReplicaServerId)
						if parseErr == nil {
							// Create a new ID for the endpoint on the replica server
							newEndpointId := virtualendpoints.NewVirtualEndpointID(
								replicaServerId.SubscriptionId,
								replicaServerId.ResourceGroupName,
								replicaServerId.FlexibleServerName,
								id.VirtualEndpointName,
							)

							newResp, newErr := client.Get(ctx, newEndpointId)
							if newErr == nil && newResp.Model != nil && newResp.Model.Properties != nil {
								// Endpoint found on the replica - a failover likely occurred
								log.Printf("[INFO] Virtual endpoint %s found after failover on server %s", id.VirtualEndpointName, replicaServerId.FlexibleServerName)
								// Update the ID to the new endpoint location
								metadata.ResourceData.SetId(newEndpointId.ID())
								id = &newEndpointId
								resp = newResp
								err = nil
							}
						}
					}

					// If we still haven't found the endpoint, mark it as gone
					if err != nil {
						log.Printf("[INFO] %s does not exist - removing from state", metadata.ResourceData.Id())
						return metadata.MarkAsGone(id)
					}
				} else {
					return fmt.Errorf("retrieving %s: %+v", id, err)
				}
			}

			state.Name = id.VirtualEndpointName

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.Type = string(pointer.From(props.EndpointType))

					if props.Members == nil || len(*props.Members) == 0 {
						// if members list is nil or empty, this is an endpoint that was previously deleted
						log.Printf("[INFO] Postgresql Flexible Server Endpoint %q was previously deleted - removing from state", id.ID())
						return metadata.MarkAsGone(id)
					}

					// After a failover, the "source" server is now what used to be the replica server
					state.SourceServerId = servers.NewFlexibleServerID(id.SubscriptionId, id.ResourceGroupName, id.FlexibleServerName).ID()

					// Model.Properties.Members can contain 1 member which means source and replica are identical, or it can contain
					// 2 members when source and replica are different => [source_server_id, replication_server_name]
					replicaServerId := state.SourceServerId

					if len(*props.Members) > 0 {
						// Directly use the first member from the API response
						memberName := (*props.Members)[0]
						replicaServer, err := lookupFlexibleServerByName(ctx, flexibleServerClient, id, memberName, state.SourceServerId)
						if err != nil {
							// Failing to find the replica server doesn't need to be a hard error since we still have a functioning endpoint
							log.Printf("[WARN] Could not locate replica server with name %s: %+v", memberName, err)
						} else if replicaServer != nil {
							replicaId, err := servers.ParseFlexibleServerID(*replicaServer.Id)
							if err != nil {
								return err
							}
							replicaServerId = replicaId.ID()
						}
					}

					state.ReplicaServerId = replicaServerId
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r PostgresqlFlexibleServerVirtualEndpointResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Postgres.VirtualEndpointClient

			id, err := virtualendpoints.ParseVirtualEndpointID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByName(id.FlexibleServerName, postgresqlFlexibleServerResourceName)
			defer locks.UnlockByName(id.FlexibleServerName, postgresqlFlexibleServerResourceName)

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r PostgresqlFlexibleServerVirtualEndpointResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var virtualEndpoint PostgresqlFlexibleServerVirtualEndpointModel
			client := metadata.Client.Postgres.VirtualEndpointClient

			if err := metadata.Decode(&virtualEndpoint); err != nil {
				return err
			}

			id, err := virtualendpoints.ParseVirtualEndpointID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			// Get current state from API to handle failover scenarios
			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					// The endpoint might have moved to the replica server after failover
					// Check if the endpoint exists on replica instead
					replicaServerId, parseErr := servers.ParseFlexibleServerID(virtualEndpoint.ReplicaServerId)
					if parseErr == nil {
						// Construct endpoint ID using the replica server
						potentialEndpointId := virtualendpoints.NewVirtualEndpointID(
							replicaServerId.SubscriptionId,
							replicaServerId.ResourceGroupName,
							replicaServerId.FlexibleServerName,
							id.VirtualEndpointName,
						)

						// Check if endpoint exists on the replica server
						replicaResp, replicaErr := client.Get(ctx, potentialEndpointId)
						if replicaErr == nil && replicaResp.Model != nil {
							// Found on replica - update the ID to point to the new location
							log.Printf("[INFO] Found virtual endpoint %s on server %s after failover", id.VirtualEndpointName, replicaServerId.FlexibleServerName)
							metadata.ResourceData.SetId(potentialEndpointId.ID())
							id = &potentialEndpointId
							resp = replicaResp
							err = nil
						} else {
							return fmt.Errorf("getting %s: %+v", id, err)
						}
					} else {
						return fmt.Errorf("getting %s: %+v", id, err)
					}
				} else {
					return fmt.Errorf("getting %s: %+v", id, err)
				}
			}

			// Parse the replica server ID from the config
			replicaServerId, err := servers.ParseFlexibleServerID(virtualEndpoint.ReplicaServerId)
			if err != nil {
				return err
			}

			// Before updating, check if we're in a post-failover state
			// where the source and replica have been swapped
			sourceServerId, err := servers.ParseFlexibleServerID(virtualEndpoint.SourceServerId)
			if err != nil {
				return err
			}

			// If the resource ID's server doesn't match the source server, it likely means a failover occurred
			// Detect the swap and adjust our update request to match the new reality
			memberName := replicaServerId.FlexibleServerName
			if id.FlexibleServerName != sourceServerId.FlexibleServerName {
				log.Printf("[INFO] Detected server role swap in virtual endpoint - adapting update")

				// After failover, the former replica is now the primary (where the endpoint exists)
				// and we need to target the former primary as the replica
				memberName = sourceServerId.FlexibleServerName
			}

			locks.ByName(id.FlexibleServerName, postgresqlFlexibleServerResourceName)
			defer locks.UnlockByName(id.FlexibleServerName, postgresqlFlexibleServerResourceName)

			if memberName != id.FlexibleServerName {
				locks.ByName(memberName, postgresqlFlexibleServerResourceName)
				defer locks.UnlockByName(memberName, postgresqlFlexibleServerResourceName)
			}

			// Update with the correct member name based on current state
			if err := client.UpdateThenPoll(ctx, *id, virtualendpoints.VirtualEndpointResourceForPatch{
				Properties: &virtualendpoints.VirtualEndpointResourceProperties{
					EndpointType: pointer.To(virtualendpoints.VirtualEndpointType(virtualEndpoint.Type)),
					Members:      pointer.To([]string{memberName}),
				},
			}); err != nil {
				return fmt.Errorf("updating %q: %+v", id, err)
			}

			return nil
		},
	}
}

// The flexible endpoint API does not store the location/rg information on replicas it only stores the name.
// This lookup is safe because replicas for a given source server are *not* allowed to have identical names
func lookupFlexibleServerByName(ctx context.Context, flexibleServerClient *servers.ServersClient, virtualEndpointId *virtualendpoints.VirtualEndpointId, replicaServerName string, sourceServerId string) (*servers.Server, error) {
	postgresServers, err := flexibleServerClient.ListCompleteMatchingPredicate(ctx, commonids.NewSubscriptionID(virtualEndpointId.SubscriptionId), servers.ServerOperationPredicate{
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

	return nil, fmt.Errorf("could not locate postgres replica server with name %s", replicaServerName)
}
