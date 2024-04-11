
## `github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/networkinterfaces` Documentation

The `networkinterfaces` SDK allows for interaction with the Azure Resource Manager Service `azurestackhci` (API Version `2024-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/networkinterfaces"
```


### Client Initialization

```go
client := networkinterfaces.NewNetworkInterfacesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `NetworkInterfacesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := networkinterfaces.NewNetworkInterfaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkInterfaceValue")

payload := networkinterfaces.NetworkInterfaces{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `NetworkInterfacesClient.Delete`

```go
ctx := context.TODO()
id := networkinterfaces.NewNetworkInterfaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkInterfaceValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `NetworkInterfacesClient.Get`

```go
ctx := context.TODO()
id := networkinterfaces.NewNetworkInterfaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkInterfaceValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NetworkInterfacesClient.List`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `NetworkInterfacesClient.ListAll`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListAll(ctx, id)` can be used to do batched pagination
items, err := client.ListAllComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `NetworkInterfacesClient.Update`

```go
ctx := context.TODO()
id := networkinterfaces.NewNetworkInterfaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkInterfaceValue")

payload := networkinterfaces.NetworkInterfacesUpdateRequest{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
