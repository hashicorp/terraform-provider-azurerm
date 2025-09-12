
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-10-01/networksecurityperimeters` Documentation

The `networksecurityperimeters` SDK allows for interaction with Azure Resource Manager `network` (API Version `2024-10-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-10-01/networksecurityperimeters"
```


### Client Initialization

```go
client := networksecurityperimeters.NewNetworkSecurityPerimetersClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `NetworkSecurityPerimetersClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := networksecurityperimeters.NewNetworkSecurityPerimeterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkSecurityPerimeterName")

payload := networksecurityperimeters.NetworkSecurityPerimeter{
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


### Example Usage: `NetworkSecurityPerimetersClient.Delete`

```go
ctx := context.TODO()
id := networksecurityperimeters.NewNetworkSecurityPerimeterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkSecurityPerimeterName")

if err := client.DeleteThenPoll(ctx, id, networksecurityperimeters.DefaultDeleteOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `NetworkSecurityPerimetersClient.Get`

```go
ctx := context.TODO()
id := networksecurityperimeters.NewNetworkSecurityPerimeterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkSecurityPerimeterName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NetworkSecurityPerimetersClient.List`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.List(ctx, id, networksecurityperimeters.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, networksecurityperimeters.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `NetworkSecurityPerimetersClient.ListBySubscription`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id, networksecurityperimeters.DefaultListBySubscriptionOperationOptions())` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id, networksecurityperimeters.DefaultListBySubscriptionOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `NetworkSecurityPerimetersClient.Patch`

```go
ctx := context.TODO()
id := networksecurityperimeters.NewNetworkSecurityPerimeterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkSecurityPerimeterName")

payload := networksecurityperimeters.UpdateTagsRequest{
	// ...
}


read, err := client.Patch(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
