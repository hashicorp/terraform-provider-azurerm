// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssqlmanagedinstance

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v5.0/sql" // nolint: staticcheck
	"github.com/gofrs/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssqlmanagedinstance/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
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

			managedInstanceId, err := parse.ManagedInstanceID(model.ManagedInstanceId)
			if err != nil {
				return fmt.Errorf("parsing `managed_instance_id`: %v", err)
			}

			id := parse.NewManagedInstanceAzureActiveDirectoryAdministratorID(managedInstanceId.SubscriptionId,
				managedInstanceId.ResourceGroup, managedInstanceId.Name, string(sql.AdministratorTypeActiveDirectory))

			metadata.Logger.Infof("Import check for %s", id)
			existing, err := client.Get(ctx, id.ResourceGroup, id.ManagedInstanceName)
			if err != nil && !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if !utils.ResponseWasNotFound(existing.Response) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			sid, err := uuid.FromString(model.ObjectId)
			if err != nil {
				return fmt.Errorf("parsing `object_id` for %s", id)
			}

			tid, err := uuid.FromString(model.TenantId)
			if err != nil {
				return fmt.Errorf("parsing `tenant_id` for %s", id)
			}

			parameters := sql.ManagedInstanceAdministrator{
				ManagedInstanceAdministratorProperties: &sql.ManagedInstanceAdministratorProperties{
					AdministratorType: utils.String(string(sql.AdministratorTypeActiveDirectory)),
					Login:             &model.LoginUsername,
					Sid:               &sid,
					TenantID:          &tid,
				},
			}

			metadata.Logger.Infof("Creating %s", id)

			future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ManagedInstanceName, parameters)
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for creation of %s: %+v", id, err)
			}

			metadata.SetID(id)

			aadAuthOnlyParams := sql.ManagedInstanceAzureADOnlyAuthentication{
				ManagedInstanceAzureADOnlyAuthProperties: &sql.ManagedInstanceAzureADOnlyAuthProperties{
					AzureADOnlyAuthentication: &model.AzureADAuthenticationOnly,
				},
			}

			aadAuthOnlyFuture, err := aadAuthOnlyClient.CreateOrUpdate(ctx, id.ResourceGroup, id.ManagedInstanceName, aadAuthOnlyParams)
			if err != nil {
				return fmt.Errorf("setting `azuread_authentication_only` for %s: %+v", id, err)
			}

			if err = aadAuthOnlyFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting to set `azuread_authentication_only` for %s: %+v", id, err)
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

			sid, err := uuid.FromString(state.ObjectId)
			if err != nil {
				return fmt.Errorf("parsing `object_id` for %s", id)
			}

			tid, err := uuid.FromString(state.TenantId)
			if err != nil {
				return fmt.Errorf("parsing `tenant_id` for %s", id)
			}

			properties := sql.ManagedInstanceAdministrator{
				ManagedInstanceAdministratorProperties: &sql.ManagedInstanceAdministratorProperties{
					AdministratorType: utils.String(string(sql.AdministratorTypeActiveDirectory)),
					Login:             &state.LoginUsername,
					Sid:               &sid,
					TenantID:          &tid,
				},
			}

			metadata.Logger.Infof("Updating %s", id)

			future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ManagedInstanceName, properties)
			if err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for update of %s: %+v", id, err)
			}

			aadAuthOnlyProperties := sql.ManagedInstanceAzureADOnlyAuthentication{
				ManagedInstanceAzureADOnlyAuthProperties: &sql.ManagedInstanceAzureADOnlyAuthProperties{
					AzureADOnlyAuthentication: &state.AzureADAuthenticationOnly,
				},
			}

			aadAuthOnlyFuture, err := aadAuthOnlyClient.CreateOrUpdate(ctx, id.ResourceGroup, id.ManagedInstanceName, aadAuthOnlyProperties)
			if err != nil {
				return fmt.Errorf("setting `azuread_authentication_only` for %s: %+v", id, err)
			}

			if err = aadAuthOnlyFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting to set `azuread_authentication_only` for %s: %+v", id, err)
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

			result, err := client.Get(ctx, id.ResourceGroup, id.ManagedInstanceName)
			if err != nil {
				if utils.ResponseWasNotFound(result.Response) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %v", id, err)
			}

			instanceId := parse.NewManagedInstanceID(id.SubscriptionId, id.ResourceGroup, id.ManagedInstanceName)

			model := MsSqlManagedInstanceActiveDirectoryAdministratorModel{
				ManagedInstanceId:         instanceId.ID(),
				AzureADAuthenticationOnly: false,
			}

			if props := result.ManagedInstanceAdministratorProperties; props != nil {
				if props.Login != nil {
					model.LoginUsername = *props.Login
				}
				if props.Sid != nil {
					model.ObjectId = props.Sid.String()
				}
				if props.TenantID != nil {
					model.TenantId = props.TenantID.String()
				}
			}

			aadAuthOnlyResult, err := aadAuthOnlyClient.Get(ctx, id.ResourceGroup, id.ManagedInstanceName)
			if err != nil && !utils.ResponseWasNotFound(result.Response) {
				return fmt.Errorf("retrieving `azuread_authentication_only` for %s: %v", id, err)
			}

			if props := aadAuthOnlyResult.ManagedInstanceAzureADOnlyAuthProperties; props != nil && props.AzureADOnlyAuthentication != nil {
				model.AzureADAuthenticationOnly = *props.AzureADOnlyAuthentication
			}

			return metadata.Encode(&model)
		},
	}
}

func (r MsSqlManagedInstanceActiveDirectoryAdministratorResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQLManagedInstance.ManagedInstanceAdministratorsClient

			id, err := parse.ManagedInstanceAzureActiveDirectoryAdministratorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			future, err := client.Delete(ctx, id.ResourceGroup, id.ManagedInstanceName)
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for deletion of %s: %+v", id, err)
			}

			return nil
		},
	}
}
