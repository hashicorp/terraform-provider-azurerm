## Queue Storage Queues SDK for API version 2023-11-03

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

	"github.com/hashicorp/go-azure-sdk/sdk/auth"
	"github.com/jackofallops/giovanni/storage/2023-11-03/queue/queues"
)

func Example() error {
	accountName := "storageaccount1"
    storageAccountKey := "ABC123...."
    queueName := "myqueue"
	domainSuffix := "core.windows.net"

	auth, err := auth.NewSharedKeyAuthorizer(accountName, storageAccountKey, auth.SharedKey)
	if err != nil {
		return fmt.Errorf("building SharedKey authorizer: %+v", err)
	}
	
    queuesClient, err := queues.NewWithBaseUri(fmt.Sprintf("https://%s.queue.%s", accountName, domainSuffix))
	if err != nil {
		return fmt.Errorf("building client for environment: %+v", err)
	}
    queuesClient.Client.SetAuthorizer(auth)
    
    ctx := context.TODO()
    metadata := map[string]string{
    	"hello": "world",
    }
	input := queues.CreateInput{
		Metadata: metadata,
    }
    if _, err := queuesClient.Create(ctx, queueName, input); err != nil {
        return fmt.Errorf("Error creating Queue: %s", err)
    }
    
    return nil 
}
```