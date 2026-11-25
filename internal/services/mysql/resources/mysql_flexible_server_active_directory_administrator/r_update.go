// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mysql_flexible_server_active_directory_administrator

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/azureadadministrators"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mysql/parse"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MySQL.AzureADAdministratorsClient

			id, err := parse.FlexibleServerAzureActiveDirectoryAdministratorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model MySQLFlexibleServerAdministratorModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			flexibleServerId := azureadadministrators.NewFlexibleServerID(id.SubscriptionId, id.ResourceGroup, id.FlexibleServerName)

			resp, err := client.Get(ctx, flexibleServerId)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := resp.Model
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			if metadata.ResourceData.HasChange("identity_id") {
				properties.Properties.IdentityResourceId = pointer.To(model.IdentityId)
			}

			if metadata.ResourceData.HasChange("login") {
				properties.Properties.Login = pointer.To(model.Login)
			}

			if metadata.ResourceData.HasChange("object_id") {
				properties.Properties.Sid = pointer.To(model.ObjectId)
			}

			if metadata.ResourceData.HasChange("tenant_id") {
				properties.Properties.TenantId = pointer.To(model.TenantId)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, flexibleServerId, *properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
