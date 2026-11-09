// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type AccountStaticWebsiteDataSource struct{}

var _ sdk.DataSource = AccountStaticWebsiteDataSource{}

type AccountStaticWebsiteDataSourceModel struct {
	StorageAccountId string `tfschema:"storage_account_id"`
	Error404Document string `tfschema:"error_404_document"`
	IndexDocument    string `tfschema:"index_document"`
}

func (d AccountStaticWebsiteDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"storage_account_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: commonids.ValidateStorageAccountID,
		},
	}
}

func (d AccountStaticWebsiteDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"error_404_document": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"index_document": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (d AccountStaticWebsiteDataSource) ModelObject() any {
	return &AccountStaticWebsiteDataSourceModel{}
}

func (d AccountStaticWebsiteDataSource) ResourceType() string {
	return "azurerm_storage_account_static_website"
}

func (d AccountStaticWebsiteDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Storage

			var state AccountStaticWebsiteDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id, err := commonids.ParseStorageAccountID(state.StorageAccountId)
			if err != nil {
				return err
			}

			accountDetails, err := client.GetAccount(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			accountsClient, err := client.AccountsDataPlaneClient(ctx, *accountDetails, client.DataPlaneOperationSupportingAnyAuthMethod())
			if err != nil {
				return fmt.Errorf("building Accounts Data Plane Client for %s: %+v", *id, err)
			}

			props, err := accountsClient.GetServiceProperties(ctx, id.StorageAccountName)
			if err != nil {
				return fmt.Errorf("retrieving static website properties for %s: %+v", *id, err)
			}

			website := props.StaticWebsite
			if website == nil || !website.Enabled {
				return fmt.Errorf("static website is not enabled for %s", *id)
			}

			state.IndexDocument = website.IndexDocument
			state.Error404Document = website.ErrorDocument404Path

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
