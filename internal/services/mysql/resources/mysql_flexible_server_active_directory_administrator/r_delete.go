// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mysql_flexible_server_active_directory_administrator

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/azureadadministrators"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mysql/parse"
)

func (r Resource) Delete() sdk.ResourceFunc {
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
