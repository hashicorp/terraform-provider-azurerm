package storage

import (
	"context"
	"fmt"

	legacy "github.com/Azure/azure-sdk-for-go/storage"
)

// NOTE: we'll remove this once #4238 and #4235 have been merged to avoid conflicts

// nolint: deadcode unused
func (client *Client) LegacyBlobClient(ctx context.Context, resourceGroup, accountName string) (*legacy.BlobStorageClient, bool, error) {
	accountKey, err := client.findAccountKey(ctx, resourceGroup, accountName)
	if err != nil {
		return nil, false, fmt.Errorf("Error retrieving Account Key: %s", err)
	}

	apiVersion := legacy.DefaultAPIVersion
	storageClient, err := legacy.NewClient(accountName, *accountKey, client.environment.StorageEndpointSuffix, apiVersion, true)
	if err != nil {
		return nil, true, fmt.Errorf("Error creating storage client for storage account %q: %s", accountName, err)
	}

	blobClient := storageClient.GetBlobService()
	return &blobClient, true, nil
}
