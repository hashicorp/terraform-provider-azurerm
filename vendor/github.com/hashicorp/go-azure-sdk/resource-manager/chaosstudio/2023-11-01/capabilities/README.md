
## `github.com/hashicorp/go-azure-sdk/resource-manager/chaosstudio/2023-11-01/capabilities` Documentation

The `capabilities` SDK allows for interaction with Azure Resource Manager `chaosstudio` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/chaosstudio/2023-11-01/capabilities"
```


### Client Initialization

```go
client := capabilities.NewCapabilitiesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `CapabilitiesClient.CapabilityTypesGet`

```go
ctx := context.TODO()
id := capabilities.NewCapabilityTypeID("12345678-1234-9876-4563-123456789012", "locationName", "targetTypeName", "capabilityTypeName")

read, err := client.CapabilityTypesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CapabilitiesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := commonids.NewChaosStudioCapabilityID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "targetName", "capabilityName")

payload := capabilities.Capability{
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


### Example Usage: `CapabilitiesClient.Delete`

```go
ctx := context.TODO()
id := commonids.NewChaosStudioCapabilityID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "targetName", "capabilityName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CapabilitiesClient.Get`

```go
ctx := context.TODO()
id := commonids.NewChaosStudioCapabilityID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "targetName", "capabilityName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CapabilitiesClient.List`

```go
ctx := context.TODO()
id := commonids.NewChaosStudioTargetID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "targetName")

// alternatively `client.List(ctx, id, capabilities.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, capabilities.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
