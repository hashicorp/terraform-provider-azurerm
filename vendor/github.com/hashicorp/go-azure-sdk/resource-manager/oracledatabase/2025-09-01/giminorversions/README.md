
## `github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/giminorversions` Documentation

The `giminorversions` SDK allows for interaction with Azure Resource Manager `oracledatabase` (API Version `2025-09-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/giminorversions"
```


### Client Initialization

```go
client := giminorversions.NewGiMinorVersionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `GiMinorVersionsClient.Get`

```go
ctx := context.TODO()
id := giminorversions.NewGiMinorVersionID("12345678-1234-9876-4563-123456789012", "locationName", "giVersionName", "giMinorVersionName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `GiMinorVersionsClient.ListByParent`

```go
ctx := context.TODO()
id := giminorversions.NewGiVersionID("12345678-1234-9876-4563-123456789012", "locationName", "giVersionName")

// alternatively `client.ListByParent(ctx, id, giminorversions.DefaultListByParentOperationOptions())` can be used to do batched pagination
items, err := client.ListByParentComplete(ctx, id, giminorversions.DefaultListByParentOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
