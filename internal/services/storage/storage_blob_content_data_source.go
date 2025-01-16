package storage

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/accounts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/blobs"
)

const (
	mibSizeBytes           = 1048576
	maxBlobSizeMib         = 1
	maxBlobSizeBytes int64 = mibSizeBytes * maxBlobSizeMib
)

type StorageBlobContentDataSource struct{}

type StorageBlobContentDataSourceModel struct {
	Name                 string `tfschema:"name"`
	StorageAccountName   string `tfschema:"storage_account_name"`
	StorageContainerName string `tfschema:"storage_container_name"`
	Content              string `tfschema:"content"`
}

func (d StorageBlobContentDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"storage_account_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"storage_container_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
	}
}

func (d StorageBlobContentDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"content": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (d StorageBlobContentDataSource) ModelObject() interface{} {
	return &StorageBlobContentDataSourceModel{}
}

func (d StorageBlobContentDataSource) ResourceType() string {
	return "azurerm_storage_blob_content"
}

func (d StorageBlobContentDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Storage
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state StorageBlobContentDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			account, err := metadata.Client.Storage.FindAccount(ctx, subscriptionId, state.StorageAccountName)
			if err != nil {
				return fmt.Errorf("retrieving Storage Account %q: %v", state.StorageAccountName, err)
			}
			if account == nil {
				return fmt.Errorf("locating Storage Account %q", state.StorageAccountName)
			}

			endpoint, err := account.DataPlaneEndpoint("blob")
			if err != nil {
				return fmt.Errorf("determining Blob endpoint: %v", err)
			}

			accountId, err := accounts.ParseAccountID(*endpoint, client.StorageDomainSuffix)
			if err != nil {
				return fmt.Errorf("parsing Account ID: %v", err)
			}

			id := blobs.NewBlobID(*accountId, state.StorageContainerName, state.Name)
			log.Printf("[INFO] Retrieving %s", id)
			metadata.SetID(id)

			blobsClient, err := client.BlobsDataPlaneClient(ctx, *account, client.DataPlaneOperationSupportingAnyAuthMethod())
			if err != nil {
				return fmt.Errorf("building Blobs Client: %v", err)
			}

			propertiesInput := blobs.GetPropertiesInput{}
			properties, err := blobsClient.GetProperties(ctx, state.StorageContainerName, state.Name, propertiesInput)
			if err != nil {
				if response.WasNotFound(properties.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving content for %s: %v", id, err)
			}
			if properties.ContentLength > maxBlobSizeBytes {
				return fmt.Errorf("size of blob '%s' exceeds maximum size limit of %d bytes", state.Name, maxBlobSizeBytes)
			}

			input := blobs.GetInput{}
			content, err := blobsClient.Get(ctx, state.StorageContainerName, state.Name, input)
			if err != nil {
				if response.WasNotFound(content.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving content for %s: %v", id, err)
			}

			state.Content = base64.StdEncoding.EncodeToString(*content.Contents)
			return metadata.Encode(&state)
		},
	}
}
