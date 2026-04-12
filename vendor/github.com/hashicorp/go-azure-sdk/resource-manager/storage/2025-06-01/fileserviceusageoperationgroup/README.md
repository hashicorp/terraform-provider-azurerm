
## `github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-06-01/fileserviceusageoperationgroup` Documentation

The `fileserviceusageoperationgroup` SDK allows for interaction with Azure Resource Manager `storage` (API Version `2025-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-06-01/fileserviceusageoperationgroup"
```


### Client Initialization

```go
client := fileserviceusageoperationgroup.NewFileServiceUsageOperationGroupClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `FileServiceUsageOperationGroupClient.FileServicesGetServiceUsage`

```go
ctx := context.TODO()
id := commonids.NewStorageAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountName")

read, err := client.FileServicesGetServiceUsage(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FileServiceUsageOperationGroupClient.FileServicesListServiceUsages`

```go
ctx := context.TODO()
id := commonids.NewStorageAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountName")

// alternatively `client.FileServicesListServiceUsages(ctx, id, fileserviceusageoperationgroup.DefaultFileServicesListServiceUsagesOperationOptions())` can be used to do batched pagination
items, err := client.FileServicesListServiceUsagesComplete(ctx, id, fileserviceusageoperationgroup.DefaultFileServicesListServiceUsagesOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
