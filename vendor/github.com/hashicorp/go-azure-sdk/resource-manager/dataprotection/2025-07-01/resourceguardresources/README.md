
## `github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/resourceguardresources` Documentation

The `resourceguardresources` SDK allows for interaction with Azure Resource Manager `dataprotection` (API Version `2025-07-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/resourceguardresources"
```


### Client Initialization

```go
client := resourceguardresources.NewResourceGuardResourcesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ResourceGuardResourcesClient.ResourceGuardsDelete`

```go
ctx := context.TODO()
id := resourceguardresources.NewResourceGuardID("12345678-1234-9876-4563-123456789012", "example-resource-group", "resourceGuardName")

read, err := client.ResourceGuardsDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ResourceGuardResourcesClient.ResourceGuardsGet`

```go
ctx := context.TODO()
id := resourceguardresources.NewResourceGuardID("12345678-1234-9876-4563-123456789012", "example-resource-group", "resourceGuardName")

read, err := client.ResourceGuardsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ResourceGuardResourcesClient.ResourceGuardsGetResourcesInResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ResourceGuardsGetResourcesInResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ResourceGuardsGetResourcesInResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ResourceGuardResourcesClient.ResourceGuardsGetResourcesInSubscription`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ResourceGuardsGetResourcesInSubscription(ctx, id)` can be used to do batched pagination
items, err := client.ResourceGuardsGetResourcesInSubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ResourceGuardResourcesClient.ResourceGuardsPatch`

```go
ctx := context.TODO()
id := resourceguardresources.NewResourceGuardID("12345678-1234-9876-4563-123456789012", "example-resource-group", "resourceGuardName")

payload := resourceguardresources.PatchResourceGuardInput{
	// ...
}


read, err := client.ResourceGuardsPatch(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ResourceGuardResourcesClient.ResourceGuardsPut`

```go
ctx := context.TODO()
id := resourceguardresources.NewResourceGuardID("12345678-1234-9876-4563-123456789012", "example-resource-group", "resourceGuardName")

payload := resourceguardresources.ResourceGuardResource{
	// ...
}


read, err := client.ResourceGuardsPut(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
