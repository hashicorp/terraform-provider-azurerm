
## `github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/edgedevices` Documentation

The `edgedevices` SDK allows for interaction with Azure Resource Manager `azurestackhci` (API Version `2024-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/edgedevices"
```


### Client Initialization

```go
client := edgedevices.NewEdgeDevicesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `EdgeDevicesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := edgedevices.NewScopedEdgeDeviceID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "edgeDeviceName")

payload := edgedevices.EdgeDevice{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `EdgeDevicesClient.Delete`

```go
ctx := context.TODO()
id := edgedevices.NewScopedEdgeDeviceID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "edgeDeviceName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `EdgeDevicesClient.Get`

```go
ctx := context.TODO()
id := edgedevices.NewScopedEdgeDeviceID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "edgeDeviceName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EdgeDevicesClient.List`

```go
ctx := context.TODO()
id := commonids.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `EdgeDevicesClient.Validate`

```go
ctx := context.TODO()
id := edgedevices.NewScopedEdgeDeviceID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "edgeDeviceName")

payload := edgedevices.ValidateRequest{
	// ...
}


if err := client.ValidateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
