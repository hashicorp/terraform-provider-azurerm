
## `github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/storageinsights` Documentation

The `storageinsights` SDK allows for interaction with the Azure Resource Manager Service `operationalinsights` (API Version `2020-08-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/storageinsights"
```


### Client Initialization

```go
client := storageinsights.NewStorageInsightsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `StorageInsightsClient.StorageInsightConfigsCreateOrUpdate`

```go
ctx := context.TODO()
id := storageinsights.NewStorageInsightConfigID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "storageInsightConfigValue")

payload := storageinsights.StorageInsight{
	// ...
}


read, err := client.StorageInsightConfigsCreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StorageInsightsClient.StorageInsightConfigsDelete`

```go
ctx := context.TODO()
id := storageinsights.NewStorageInsightConfigID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "storageInsightConfigValue")

read, err := client.StorageInsightConfigsDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StorageInsightsClient.StorageInsightConfigsGet`

```go
ctx := context.TODO()
id := storageinsights.NewStorageInsightConfigID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "storageInsightConfigValue")

read, err := client.StorageInsightConfigsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StorageInsightsClient.StorageInsightConfigsListByWorkspace`

```go
ctx := context.TODO()
id := storageinsights.NewWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue")

// alternatively `client.StorageInsightConfigsListByWorkspace(ctx, id)` can be used to do batched pagination
items, err := client.StorageInsightConfigsListByWorkspaceComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
