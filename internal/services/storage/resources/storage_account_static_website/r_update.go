// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_account_static_website

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/accounts"
)

func (a Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			storageClient := metadata.Client.Storage
			var model AccountStaticWebsiteResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id, err := commonids.ParseStorageAccountID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			accountDetails, err := storageClient.GetAccount(ctx, *id)
			if err != nil {
				return err
			}
			if accountDetails == nil {
				return fmt.Errorf("unable to locate %s", *id)
			}

			client, err := storageClient.AccountsDataPlaneClient(ctx, *accountDetails, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
			if err != nil {
				return fmt.Errorf("building Accounts Data Plane Client: %s", err)
			}

			props, err := client.GetServiceProperties(ctx, id.StorageAccountName)
			if err != nil || props.StaticWebsite == nil {
				return fmt.Errorf("retrieving static website properties for %s: %+v", *id, err)
			}

			properties := accounts.StorageServiceProperties{
				StaticWebsite: props.StaticWebsite,
			}

			if metadata.ResourceData.HasChange("index_document") {
				properties.StaticWebsite.IndexDocument = model.IndexDocument
			}

			if metadata.ResourceData.HasChange("error_404_document") {
				properties.StaticWebsite.ErrorDocument404Path = model.Error404Document
			}

			if _, err = client.SetServiceProperties(ctx, id.StorageAccountName, properties); err != nil {
				return fmt.Errorf("updating static website for %s: %+v", *id, err)
			}

			return nil
		},
	}
}
