
## `github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/networkconnections` Documentation

The `networkconnections` SDK allows for interaction with Azure Resource Manager `devcenter` (API Version `2025-02-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/networkconnections"
```


### Client Initialization

```go
client := networkconnections.NewNetworkConnectionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `NetworkConnectionsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := networkconnections.NewNetworkConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkConnectionName")

payload := networkconnections.NetworkConnection{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `NetworkConnectionsClient.Delete`

```go
ctx := context.TODO()
id := networkconnections.NewNetworkConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkConnectionName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `NetworkConnectionsClient.Get`

```go
ctx := context.TODO()
id := networkconnections.NewNetworkConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkConnectionName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NetworkConnectionsClient.GetHealthDetails`

```go
ctx := context.TODO()
id := networkconnections.NewNetworkConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkConnectionName")

read, err := client.GetHealthDetails(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NetworkConnectionsClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `NetworkConnectionsClient.ListBySubscription`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `NetworkConnectionsClient.ListHealthDetails`

```go
ctx := context.TODO()
id := networkconnections.NewNetworkConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkConnectionName")

// alternatively `client.ListHealthDetails(ctx, id)` can be used to do batched pagination
items, err := client.ListHealthDetailsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `NetworkConnectionsClient.Update`

```go
ctx := context.TODO()
id := networkconnections.NewNetworkConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkConnectionName")

payload := networkconnections.NetworkConnectionUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
