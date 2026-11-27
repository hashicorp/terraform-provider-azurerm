// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_account_static_website

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/storageaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/accounts"
)

func (a Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			storageClient := metadata.Client.Storage
			var model AccountStaticWebsiteResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			accountID, err := commonids.ParseStorageAccountID(model.StorageAccountId)
			if err != nil {
				return err
			}

			// Get the target account to ensure it supports queues
			account, err := storageClient.ResourceManager.StorageAccounts.GetProperties(ctx, *accountID, storageaccounts.DefaultGetPropertiesOperationOptions())
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *accountID, err)
			}
			if account.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *accountID)
			}

			if account.Model.Sku == nil || account.Model.Sku.Tier == nil || string(account.Model.Sku.Name) == "" {
				return fmt.Errorf("could not read SKU details for %s", *accountID)
			}

			accountTier := *account.Model.Sku.Tier
			accountReplicationTypeParts := strings.Split(string(account.Model.Sku.Name), "_")
			if len(accountReplicationTypeParts) != 2 {
				return fmt.Errorf("could not read SKU replication type for %s", *accountID)
			}
			accountReplicationType := accountReplicationTypeParts[1]

			accountDetails, err := storageClient.GetAccount(ctx, *accountID)
			if err != nil {
				return err
			}
			if accountDetails == nil {
				return fmt.Errorf("unable to locate %s", *accountID)
			}

			supportLevel := helpers.AvailableFunctionalityForAccount(accountDetails.Kind, accountTier, accountReplicationType)

			if !supportLevel.SupportStaticWebsite {
				return fmt.Errorf("account %s does not support Static Websites", *accountID)
			}

			client, err := storageClient.AccountsDataPlaneClient(ctx, *accountDetails, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
			if err != nil {
				return fmt.Errorf("building Accounts Data Plane Client: %s", err)
			}

			properties := accounts.StorageServiceProperties{
				StaticWebsite: &accounts.StaticWebsite{
					Enabled: true,
				},
			}
			if model.IndexDocument != "" {
				properties.StaticWebsite.IndexDocument = model.IndexDocument
			}
			if model.Error404Document != "" {
				properties.StaticWebsite.ErrorDocument404Path = model.Error404Document
			}

			if _, err = client.SetServiceProperties(ctx, accountID.StorageAccountName, properties); err != nil {
				return fmt.Errorf("creating static website for %s: %+v", accountID, err)
			}

			metadata.SetID(accountID)
			return pluginsdk.SetResourceIdentityData(metadata.ResourceData, accountID, pluginsdk.ResourceTypeForIdentityVirtual)
		},
	}
}
