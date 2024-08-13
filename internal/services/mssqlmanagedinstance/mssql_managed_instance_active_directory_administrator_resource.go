// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssqlmanagedinstance

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/managedinstanceadministrators"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/managedinstanceazureadonlyauthentications"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssqlmanagedinstance/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssqlmanagedinstance/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type MsSqlManagedInstanceActiveDirectoryAdministratorModel struct {
	ManagedInstanceId         string `tfschema:"managed_instance_id"`
	AzureADAuthenticationOnly bool   `tfschema:"azuread_authentication_only"`
	LoginUsername             string `tfschema:"login_username"`
	ObjectId                  string `tfschema:"object_id"`
	TenantId                  string `tfschema:"tenant_id"`
}

var _ sdk.Resource = MsSqlManagedInstanceActiveDirectoryAdministratorResource{}
var _ sdk.ResourceWithUpdate = MsSqlManagedInstanceActiveDirectoryAdministratorResource{}

type MsSqlManagedInstanceActiveDirectoryAdministratorResource struct{}

func (r MsSqlManagedInstanceActiveDirectoryAdministratorResource) ResourceType() string {
	return "azurerm_mssql_managed_instance_active_directory_administrator"
}

func (r MsSqlManagedInstanceActiveDirectoryAdministratorResource) ModelObject() interface{} {
	return &MsSqlManagedInstanceActiveDirectoryAdministratorModel{}
}

func (r MsSqlManagedInstanceActiveDirectoryAdministratorResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.ManagedInstanceAzureActiveDirectoryAdministratorID
}

func (r MsSqlManagedInstanceActiveDirectoryAdministratorResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"managed_instance_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ManagedInstanceID,
		},

		"login_username": {
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

		"azuread_authentication_only": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},
	}
}

