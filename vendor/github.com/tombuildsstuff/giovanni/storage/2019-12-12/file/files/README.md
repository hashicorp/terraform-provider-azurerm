## File Storage Files SDK for API version 2019-12-12

This package allows you to interact with the Files File Storage API

### Supported Authorizers

* Azure Active Directory (for the Resource Endpoint `https://storage.azure.com`)
* SharedKeyLite (Blob, File & Queue)

### Limitations

* At this time the headers `x-ms-file-permission` and `x-ms-file-attributes` are hard-coded (to `inherit` and `None`, respectively).

### Example Usage

```go
package main

import (
	"context"
	"fmt"
	"time"
	
	"github.com/Azure/go-autorest/autorest"
	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/file/files"
)

func Example() error {
	accountName := "storageaccount1"
    storageAccountKey := "ABC123...."
    shareName := "myshare"
    directoryName := "myfiles"
    fileName := "example.txt"
    
    storageAuth := autorest.NewSharedKeyLiteAuthorizer(accountName, storageAccountKey)
    filesClient := files.New()
    filesClient.Client.Authorizer = storageAuth
    
    ctx := context.TODO()
    input := files.CreateInput{}
    if _, err := filesClient.Create(ctx, accountName, shareName, directoryName, fileName, input); err != nil {
        return fmt.Errorf("Error creating File: %s", err)
    }
    
    return nil 
}
```