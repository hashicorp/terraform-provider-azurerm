## Queue Storage Queues SDK for API version 2019-12-12

This package allows you to interact with the Queues Queue Storage API

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
	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/queue/queues"
)

func Example() error {
	accountName := "storageaccount1"
    storageAccountKey := "ABC123...."
    queueName := "myqueue"
    
    storageAuth := autorest.NewSharedKeyLiteAuthorizer(accountName, storageAccountKey)
    queuesClient := queues.New()
    queuesClient.Client.Authorizer = storageAuth
    
    ctx := context.TODO()
    metadata := map[string]string{
    	"hello": "world",
    }
    if _, err := queuesClient.Create(ctx, accountName, queueName, metadata); err != nil {
        return fmt.Errorf("Error creating Queue: %s", err)
    }
    
    return nil 
}
```