// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package postgresql_flexible_server_virtual_endpoint

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2025-08-01/servers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2025-08-01/virtualendpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
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

			id := commonids.NewCompositeResourceID(&sourceEndpointId, &replicaEndpointId)

			payload := virtualendpoints.VirtualEndpoint{
				Name: &virtualEndpoint.Name,
				Properties: &virtualendpoints.VirtualEndpointResourceProperties{
					EndpointType: pointer.ToEnum[virtualendpoints.VirtualEndpointType](virtualEndpoint.Type),
					Members:      &[]string{replicaServerId.FlexibleServerName},
				},
			}

			if err = client.CreateCallbackThenPoll(ctx, sourceEndpointId, payload, metadata.SetIDCallback(id)); err != nil {
				return fmt.Errorf("creating %s: %+v", sourceEndpointId, err)
			}
			metadata.SetID(id)

			return nil
		},
	}
}
