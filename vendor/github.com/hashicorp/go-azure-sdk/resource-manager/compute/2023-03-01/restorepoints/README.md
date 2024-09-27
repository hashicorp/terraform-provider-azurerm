
## `github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-03-01/restorepoints` Documentation

The `restorepoints` SDK allows for interaction with Azure Resource Manager `compute` (API Version `2023-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-03-01/restorepoints"
```


### Client Initialization

```go
client := restorepoints.NewRestorePointsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `RestorePointsClient.Create`

```go
ctx := context.TODO()
id := restorepoints.NewRestorePointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "restorePointCollectionName", "restorePointName")

payload := restorepoints.RestorePoint{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `RestorePointsClient.Delete`

```go
ctx := context.TODO()
id := restorepoints.NewRestorePointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "restorePointCollectionName", "restorePointName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `RestorePointsClient.Get`

```go
ctx := context.TODO()
id := restorepoints.NewRestorePointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "restorePointCollectionName", "restorePointName")

read, err := client.Get(ctx, id, restorepoints.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
