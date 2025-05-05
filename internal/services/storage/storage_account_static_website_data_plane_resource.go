// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-05-01/storageaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/accounts"
)

type AccountStaticWebsiteResource struct{}

var _ sdk.ResourceWithUpdate = AccountStaticWebsiteResource{}

type AccountStaticWebsiteResourceModel struct {
	StorageAccountId string `tfschema:"storage_account_id"`
	Error404Document string `tfschema:"error_404_document"`
	IndexDocument    string `tfschema:"index_document"`
}

func (a AccountStaticWebsiteResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"storage_account_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateStorageAccountID,
		},

		"error_404_document": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			AtLeastOneOf: []string{"error_404_document", "index_document"},
		},

		"index_document": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.All(
				validation.StringDoesNotContainAny("/"),
				validation.StringLenBetween(3, 255),
			),
			AtLeastOneOf: []string{"error_404_document", "index_document"},
		},
	}
}

func (a AccountStaticWebsiteResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (a AccountStaticWebsiteResource) ModelObject() interface{} {
	return &AccountStaticWebsiteResourceModel{}
}

func (a AccountStaticWebsiteResource) ResourceType() string {
	return "azurerm_storage_account_static_website"
}

func (a AccountStaticWebsiteResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return commonids.ValidateStorageAccountID
}

func (a AccountStaticWebsiteResource) Create() sdk.ResourceFunc {
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

			supportLevel := availableFunctionalityForAccount(accountDetails.Kind, accountTier, accountReplicationType)

			if !supportLevel.supportStaticWebsite {
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

			return nil
		},
	}
}

func (a AccountStaticWebsiteResource) Read() sdk.ResourceFunc {
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

			return metadata.Encode(&state)
		},
	}
}

func (a AccountStaticWebsiteResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			storageClient := metadata.Client.Storage

			id, err := commonids.ParseStorageAccountID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			accountDetails, err := storageClient.GetAccount(ctx, *id)
			if err != nil {
				return nil // lint:ignore nilerr If we don't find the account we can safely assume we don't need to remove the website since it must already be deleted
			}

			properties := accounts.StorageServiceProperties{
				StaticWebsite: &accounts.StaticWebsite{
					Enabled: false,
				},
			}

			client, err := storageClient.AccountsDataPlaneClient(ctx, *accountDetails, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
			if err != nil {
				return fmt.Errorf("building Accounts Data Plane Client: %s", err)
			}

			if _, err = client.SetServiceProperties(ctx, id.StorageAccountName, properties); err != nil {
				return fmt.Errorf("deleting static website for %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (a AccountStaticWebsiteResource) Update() sdk.ResourceFunc {
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
