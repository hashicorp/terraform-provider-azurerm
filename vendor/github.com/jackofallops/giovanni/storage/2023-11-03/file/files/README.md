## File Storage Files SDK for API version 2023-11-03

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

	"github.com/hashicorp/go-azure-sdk/sdk/auth"
	"github.com/jackofallops/giovanni/storage/2023-11-03/file/files"
)

func Example() error {
	accountName := "storageaccount1"
    storageAccountKey := "ABC123...."
    shareName := "myshare"
    directoryName := "myfiles"
    fileName := "example.txt"
	domainSuffix := "core.windows.net"

	auth, err := auth.NewSharedKeyAuthorizer(accountName, storageAccountKey, auth.SharedKey)
	if err != nil {
		return fmt.Errorf("building SharedKey authorizer: %+v", err)
	}
	
    filesClient, err := files.NewWithBaseUri(fmt.Sprintf("https://%s.file.%s", accountName, domainSuffix))
	if err != nil {
		return fmt.Errorf("building client for environment: %+v", err)
	}
    filesClient.Client.SetAuthorizer(auth)
    
    ctx := context.TODO()
    input := files.CreateInput{}
    if _, err := filesClient.Create(ctx, shareName, directoryName, fileName, input); err != nil {
        return fmt.Errorf("Error creating File: %s", err)
    }
    
    return nil 
}
```