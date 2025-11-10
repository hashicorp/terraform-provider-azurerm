## Data Lake Storage Gen2 File Systems SDK for API version 2023-11-03

This package allows you to interact with the Data Lake Storage Gen2 File Systems API

### Supported Authorizers

* Azure Active Directory (for the Resource Endpoint `https://storage.azure.com`)

### Example Usage

```go
package main

import (
	"context"
	"fmt"
	
    "github.com/hashicorp/go-azure-helpers/authentication"
    "github.com/hashicorp/go-azure-helpers/sender"
	"github.com/hashicorp/go-azure-sdk/sdk/auth"
    "github.com/jackofallops/giovanni/storage/2023-11-03/datalakestore/filesystems"
)

func Example() error {
	accountName := "storageaccount1"
    fileSystemName := "filesystem1"
	storageAccountKey := "ABC123...."
	domainSuffix := "core.windows.net"
	

	auth, err := auth.NewSharedKeyAuthorizer(accountName, storageAccountKey, auth.SharedKey)
	if err != nil {
		return fmt.Errorf("building SharedKey authorizer: %+v", err)
	}
	
    fileSystemsClient, err := filesystems.NewWithBaseUri(fmt.Sprintf("https://%s.dfs.%s", accountName, domainSuffix))
	if err != nil {
		return fmt.Errorf("building client for environment: %+v", err)
	}
	fileSystemsClient.Client.SetAuthorizer(auth)

	input := filesystems.CreateInput{
		Properties: map[string]string{},
	}

	ctx := context.Background()
	if _, err = fileSystemsClient.Create(ctx, fileSystemName, input); err != nil {
		return fmt.Errorf("Error creating: %s", err)
	}
	
    return nil 
}
```