## Table Storage Entities SDK for API version 2019-12-12

This package allows you to interact with the Entities Table Storage API

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
	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/table/entities"
)

func Example() error {
	accountName := "storageaccount1"
    storageAccountKey := "ABC123...."
    tableName := "mytable"
    
    storageAuth := autorest.NewSharedKeyLiteTableAuthorizer(accountName, storageAccountKey)
    entitiesClient := entities.New()
    entitiesClient.Client.Authorizer = storageAuth
    
    ctx := context.TODO()
    input := entities.InsertEntityInput{
    	PartitionKey: "abc",
    	RowKey: "123",
    	MetaDataLevel: entities.NoMetaData,
    	Entity: map[string]interface{}{
    	    "title": "Don't Kill My Vibe",
    	    "artist": "Sigrid",
    	},
    }
    if _, err := entitiesClient.Insert(ctx, accountName, tableName, input); err != nil {
        return fmt.Errorf("Error creating Entity: %s", err)
    }
    
    return nil 
}
```