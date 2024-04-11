
## `github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-02-01-preview/transparentdataencryptions` Documentation

The `transparentdataencryptions` SDK allows for interaction with the Azure Resource Manager Service `sql` (API Version `2023-02-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-02-01-preview/transparentdataencryptions"
```


### Client Initialization

```go
client := transparentdataencryptions.NewTransparentDataEncryptionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `TransparentDataEncryptionsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := commonids.NewSqlDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue", "databaseValue")

payload := transparentdataencryptions.LogicalDatabaseTransparentDataEncryption{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `TransparentDataEncryptionsClient.Get`

```go
ctx := context.TODO()
id := commonids.NewSqlDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue", "databaseValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TransparentDataEncryptionsClient.ListByDatabase`

```go
ctx := context.TODO()
id := commonids.NewSqlDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue", "databaseValue")

// alternatively `client.ListByDatabase(ctx, id)` can be used to do batched pagination
items, err := client.ListByDatabaseComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
