
## `github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/resourceguards` Documentation

The `resourceguards` SDK allows for interaction with Azure Resource Manager `dataprotection` (API Version `2025-07-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/resourceguards"
```


### Client Initialization

```go
client := resourceguards.NewResourceGuardsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ResourceGuardsClient.GetDefaultDeleteResourceGuardProxyRequestsObject`

```go
ctx := context.TODO()
id := resourceguards.NewDeleteResourceGuardProxyRequestID("12345678-1234-9876-4563-123456789012", "example-resource-group", "resourceGuardName", "deleteResourceGuardProxyRequestName")

read, err := client.GetDefaultDeleteResourceGuardProxyRequestsObject(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ResourceGuardsClient.GetDeleteResourceGuardProxyRequestsObjects`

```go
ctx := context.TODO()
id := resourceguards.NewResourceGuardID("12345678-1234-9876-4563-123456789012", "example-resource-group", "resourceGuardName")

// alternatively `client.GetDeleteResourceGuardProxyRequestsObjects(ctx, id)` can be used to do batched pagination
items, err := client.GetDeleteResourceGuardProxyRequestsObjectsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
