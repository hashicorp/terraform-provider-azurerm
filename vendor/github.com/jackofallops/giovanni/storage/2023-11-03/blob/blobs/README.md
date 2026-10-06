## Blob Storage Blobs SDK for API version 2023-11-03

This package allows you to interact with the Blobs Blob Storage API

### Supported Authorizers

* Azure Active Directory (for the Resource Endpoint `https://storage.azure.com`)
* SharedKeyLite (Blob, File & Queue)

### Example Usage

```go
package main

import (
	"context"
	"fmt"
	"time"
	
	"github.com/hashicorp/go-azure-sdk/sdk/auth"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/blobs"
)

func Example() error {
	accountName := "storageaccount1"
    storageAccountKey := "ABC123...."
    containerName := "mycontainer"
	domainSuffix := "core.windows.net"
    fileName := "example-large-file.iso"

	blobClient, err := blobs.NewWithBaseUri(fmt.Sprintf("https://%s.blob.%s", accountName, domainSuffix))
	if err != nil {
		return fmt.Errorf("building client for environment: %v", err)
	}
	
	auth, err := auth.NewSharedKeyAuthorizer(accountName, storageAccountKey, auth.SharedKey)
	if err != nil {
		return fmt.Errorf("building SharedKey authorizer: %+v", err)
	}
	blobClient.Client.SetAuthorizer(auth)
    
    ctx := context.TODO()
    copyInput := blobs.CopyInput{
        CopySource: "http://releases.ubuntu.com/14.04/ubuntu-14.04.6-desktop-amd64.iso",
    }
    refreshInterval := 5 * time.Second
    if err := blobClient.CopyAndWait(ctx, containerName, fileName, copyInput, refreshInterval); err != nil {
        return fmt.Errorf("Error copying: %s", err)
    }
    
    return nil 
}

```