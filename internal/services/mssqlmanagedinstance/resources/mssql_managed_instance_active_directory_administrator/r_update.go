// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_managed_instance_active_directory_administrator

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/managedinstanceadministrators"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/managedinstanceazureadonlyauthentications"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssqlmanagedinstance/parse"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQLManagedInstance.ManagedInstanceAdministratorsClient
			aadAuthOnlyClient := metadata.Client.MSSQLManagedInstance.ManagedInstanceAzureADOnlyAuthenticationsClient

			id, err := parse.ManagedInstanceAzureActiveDirectoryAdministratorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

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

			if err = client.CreateOrUpdateThenPoll(ctx, *managedInstanceId, parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			aadAuthOnlyProperties := managedinstanceazureadonlyauthentications.ManagedInstanceAzureADOnlyAuthentication{
				Properties: &managedinstanceazureadonlyauthentications.ManagedInstanceAzureADOnlyAuthProperties{
					AzureADOnlyAuthentication: state.AzureADAuthenticationOnly,
				},
			}

			if err = aadAuthOnlyClient.CreateOrUpdateThenPoll(ctx, *managedInstanceId, aadAuthOnlyProperties); err != nil {
				return fmt.Errorf("setting `azuread_authentication_only` for %s: %+v", id, err)
			}

			return nil
		},
	}
}
