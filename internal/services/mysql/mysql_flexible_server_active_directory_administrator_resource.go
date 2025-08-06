// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mysql

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/azureadadministrators"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mysql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mysql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type MySQLFlexibleServerAdministratorModel struct {
	ServerId   string `tfschema:"server_id"`
	IdentityId string `tfschema:"identity_id"`
	Login      string `tfschema:"login"`
	ObjectId   string `tfschema:"object_id"`
	TenantId   string `tfschema:"tenant_id"`
}

type MySQLFlexibleServerAdministratorResource struct{}

var _ sdk.ResourceWithUpdate = MySQLFlexibleServerAdministratorResource{}

func (r MySQLFlexibleServerAdministratorResource) ResourceType() string {
	return "azurerm_mysql_flexible_server_active_directory_administrator"
}

func (r MySQLFlexibleServerAdministratorResource) ModelObject() interface{} {
	return &MySQLFlexibleServerAdministratorModel{}
}

func (r MySQLFlexibleServerAdministratorResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.FlexibleServerAzureActiveDirectoryAdministratorID
}

func (r MySQLFlexibleServerAdministratorResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"server_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: azureadadministrators.ValidateFlexibleServerID,
		},

		"identity_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: commonids.ValidateUserAssignedIdentityID,
		},

		"login": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"object_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.IsUUID,
		},

		"tenant_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.IsUUID,
		},
	}
}

func (r MySQLFlexibleServerAdministratorResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r MySQLFlexibleServerAdministratorResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model MySQLFlexibleServerAdministratorModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.MySQL.AzureADAdministratorsClient
			flexibleServerId, err := azureadadministrators.ParseFlexibleServerID(model.ServerId)
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *flexibleServerId)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", flexibleServerId, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), flexibleServerId)
			}

			properties := azureadadministrators.AzureADAdministrator{
				Properties: &azureadadministrators.AdministratorProperties{
					AdministratorType:  pointer.To(azureadadministrators.AdministratorTypeActiveDirectory),
					IdentityResourceId: pointer.To(model.IdentityId),
					Login:              pointer.To(model.Login),
					Sid:                pointer.To(model.ObjectId),
					TenantId:           pointer.To(model.TenantId),
				},
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *flexibleServerId, properties); err != nil {
				return fmt.Errorf("creating %s: %+v", flexibleServerId, err)
			}

			id := parse.NewFlexibleServerAzureActiveDirectoryAdministratorID(flexibleServerId.SubscriptionId, flexibleServerId.ResourceGroupName, flexibleServerId.FlexibleServerName, string(azureadadministrators.AdministratorTypeActiveDirectory))

			metadata.SetID(id)
			return nil
		},
	}
}

func (r MySQLFlexibleServerAdministratorResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MySQL.AzureADAdministratorsClient

			id, err := parse.FlexibleServerAzureActiveDirectoryAdministratorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model MySQLFlexibleServerAdministratorModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			flexibleServerId := azureadadministrators.NewFlexibleServerID(id.SubscriptionId, id.ResourceGroup, id.FlexibleServerName)

			resp, err := client.Get(ctx, flexibleServerId)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := resp.Model
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			if metadata.ResourceData.HasChange("identity_id") {
				properties.Properties.IdentityResourceId = pointer.To(model.IdentityId)
			}

			if metadata.ResourceData.HasChange("login") {
				properties.Properties.Login = pointer.To(model.Login)
			}

			if metadata.ResourceData.HasChange("object_id") {
				properties.Properties.Sid = pointer.To(model.ObjectId)
			}

			if metadata.ResourceData.HasChange("tenant_id") {
				properties.Properties.TenantId = pointer.To(model.TenantId)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, flexibleServerId, *properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r MySQLFlexibleServerAdministratorResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MySQL.AzureADAdministratorsClient

			id, err := parse.FlexibleServerAzureActiveDirectoryAdministratorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			flexibleServerId := azureadadministrators.NewFlexibleServerID(id.SubscriptionId, id.ResourceGroup, id.FlexibleServerName)

			resp, err := client.Get(ctx, flexibleServerId)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := MySQLFlexibleServerAdministratorModel{
				ServerId: flexibleServerId.ID(),
			}

			if model := resp.Model; model != nil {
				if properties := model.Properties; properties != nil {
					identity, err := commonids.ParseUserAssignedIdentityIDInsensitively(*properties.IdentityResourceId)
					if err != nil {
						return err
					}

					state.IdentityId = identity.ID()
					state.Login = pointer.From(properties.Login)
					state.ObjectId = pointer.From(properties.Sid)
					state.TenantId = pointer.From(properties.TenantId)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r MySQLFlexibleServerAdministratorResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MySQL.AzureADAdministratorsClient

			id, err := parse.FlexibleServerAzureActiveDirectoryAdministratorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			flexibleServerId := azureadadministrators.NewFlexibleServerID(id.SubscriptionId, id.ResourceGroup, id.FlexibleServerName)

			if err := client.DeleteThenPoll(ctx, flexibleServerId); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
