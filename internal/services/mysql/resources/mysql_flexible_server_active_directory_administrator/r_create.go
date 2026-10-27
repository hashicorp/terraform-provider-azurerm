// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mysql_flexible_server_active_directory_administrator

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/azureadadministrators"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mysql/parse"
)

func (r Resource) Create() sdk.ResourceFunc {
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

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, *flexibleServerId)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", flexibleServerId, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), flexibleServerId)
				}
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

			// TODO: implement `CallbackThenPoll`, requires migrating to an ID that implements `resourceids.ResourceId`
			if err := client.CreateOrUpdateThenPoll(ctx, *flexibleServerId, properties); err != nil {
				return fmt.Errorf("creating %s: %+v", flexibleServerId, err)
			}

			metadata.SetID(parse.NewFlexibleServerAzureActiveDirectoryAdministratorID(flexibleServerId.SubscriptionId, flexibleServerId.ResourceGroupName, flexibleServerId.FlexibleServerName, string(azureadadministrators.AdministratorTypeActiveDirectory)))
			return nil
		},
	}
}
