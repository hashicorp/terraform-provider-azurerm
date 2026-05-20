## File Storage Directories SDK for API version 2023-11-03

This package allows you to interact with the Directories File Storage API

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
	"github.com/jackofallops/giovanni/storage/2023-11-03/file/directories"
)

func Example() error {
	accountName := "storageaccount1"
    storageAccountKey := "ABC123...."
    shareName := "myshare"
    directoryName := "myfiles"
	domainSuffix := "core.windows.net"

	auth, err := auth.NewSharedKeyAuthorizer(accountName, storageAccountKey, auth.SharedKey)
	if err != nil {
		return fmt.Errorf("building SharedKey authorizer: %+v", err)
	}
	
    directoriesClient, err := directories.NewWithBaseUri(fmt.Sprintf("https://%s.dfs.%s", accountName, domainSuffix))
	if err != nil {
		return fmt.Errorf("building client for environment: %+v", err)
	}
    directoriesClient.Client.SetAuthorizer(auth)
    
    ctx := context.TODO()
    metadata := map[string]string{
    	"hello": "world",
    }
	
	input := directories.CreateDirectoryInput{
		MetaData: metadata,
    }
    if _, err := directoriesClient.Create(ctx, shareName, directoryName, input); err != nil {
        return fmt.Errorf("Error creating Directory: %s", err)
    }
    
    return nil 
}
```