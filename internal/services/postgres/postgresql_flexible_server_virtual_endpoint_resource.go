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
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/postgres/migration"
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

var _ sdk.ResourceWithStateMigration = PostgresqlFlexibleServerVirtualEndpointResource{}

func (r PostgresqlFlexibleServerVirtualEndpointResource) ModelObject() interface{} {
	return &PostgresqlFlexibleServerVirtualEndpointModel{}
}

func (r PostgresqlFlexibleServerVirtualEndpointResource) ResourceType() string {
	return "azurerm_postgresql_flexible_server_virtual_endpoint"
}

func (r PostgresqlFlexibleServerVirtualEndpointResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return func(input interface{}, key string) (warnings []string, errors []error) {
		v, ok := input.(string)
		if !ok {
			errors = append(errors, fmt.Errorf("expected %q to be a string", key))
			return
		}

		if _, err := commonids.ParseCompositeResourceID(v, &virtualendpoints.VirtualEndpointId{}, &virtualendpoints.VirtualEndpointId{}); err != nil {
			errors = append(errors, err)
		}
		return
	}
}

func (r PostgresqlFlexibleServerVirtualEndpointResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r PostgresqlFlexibleServerVirtualEndpointResource) StateUpgraders() sdk.StateUpgradeData {
	return sdk.StateUpgradeData{
		SchemaVersion: 1,
		Upgraders: map[int]pluginsdk.StateUpgrade{
			0: migration.PostgresqlFlexibleServerVirtualEndpointV0toV1{},
		},
	}
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

			sourceEndpointId := virtualendpoints.NewVirtualEndpointID(sourceServerId.SubscriptionId, sourceServerId.ResourceGroupName, sourceServerId.FlexibleServerName, virtualEndpoint.Name)
			replicaEndpointId := virtualendpoints.NewVirtualEndpointID(replicaServerId.SubscriptionId, replicaServerId.ResourceGroupName, replicaServerId.FlexibleServerName, virtualEndpoint.Name)

			locks.ByName(sourceEndpointId.FlexibleServerName, postgresqlFlexibleServerResourceName)
			defer locks.UnlockByName(sourceEndpointId.FlexibleServerName, postgresqlFlexibleServerResourceName)

			if replicaServerId.FlexibleServerName != replicaEndpointId.FlexibleServerName {
				locks.ByName(replicaServerId.FlexibleServerName, postgresqlFlexibleServerResourceName)
				defer locks.UnlockByName(replicaServerId.FlexibleServerName, postgresqlFlexibleServerResourceName)
			}

			// This API can be a bit flaky if the same named resource is created/destroyed quickly
			// usually waiting a minute or two before redeploying is enough to resolve the conflict
			if err = client.CreateThenPoll(ctx, sourceEndpointId, virtualendpoints.VirtualEndpointResource{
				Name: &virtualEndpoint.Name,
				Properties: &virtualendpoints.VirtualEndpointResourceProperties{
					EndpointType: pointer.To(virtualendpoints.VirtualEndpointType(virtualEndpoint.Type)),
					Members:      &[]string{replicaServerId.FlexibleServerName},
				},
			}); err != nil {
				return fmt.Errorf("creating %s: %+v", sourceEndpointId, err)
			}

			id := commonids.NewCompositeResourceID(&sourceEndpointId, &replicaEndpointId)

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

			id, err := commonids.ParseCompositeResourceID(metadata.ResourceData.Id(), &virtualendpoints.VirtualEndpointId{}, &virtualendpoints.VirtualEndpointId{})
			if err != nil {
				return err
			}

			// In case of a fail-over, we need to see if the endpoint lives under the source id or the replica id
			failOverHasOccurred := false
			virtualEndpointId := *id.First
			resp, err := client.Get(ctx, virtualEndpointId)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					virtualEndpointId = *id.Second
					// if the endpoint doesn't exist under the source server, look for it under the replica server
					resp, err = client.Get(ctx, virtualEndpointId)
					if err != nil {
						if response.WasNotFound(resp.HttpResponse) {
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
					state.Type = string(pointer.From(props.EndpointType))

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

func (r PostgresqlFlexibleServerVirtualEndpointResource) Delete() sdk.ResourceFunc {
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
						if response.WasNotFound(resp.HttpResponse) {
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

func (r PostgresqlFlexibleServerVirtualEndpointResource) Update() sdk.ResourceFunc {
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
					_, err = client.Get(ctx, virtualEndpointId)
					if err != nil {
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
					EndpointType: pointer.To(virtualendpoints.VirtualEndpointType(virtualEndpoint.Type)),
					Members:      pointer.To([]string{replicaServerId.FlexibleServerName}),
				},
			}); err != nil {
				return fmt.Errorf("updating %q: %+v", id, err)
			}

			// the id has changed and needs to be updated
			replicaEndpointId := virtualendpoints.NewVirtualEndpointID(replicaServerId.SubscriptionId, replicaServerId.ResourceGroupName, replicaServerId.FlexibleServerName, virtualEndpoint.Name)
			endPointId := commonids.NewCompositeResourceID(&virtualEndpointId, &replicaEndpointId)
			metadata.SetID(endPointId)

			return nil
		},
	}
}

// The flexible endpoint API does not store the location/rg information on replicas it only stores the name.
// This lookup is safe because replicas for a given source server are *not* allowed to have identical names
func lookupFlexibleServerByName(ctx context.Context, flexibleServerClient *servers.ServersClient, virtualEndpointId virtualendpoints.VirtualEndpointId, replicaServerName string, sourceServerId string) (*servers.Server, error) {
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
