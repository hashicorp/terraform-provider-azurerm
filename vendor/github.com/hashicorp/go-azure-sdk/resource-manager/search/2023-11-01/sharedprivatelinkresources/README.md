
## `github.com/hashicorp/go-azure-sdk/resource-manager/search/2023-11-01/sharedprivatelinkresources` Documentation

The `sharedprivatelinkresources` SDK allows for interaction with Azure Resource Manager `search` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/search/2023-11-01/sharedprivatelinkresources"
```


### Client Initialization

```go
client := sharedprivatelinkresources.NewSharedPrivateLinkResourcesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SharedPrivateLinkResourcesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := sharedprivatelinkresources.NewSharedPrivateLinkResourceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "searchServiceName", "sharedPrivateLinkResourceName")

payload := sharedprivatelinkresources.SharedPrivateLinkResource{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload, sharedprivatelinkresources.DefaultCreateOrUpdateOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `SharedPrivateLinkResourcesClient.Delete`

```go
ctx := context.TODO()
id := sharedprivatelinkresources.NewSharedPrivateLinkResourceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "searchServiceName", "sharedPrivateLinkResourceName")

if err := client.DeleteThenPoll(ctx, id, sharedprivatelinkresources.DefaultDeleteOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `SharedPrivateLinkResourcesClient.Get`

```go
ctx := context.TODO()
id := sharedprivatelinkresources.NewSharedPrivateLinkResourceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "searchServiceName", "sharedPrivateLinkResourceName")

read, err := client.Get(ctx, id, sharedprivatelinkresources.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SharedPrivateLinkResourcesClient.ListByService`

```go
ctx := context.TODO()
id := sharedprivatelinkresources.NewSearchServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "searchServiceName")

// alternatively `client.ListByService(ctx, id, sharedprivatelinkresources.DefaultListByServiceOperationOptions())` can be used to do batched pagination
items, err := client.ListByServiceComplete(ctx, id, sharedprivatelinkresources.DefaultListByServiceOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
