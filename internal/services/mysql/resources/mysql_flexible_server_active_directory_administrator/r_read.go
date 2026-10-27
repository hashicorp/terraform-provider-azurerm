// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mysql_flexible_server_active_directory_administrator

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
)

func (r Resource) Read() sdk.ResourceFunc {
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
