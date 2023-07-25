// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cosmos

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresqlhsc/2022-11-08/roles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type CosmosDbPostgreSQLRoleResourceModel struct {
	Name      string `tfschema:"name"`
	ClusterId string `tfschema:"cluster_id"`
	Password  string `tfschema:"password"`
}

type CosmosDbPostgreSQLRoleResource struct{}

var _ sdk.Resource = CosmosDbPostgreSQLClusterResource{}

func (r CosmosDbPostgreSQLRoleResource) ResourceType() string {
	return "azurerm_cosmosdb_postgresql_role"
}

func (r CosmosDbPostgreSQLRoleResource) ModelObject() interface{} {
	return &CosmosDbPostgreSQLRoleResourceModel{}
}

func (r CosmosDbPostgreSQLRoleResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return roles.ValidateRoleID
}

func (r CosmosDbPostgreSQLRoleResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.RoleName,
		},

		"cluster_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: roles.ValidateServerGroupsv2ID,
		},

		"password": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			Sensitive:    true,
			ValidateFunc: validate.RolePassword,
		},
	}
}

func (r CosmosDbPostgreSQLRoleResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r CosmosDbPostgreSQLRoleResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model CosmosDbPostgreSQLRoleResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Cosmos.RolesClient
			clusterId, err := roles.ParseServerGroupsv2ID(model.ClusterId)
			if err != nil {
				return err
			}

			id := roles.NewRoleID(clusterId.SubscriptionId, clusterId.ResourceGroupName, clusterId.ServerGroupsv2Name, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			parameters := roles.Role{
				Properties: roles.RoleProperties{
					Password: model.Password,
				},
			}

			if err := client.CreateThenPoll(ctx, id, parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r CosmosDbPostgreSQLRoleResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cosmos.RolesClient

			id, err := roles.ParseRoleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			state := CosmosDbPostgreSQLRoleResourceModel{
				Name:      id.RoleName,
				ClusterId: roles.NewServerGroupsv2ID(id.SubscriptionId, id.ResourceGroupName, id.ServerGroupsv2Name).ID(),
				Password:  metadata.ResourceData.Get("password").(string),
			}

			return metadata.Encode(&state)
		},
	}
}

func (r CosmosDbPostgreSQLRoleResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cosmos.RolesClient

			id, err := roles.ParseRoleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
