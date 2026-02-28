// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2024-08-01/administrators"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2024-08-01/servers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type PostgresqlFlexibleServerActiveDirectoryAdministratorDataSourceModel struct {
	ServerId      string `tfschema:"server_id"`
	ObjectId      string `tfschema:"object_id"`
	PrincipalName string `tfschema:"principal_name"`
	PrincipalType string `tfschema:"principal_type"`
	TenantId      string `tfschema:"tenant_id"`
}

type PostgresqlFlexibleServerActiveDirectoryAdministratorDataSource struct{}

var _ sdk.DataSource = PostgresqlFlexibleServerActiveDirectoryAdministratorDataSource{}

func (d PostgresqlFlexibleServerActiveDirectoryAdministratorDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"server_id": commonschema.ResourceIDReferenceRequired(&servers.FlexibleServerId{}),

		"object_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.IsUUID,
		},
	}
}

func (d PostgresqlFlexibleServerActiveDirectoryAdministratorDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"principal_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"principal_type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tenant_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (d PostgresqlFlexibleServerActiveDirectoryAdministratorDataSource) ModelObject() interface{} {
	return &PostgresqlFlexibleServerActiveDirectoryAdministratorDataSourceModel{}
}

func (d PostgresqlFlexibleServerActiveDirectoryAdministratorDataSource) ResourceType() string {
	return "azurerm_postgresql_flexible_server_active_directory_administrator"
}

func (d PostgresqlFlexibleServerActiveDirectoryAdministratorDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Postgres.FlexibleServerAdministratorsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state PostgresqlFlexibleServerActiveDirectoryAdministratorDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			serverId, err := servers.ParseFlexibleServerID(state.ServerId)
			if err != nil {
				return err
			}

			id := administrators.NewAdministratorID(subscriptionId, serverId.ResourceGroupName, serverId.FlexibleServerName, state.ObjectId)

			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			metadata.SetID(id)

			state.ServerId = servers.NewFlexibleServerID(id.SubscriptionId, id.ResourceGroupName, id.FlexibleServerName).ID()
			state.ObjectId = id.ObjectId

			if model := resp.Model; model != nil {
				state.PrincipalName = pointer.From(model.Properties.PrincipalName)
				state.PrincipalType = string(pointer.From(model.Properties.PrincipalType))
				state.TenantId = pointer.From(model.Properties.TenantId)
			}

			return metadata.Encode(&state)
		},
	}
}
