package mysql

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2022-01-01/azureadadministrators"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
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
	return validate.MySQLFlexibleServerAzureActiveDirectoryAdministratorID
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
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"login": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"object_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsUUID,
		},

		"tenant_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
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

			id := azureadadministrators.NewFlexibleServerID(flexibleServerId.SubscriptionId, flexibleServerId.ResourceGroupName, flexibleServerId.FlexibleServerName)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := &azureadadministrators.AzureADAdministrator{
				Properties: &azureadadministrators.AdministratorProperties{
					AdministratorType: pointer.To(azureadadministrators.AdministratorTypeActiveDirectory),
				},
			}

			if model.IdentityId != "" {
				properties.Properties.IdentityResourceId = &model.IdentityId
			}

			if model.Login != "" {
				properties.Properties.Login = &model.Login
			}

			if model.ObjectId != "" {
				properties.Properties.Sid = &model.ObjectId
			}

			if model.TenantId != "" {
				properties.Properties.TenantId = &model.TenantId
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, *properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

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

			id, err := azureadadministrators.ParseFlexibleServerID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model MySQLFlexibleServerAdministratorModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := resp.Model
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			if metadata.ResourceData.HasChange("identity_id") {
				if model.IdentityId != "" {
					properties.Properties.IdentityResourceId = &model.IdentityId
				} else {
					properties.Properties.IdentityResourceId = nil
				}
			}

			if metadata.ResourceData.HasChange("login") {
				if model.Login != "" {
					properties.Properties.Login = &model.Login
				} else {
					properties.Properties.Login = nil
				}
			}

			if metadata.ResourceData.HasChange("object_id") {
				if model.ObjectId != "" {
					properties.Properties.Sid = &model.ObjectId
				} else {
					properties.Properties.Sid = nil
				}
			}

			if metadata.ResourceData.HasChange("tenant_id") {
				if model.TenantId != "" {
					properties.Properties.TenantId = &model.TenantId
				} else {
					properties.Properties.TenantId = nil
				}
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *properties); err != nil {
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

			id, err := azureadadministrators.ParseFlexibleServerID(metadata.ResourceData.Id())
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

			state := MySQLFlexibleServerAdministratorModel{
				ServerId: azureadadministrators.NewFlexibleServerID(id.SubscriptionId, id.ResourceGroupName, id.FlexibleServerName).ID(),
			}

			if model := resp.Model; model != nil {
				if properties := model.Properties; properties != nil {
					if properties.IdentityResourceId != nil {
						state.IdentityId = pointer.From(properties.IdentityResourceId)
					}

					if properties.Login != nil {
						state.Login = pointer.From(properties.Login)
					}

					if properties.Sid != nil {
						state.ObjectId = pointer.From(properties.Sid)
					}

					if properties.TenantId != nil {
						state.TenantId = pointer.From(properties.TenantId)
					}
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

			id, err := azureadadministrators.ParseFlexibleServerID(metadata.ResourceData.Id())
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
