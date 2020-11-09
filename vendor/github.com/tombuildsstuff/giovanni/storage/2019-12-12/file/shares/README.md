## File Storage Shares SDK for API version 2019-12-12

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
	"time"
	
	"github.com/Azure/go-autorest/autorest"
	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/file/shares"
)

func Example() error {
	accountName := "storageaccount1"
    storageAccountKey := "ABC123...."
    shareName := "myshare"
    
    storageAuth := autorest.NewSharedKeyLiteAuthorizer(accountName, storageAccountKey)
    sharesClient := shares.New()
    sharesClient.Client.Authorizer = storageAuth
    
    ctx := context.TODO()
    input := shares.CreateInput{
    	QuotaInGB: 2,
    }
    if _, err := sharesClient.Create(ctx, accountName, shareName, input); err != nil {
        return fmt.Errorf("Error creating Share: %s", err)
    }
    
    return nil 
}
```