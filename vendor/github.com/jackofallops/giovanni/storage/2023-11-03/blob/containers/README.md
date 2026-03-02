## Blob Storage Container SDK for API version 2023-11-03

This package allows you to interact with the Containers Blob Storage API

### Supported Authorizers

* Azure Active Directory (for the Resource Endpoint `https://storage.azure.com`)
* SharedKeyLite (Blob, File & Queue)

Note: when using the `ListBlobs` operation, only `SharedKeyLite` authentication is supported.

### Example Usage

```go
package main

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/auth"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/containers"
)

func Example() error {
	accountName := "storageaccount1"
    storageAccountKey := "ABC123...."
    containerName := "mycontainer"
	domainSuffix := "core.windows.net"

    containersClient, err := containers.NewWithBaseUri(fmt.Sprintf("https://%s.blob.%s", accountName, domainSuffix))
	if err != nil {
		return fmt.Errorf("building client for environment: %+v", err)
	}

	auth, err := auth.NewSharedKeyAuthorizer(accountName, storageAccountKey, auth.SharedKey)
	if err != nil {
		return fmt.Errorf("building SharedKey authorizer: %+v", err)
	}
	containersClient.Client.SetAuthorizer(auth)
    
    ctx := context.TODO()
    createInput := containers.CreateInput{
        AccessLevel: containers.Private,
    }
    if _, err := containersClient.Create(ctx, containerName, createInput); err != nil {
        return fmt.Errorf("Error creating Container: %s", err)
    }
    
    return nil 
}
```