
## `github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2022-06-15/privatelinkresources` Documentation

The `privatelinkresources` SDK allows for interaction with Azure Resource Manager `eventgrid` (API Version `2022-06-15`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2022-06-15/privatelinkresources"
```


### Client Initialization

```go
client := privatelinkresources.NewPrivateLinkResourcesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `PrivateLinkResourcesClient.Get`

```go
ctx := context.TODO()
id := privatelinkresources.NewScopedPrivateLinkResourceID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "privateLinkResourceName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PrivateLinkResourcesClient.ListByResource`

```go
ctx := context.TODO()
id := commonids.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

// alternatively `client.ListByResource(ctx, id, privatelinkresources.DefaultListByResourceOperationOptions())` can be used to do batched pagination
items, err := client.ListByResourceComplete(ctx, id, privatelinkresources.DefaultListByResourceOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
