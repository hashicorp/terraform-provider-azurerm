
## `github.com/hashicorp/go-azure-sdk/resource-manager/fluidrelay/2022-05-26/fluidrelayservers` Documentation

The `fluidrelayservers` SDK allows for interaction with the Azure Resource Manager Service `fluidrelay` (API Version `2022-05-26`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/fluidrelay/2022-05-26/fluidrelayservers"
```


### Client Initialization

```go
client := fluidrelayservers.NewFluidRelayServersClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `FluidRelayServersClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := fluidrelayservers.NewFluidRelayServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fluidRelayServerValue")

payload := fluidrelayservers.FluidRelayServer{
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


### Example Usage: `FluidRelayServersClient.Delete`

```go
ctx := context.TODO()
id := fluidrelayservers.NewFluidRelayServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fluidRelayServerValue")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FluidRelayServersClient.Get`

```go
ctx := context.TODO()
id := fluidrelayservers.NewFluidRelayServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fluidRelayServerValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FluidRelayServersClient.GetKeys`

```go
ctx := context.TODO()
id := fluidrelayservers.NewFluidRelayServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fluidRelayServerValue")

read, err := client.GetKeys(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FluidRelayServersClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := fluidrelayservers.NewResourceGroupID()

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `FluidRelayServersClient.ListBySubscription`

```go
ctx := context.TODO()
id := fluidrelayservers.NewSubscriptionID()

// alternatively `client.ListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `FluidRelayServersClient.ListKeys`

```go
ctx := context.TODO()
id := fluidrelayservers.NewFluidRelayServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fluidRelayServerValue")

read, err := client.ListKeys(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FluidRelayServersClient.RegenerateKey`

```go
ctx := context.TODO()
id := fluidrelayservers.NewFluidRelayServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fluidRelayServerValue")

payload := fluidrelayservers.RegenerateKeyRequest{
	// ...
}


read, err := client.RegenerateKey(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FluidRelayServersClient.Update`

```go
ctx := context.TODO()
id := fluidrelayservers.NewFluidRelayServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fluidRelayServerValue")

payload := fluidrelayservers.FluidRelayServerUpdate{
	// ...
}


read, err := client.Update(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
