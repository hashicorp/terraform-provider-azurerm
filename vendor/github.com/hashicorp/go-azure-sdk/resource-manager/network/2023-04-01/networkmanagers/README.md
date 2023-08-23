
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/networkmanagers` Documentation

The `networkmanagers` SDK allows for interaction with the Azure Resource Manager Service `network` (API Version `2023-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/networkmanagers"
```


### Client Initialization

```go
client := networkmanagers.NewNetworkManagersClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `NetworkManagersClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := networkmanagers.NewNetworkManagerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerValue")

payload := networkmanagers.NetworkManager{
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


### Example Usage: `NetworkManagersClient.Delete`

```go
ctx := context.TODO()
id := networkmanagers.NewNetworkManagerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerValue")

if err := client.DeleteThenPoll(ctx, id, networkmanagers.DefaultDeleteOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `NetworkManagersClient.Get`

```go
ctx := context.TODO()
id := networkmanagers.NewNetworkManagerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NetworkManagersClient.List`

```go
ctx := context.TODO()
id := networkmanagers.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.List(ctx, id, networkmanagers.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, networkmanagers.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `NetworkManagersClient.ListBySubscription`

```go
ctx := context.TODO()
id := networkmanagers.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id, networkmanagers.DefaultListBySubscriptionOperationOptions())` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id, networkmanagers.DefaultListBySubscriptionOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `NetworkManagersClient.NetworkManagerCommitsPost`

```go
ctx := context.TODO()
id := networkmanagers.NewNetworkManagerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerValue")

payload := networkmanagers.NetworkManagerCommit{
	// ...
}


if err := client.NetworkManagerCommitsPostThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `NetworkManagersClient.NetworkManagerDeploymentStatusList`

```go
ctx := context.TODO()
id := networkmanagers.NewNetworkManagerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerValue")

payload := networkmanagers.NetworkManagerDeploymentStatusParameter{
	// ...
}


read, err := client.NetworkManagerDeploymentStatusList(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NetworkManagersClient.Patch`

```go
ctx := context.TODO()
id := networkmanagers.NewNetworkManagerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerValue")

payload := networkmanagers.PatchObject{
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
