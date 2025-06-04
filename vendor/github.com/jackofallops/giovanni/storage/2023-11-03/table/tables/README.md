## Table Storage Tables SDK for API version 2023-11-03

This package allows you to interact with the Tables Table Storage API

### Supported Authorizers

* SharedKeyLite (Table)

### Example Usage

```go
package main

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/auth"
	"github.com/jackofallops/giovanni/storage/2023-11-03/table/tables"
)

func Example() error {
	accountName := "storageaccount1"
    storageAccountKey := "ABC123...."
    tableName := "mytable"
	domainSuffix := "core.windows.net"

	auth, err := auth.NewSharedKeyAuthorizer(accountName, storageAccountKey, auth.SharedKeyTable)
	if err != nil {
		return fmt.Errorf("building SharedKey authorizer: %+v", err)
	}
    tablesClient, err := tables.NewWithBaseUri(fmt.Sprintf("https://%s.table.%s", accountName, domainSuffix))
	if err != nil {
		return fmt.Errorf("building client for environment: %+v", err)
	}
	tablesClient.Client.SetAuthorizer(auth)
    
    ctx := context.TODO()
    if _, err := tablesClient.Create(ctx, tableName); err != nil {
        return fmt.Errorf("Error creating Table: %s", err)
    }
    
    return nil 
}
```