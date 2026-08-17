
## `github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/dbversions` Documentation

The `dbversions` SDK allows for interaction with Azure Resource Manager `oracledatabase` (API Version `2025-09-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/dbversions"
```


### Client Initialization

```go
client := dbversions.NewDbVersionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DbVersionsClient.Get`

```go
ctx := context.TODO()
id := dbversions.NewDbSystemDbVersionID("12345678-1234-9876-4563-123456789012", "locationName", "dbSystemDbVersionName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DbVersionsClient.ListByLocation`

```go
ctx := context.TODO()
id := dbversions.NewLocationID("12345678-1234-9876-4563-123456789012", "locationName")

// alternatively `client.ListByLocation(ctx, id, dbversions.DefaultListByLocationOperationOptions())` can be used to do batched pagination
items, err := client.ListByLocationComplete(ctx, id, dbversions.DefaultListByLocationOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
