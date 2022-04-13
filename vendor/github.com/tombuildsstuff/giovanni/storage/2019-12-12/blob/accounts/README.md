## Blob Storage Account SDK for API version 2019-12-12

This package allows you to interact with the Accounts Blob Storage API

### Supported Authorizers

* Azure Active Directory 

### Example Usage

```go
package main

import (
	"context"
	"fmt"
	"time"
	
	"github.com/Azure/go-autorest/autorest"
	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/blob/accounts"
)

func Example() error {
	accountName := "storageaccount1"
    
    // e.g. https://github.com/tombuildsstuff/giovanni/blob/76f5f686c99ecdcc3fa533a0330d0e1aacb1c327/example/azuread-auth/main.go#L54
    client, err := buildClient()
    if err != nil {
    	return fmt.Errorf("error building client: %s", err)
    }
    
    ctx := context.TODO()
    
    input := StorageServiceProperties{
        StaticWebsite: &StaticWebsite{
            Enabled:              true,
            IndexDocument:        index,
            ErrorDocument404Path: errorDocument,
        },
    }
    
    _, err = client.SetServiceProperties(ctx, accountName, input)
    if err != nil {
        return fmt.Errorf("error setting properties: %s", err)
    }
    
    time.Sleep(2 * time.Second)
    
    _, err = accountsClient.GetServiceProperties(ctx, accountName)
    if err != nil {
        return fmt.Errorf("error getting properties: %s", err)
    }
    
    return nil 
}

```