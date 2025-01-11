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
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2023-06-01-preview/servers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2023-06-01-preview/virtualendpoints"
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
					log.Printf("[INFO] %s does not exist - removing from state", metadata.ResourceData.Id())
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state.Name = id.VirtualEndpointName

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.Type = string(pointer.From(props.EndpointType))

					if props.Members == nil || len(*props.Members) != 2 {
						// if members list is nil, this is an endpoint that was previously deleted
						log.Printf("[INFO] Postgresql Flexible Server Endpoint %q was previously deleted - removing from state", id.ID())
						return metadata.MarkAsGone(id)
					}

					// Model.Properties.Members is a tuple => [source_server_id, replication_server_name]
					sourceServerName := (*props.Members)[0]
					replicaServerName := (*props.Members)[1]

					sourceServerId := servers.NewFlexibleServerID(id.SubscriptionId, id.ResourceGroupName, sourceServerName).ID()

					replicaServer, err := lookupFlexibleServerByName(ctx, flexibleServerClient, id, replicaServerName, sourceServerId)
					if err != nil {
						return err
					}

					state.SourceServerId = sourceServerId

					if replicaServer != nil {
						replicaId, err := servers.ParseFlexibleServerID(*replicaServer.Id)
						if err != nil {
							return err
						}

						state.ReplicaServerId = replicaId.ID()
					}
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

			replicaServerId, err := servers.ParseFlexibleServerID(virtualEndpoint.ReplicaServerId)
			if err != nil {
				return err
			}

			locks.ByName(id.FlexibleServerName, postgresqlFlexibleServerResourceName)
			defer locks.UnlockByName(id.FlexibleServerName, postgresqlFlexibleServerResourceName)

			if err := client.UpdateThenPoll(ctx, *id, virtualendpoints.VirtualEndpointResourceForPatch{
				Properties: &virtualendpoints.VirtualEndpointResourceProperties{
					EndpointType: pointer.To(virtualendpoints.VirtualEndpointType(virtualEndpoint.Type)),
					Members:      pointer.To([]string{replicaServerId.FlexibleServerName}),
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
