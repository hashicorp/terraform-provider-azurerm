// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_managed_instance_active_directory_administrator

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/managedinstanceadministrators"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/managedinstanceazureadonlyauthentications"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssqlmanagedinstance/parse"
)

func (r Resource) Create() sdk.ResourceFunc {
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

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, *managedInstanceId)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			parameters := managedinstanceadministrators.ManagedInstanceAdministrator{
				Properties: &managedinstanceadministrators.ManagedInstanceAdministratorProperties{
					AdministratorType: managedinstanceadministrators.ManagedInstanceAdministratorTypeActiveDirectory,
					Login:             model.LoginUsername,
					Sid:               model.ObjectId,
					TenantId:          &model.TenantId,
				},
			}

			// TODO: implement `CallbackThenPoll`, requires migrating to an ID that implements `resourceids.ResourceId`
			if err := client.CreateOrUpdateThenPoll(ctx, *managedInstanceId, parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
			metadata.SetID(id)

			aadAuthOnlyParams := managedinstanceazureadonlyauthentications.ManagedInstanceAzureADOnlyAuthentication{
				Properties: &managedinstanceazureadonlyauthentications.ManagedInstanceAzureADOnlyAuthProperties{
					AzureADOnlyAuthentication: model.AzureADAuthenticationOnly,
				},
			}

			if err = aadAuthOnlyClient.CreateOrUpdateThenPoll(ctx, *managedInstanceId, aadAuthOnlyParams); err != nil {
				return fmt.Errorf("setting `azuread_authentication_only` for %s: %+v", id, err)
			}

			return nil
		},
	}
}