func (r MsSqlManagedInstanceActiveDirectoryAdministratorResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r MsSqlManagedInstanceActiveDirectoryAdministratorResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQLManagedInstance.ManagedInstanceAdministratorsClient
			aadAuthOnlyClient := metadata.Client.MSSQLManagedInstance.ManagedInstanceAzureADOnlyAuthenticationsClient

			var model MsSqlManagedInstanceActiveDirectoryAdministratorModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			managedInstanceId, err := commonids.ParseSqlManagedInstanceID(model.ManagedInstanceId)
			if err != nil {
				return err
			}

			id := parse.NewManagedInstanceAzureActiveDirectoryAdministratorID(managedInstanceId.SubscriptionId,
				managedInstanceId.ResourceGroupName, managedInstanceId.ManagedInstanceName, string(managedinstanceadministrators.ManagedInstanceAdministratorTypeActiveDirectory))

			metadata.Logger.Infof("Import check for %s", id)
			existing, err := client.Get(ctx, *managedInstanceId)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			parameters := managedinstanceadministrators.ManagedInstanceAdministrator{
				Properties: &managedinstanceadministrators.ManagedInstanceAdministratorProperties{
					AdministratorType: managedinstanceadministrators.ManagedInstanceAdministratorTypeActiveDirectory,
					Login:             model.LoginUsername,
					Sid:               model.ObjectId,
					TenantId:          &model.TenantId,
				},
			}

			metadata.Logger.Infof("Creating %s", id)

			err = client.CreateOrUpdateThenPoll(ctx, *managedInstanceId, parameters)
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			aadAuthOnlyParams := managedinstanceazureadonlyauthentications.ManagedInstanceAzureADOnlyAuthentication{
				Properties: &managedinstanceazureadonlyauthentications.ManagedInstanceAzureADOnlyAuthProperties{
					AzureADOnlyAuthentication: model.AzureADAuthenticationOnly,
				},
			}

			err = aadAuthOnlyClient.CreateOrUpdateThenPoll(ctx, *managedInstanceId, aadAuthOnlyParams)
			if err != nil {
				return fmt.Errorf("setting `azuread_authentication_only` for %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r MsSqlManagedInstanceActiveDirectoryAdministratorResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQLManagedInstance.ManagedInstanceAdministratorsClient
			aadAuthOnlyClient := metadata.Client.MSSQLManagedInstance.ManagedInstanceAzureADOnlyAuthenticationsClient

			id, err := parse.ManagedInstanceAzureActiveDirectoryAdministratorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("Decoding state for %s", id)
			var state MsSqlManagedInstanceActiveDirectoryAdministratorModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			managedInstanceId, err := commonids.ParseSqlManagedInstanceID(state.ManagedInstanceId)
			if err != nil {
				return err
			}

			parameters := managedinstanceadministrators.ManagedInstanceAdministrator{
				Properties: &managedinstanceadministrators.ManagedInstanceAdministratorProperties{
					AdministratorType: managedinstanceadministrators.ManagedInstanceAdministratorTypeActiveDirectory,
					Login:             state.LoginUsername,
					Sid:               state.ObjectId,
					TenantId:          &state.TenantId,
				},
			}

			metadata.Logger.Infof("Updating %s", id)

			err = client.CreateOrUpdateThenPoll(ctx, *managedInstanceId, parameters)
			if err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			aadAuthOnlyProperties := managedinstanceazureadonlyauthentications.ManagedInstanceAzureADOnlyAuthentication{
				Properties: &managedinstanceazureadonlyauthentications.ManagedInstanceAzureADOnlyAuthProperties{
					AzureADOnlyAuthentication: state.AzureADAuthenticationOnly,
				},
			}

			err = aadAuthOnlyClient.CreateOrUpdateThenPoll(ctx, *managedInstanceId, aadAuthOnlyProperties)
			if err != nil {
				return fmt.Errorf("setting `azuread_authentication_only` for %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r MsSqlManagedInstanceActiveDirectoryAdministratorResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQLManagedInstance.ManagedInstanceAdministratorsClient
			aadAuthOnlyClient := metadata.Client.MSSQLManagedInstance.ManagedInstanceAzureADOnlyAuthenticationsClient

			id, err := parse.ManagedInstanceAzureActiveDirectoryAdministratorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("Decoding state for %s", id)
			var state MsSqlManagedInstanceActiveDirectoryAdministratorModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			managedInstanceId := commonids.NewSqlManagedInstanceID(id.SubscriptionId, id.ResourceGroup, id.ManagedInstanceName)

			result, err := client.Get(ctx, managedInstanceId)
			if err != nil {
				if response.WasNotFound(result.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %v", id, err)
			}

			model := MsSqlManagedInstanceActiveDirectoryAdministratorModel{
				ManagedInstanceId:         managedInstanceId.ID(),
				AzureADAuthenticationOnly: false,
			}

			if result.Model != nil {

				if props := result.Model.Properties; props != nil {
					model.LoginUsername = props.Login
					model.ObjectId = props.Sid
					model.TenantId = pointer.From(props.TenantId)

				}
			}

			aadAuthOnlyResult, err := aadAuthOnlyClient.Get(ctx, managedInstanceId)
			if err != nil && !response.WasNotFound(result.HttpResponse) {
				return fmt.Errorf("retrieving `azuread_authentication_only` for %s: %v", id, err)
			}

			if aadAuthOnlyModel := aadAuthOnlyResult.Model; aadAuthOnlyModel != nil {
				if props := aadAuthOnlyModel.Properties; props != nil {
					model.AzureADAuthenticationOnly = props.AzureADOnlyAuthentication
				}
			}

			return metadata.Encode(&model)
		},
	}
}

func (r MsSqlManagedInstanceActiveDirectoryAdministratorResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 180 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQLManagedInstance.ManagedInstanceAdministratorsClient
			aadAuthOnlyClient := metadata.Client.MSSQLManagedInstance.ManagedInstanceAzureADOnlyAuthenticationsClient

			id, err := parse.ManagedInstanceAzureActiveDirectoryAdministratorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			managedInstanceId := commonids.NewSqlManagedInstanceID(id.SubscriptionId, id.ResourceGroup, id.ManagedInstanceName)

			err = aadAuthOnlyClient.DeleteThenPoll(ctx, managedInstanceId)
			if err != nil {
				return fmt.Errorf("removing `azuread_authentication_only` for %s: %+v", managedInstanceId, err)
			}

			err = client.DeleteThenPoll(ctx, managedInstanceId)
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
