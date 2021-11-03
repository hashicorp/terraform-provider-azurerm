## Table Storage Tables SDK for API version 2019-12-12

This package allows you to interact with the Tables Table Storage API

### Supported Authorizers

* SharedKeyLite (Table)

### Example Usage

```go
package main

import (
	"context"
	"fmt"
	"time"
	
	"github.com/Azure/go-autorest/autorest"
	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/table/tables"
)

func Example() error {
	accountName := "storageaccount1"
    storageAccountKey := "ABC123...."
    tableName := "mytable"
    
    storageAuth := autorest.NewSharedKeyLiteTableAuthorizer(accountName, storageAccountKey)
    tablesClient := tables.New()
    tablesClient.Client.Authorizer = storageAuth
    
    ctx := context.TODO()
    if _, err := tablesClient.Insert(ctx, accountName, tableName); err != nil {
        return fmt.Errorf("Error creating Table: %s", err)
    }
    
    return nil 
}
```