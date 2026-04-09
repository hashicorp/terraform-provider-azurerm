## File Storage Shares SDK for API version 2023-11-03

This package allows you to interact with the Shares File Storage API

### Supported Authorizers

* Azure Active Directory (for the Resource Endpoint `https://storage.azure.com`)
* SharedKeyLite (Blob, File & Queue)

### Example Usage

```go
package main

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/auth"
	"github.com/jackofallops/giovanni/storage/2023-11-03/file/shares"
)

func Example() error {
	accountName := "storageaccount1"
    storageAccountKey := "ABC123...."
    shareName := "myshare"
	domainSuffix := "core.windows.net"

	auth, err := auth.NewSharedKeyAuthorizer(accountName, storageAccountKey, auth.SharedKey)
	if err != nil {
		return fmt.Errorf("building SharedKey authorizer: %+v", err)
	}
	
    sharesClient, err := shares.NewWithBaseUri(fmt.Sprintf("https://%s.file.%s", accountName, domainSuffix))
	if err != nil {
		return fmt.Errorf("building SharedKey authorizer: %+v", err)
	}
    sharesClient.Client.SetAuthorizer(auth)
    
    ctx := context.TODO()
    input := shares.CreateInput{
    	QuotaInGB: 2,
    }
    if _, err := sharesClient.Create(ctx, shareName, input); err != nil {
        return fmt.Errorf("Error creating Share: %s", err)
    }
    
    return nil 
}
```