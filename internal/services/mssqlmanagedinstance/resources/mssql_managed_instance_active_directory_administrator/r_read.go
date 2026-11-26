// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_managed_instance_active_directory_administrator

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssqlmanagedinstance/parse"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
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

			managedInstanceId := commonids.NewSqlManagedInstanceID(id.SubscriptionId, id.ResourceGroup, id.ManagedInstanceName)

			result, err := client.Get(ctx, managedInstanceId)
			if err != nil {
				if response.WasNotFound(result.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %v", id, err)
			}

			model := MsSqlManagedInstanceActiveDirectoryAdministratorModel{
				ManagedInstanceId: managedInstanceId.ID(),
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
