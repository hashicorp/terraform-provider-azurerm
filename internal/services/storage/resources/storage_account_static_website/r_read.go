// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_account_static_website

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (a Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			storageClient := metadata.Client.Storage

			var state AccountStaticWebsiteResourceModel

			id, err := commonids.ParseStorageAccountID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			state.StorageAccountId = id.ID()

			accountDetails, err := storageClient.GetAccount(ctx, *id)
			if err != nil {
				return metadata.MarkAsGone(id)
			}

			accountsClient, err := storageClient.AccountsDataPlaneClient(ctx, *accountDetails, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
			if err != nil {
				return fmt.Errorf("building Accounts Data Plane Client for %s: %+v", *id, err)
			}

			props, err := accountsClient.GetServiceProperties(ctx, id.StorageAccountName)
			if err != nil {
				return fmt.Errorf("retrieving static website properties for %s: %+v", *id, err)
			}

			if website := props.StaticWebsite; website != nil {
				state.IndexDocument = website.IndexDocument
				state.Error404Document = website.ErrorDocument404Path
			}

			if err = pluginsdk.SetResourceIdentityData(metadata.ResourceData, id, pluginsdk.ResourceTypeForIdentityVirtual); err != nil {
				return err
			}

			return metadata.Encode(&state)
		},
	}
}
