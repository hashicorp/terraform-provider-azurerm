## Table Storage Entities SDK for API version 2023-11-03

This package allows you to interact with the Entities Table Storage API

### Supported Authorizers

* SharedKeyLite (Table)

### Example Usage

```go
package main

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/auth"
	"github.com/jackofallops/giovanni/storage/2023-11-03/table/entities"
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
	
    entitiesClient, err := entities.NewWithBaseUri(fmt.Sprintf("https://%s.table.%s", accountName, domainSuffix))
	if err != nil {
		return fmt.Errorf("building client for environment: %+v", err)
	}
	entitiesClient.Client.SetAuthorizer(auth)
    
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
    if _, err := entitiesClient.Insert(ctx, tableName, input); err != nil {
        return fmt.Errorf("Error creating Entity: %s", err)
    }
    
    return nil 
}
```