
## `github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/restorepointcollections` Documentation

The `restorepointcollections` SDK allows for interaction with Azure Resource Manager `compute` (API Version `2024-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/restorepointcollections"
```


### Client Initialization

```go
client := restorepointcollections.NewRestorePointCollectionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `RestorePointCollectionsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := restorepointcollections.NewRestorePointCollectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "restorePointCollectionName")

payload := restorepointcollections.RestorePointCollection{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RestorePointCollectionsClient.Delete`

```go
ctx := context.TODO()
id := restorepointcollections.NewRestorePointCollectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "restorePointCollectionName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `RestorePointCollectionsClient.Get`

```go
ctx := context.TODO()
id := restorepointcollections.NewRestorePointCollectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "restorePointCollectionName")

read, err := client.Get(ctx, id, restorepointcollections.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RestorePointCollectionsClient.List`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `RestorePointCollectionsClient.ListAll`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListAll(ctx, id)` can be used to do batched pagination
items, err := client.ListAllComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `RestorePointCollectionsClient.Update`

```go
ctx := context.TODO()
id := restorepointcollections.NewRestorePointCollectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "restorePointCollectionName")

payload := restorepointcollections.RestorePointCollectionUpdate{
	// ...
}


read, err := client.Update(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
