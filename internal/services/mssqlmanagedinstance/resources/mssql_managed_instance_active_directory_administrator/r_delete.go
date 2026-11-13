// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_managed_instance_active_directory_administrator

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/managedinstanceazureadonlyauthentications"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssqlmanagedinstance/parse"
)

func (r Resource) Delete() sdk.ResourceFunc {
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

			// Before deleting an AAD admin, it is necessary to disable `AzureADOnlyAuthentication` first, as deleting an AAD admin when `AzureADOnlyAuthentication` feature is enabled is not supported.
			// Use `CreateOrUpdateThenPoll` instead of `DeleteThenPoll`, because the actual deletion behaviour of the API is not to really delete the record, but to update `AzureADOnlyAuthentication` to false. Therefore, using `DeleteThenPoll` will cause pull till done to never end until it times out.
			aadAuthOnlyParams := managedinstanceazureadonlyauthentications.ManagedInstanceAzureADOnlyAuthentication{
				Properties: &managedinstanceazureadonlyauthentications.ManagedInstanceAzureADOnlyAuthProperties{},
			}

			if err = aadAuthOnlyClient.CreateOrUpdateThenPoll(ctx, managedInstanceId, aadAuthOnlyParams); err != nil {
				return fmt.Errorf("disabling `azuread_authentication_only` for %s: %+v", id, err)
			}

			if err = client.DeleteThenPoll(ctx, managedInstanceId); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
