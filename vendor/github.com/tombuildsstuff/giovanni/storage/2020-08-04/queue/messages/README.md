## Queue Storage Messages SDK for API version 2020-08-04

This package allows you to interact with the Messages Queue Storage API

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
	"github.com/tombuildsstuff/giovanni/storage/2020-08-04/queue/messages"
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
    
    messagesClient, err  := messages.NewWithBaseUri(fmt.Sprintf("https://%s.queue.%s", accountName, domainSuffix))
	if err != nil {
		return fmt.Errorf("building client for environment: %+v", err)
	}
    messagesClient.Client.WithAuthorizer(auth)
    
    ctx := context.TODO()
    input := messages.PutInput{
    	Message: "<over><message>hello</message></over>",
    }
    if _, err := messagesClient.Put(ctx, queueName, input); err != nil {
        return fmt.Errorf("Error creating Message: %s", err)
    }
    
    return nil 
}
```