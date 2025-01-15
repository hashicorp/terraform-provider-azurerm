
## `github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2023-01-01/createresource` Documentation

The `createresource` SDK allows for interaction with Azure Resource Manager `datadog` (API Version `2023-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2023-01-01/createresource"
```


### Client Initialization

```go
client := createresource.NewCreateResourceClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `CreateResourceClient.CreationSupportedGet`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

read, err := client.CreationSupportedGet(ctx, id, createresource.DefaultCreationSupportedGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CreateResourceClient.CreationSupportedList`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.CreationSupportedList(ctx, id, createresource.DefaultCreationSupportedListOperationOptions())` can be used to do batched pagination
items, err := client.CreationSupportedListComplete(ctx, id, createresource.DefaultCreationSupportedListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
