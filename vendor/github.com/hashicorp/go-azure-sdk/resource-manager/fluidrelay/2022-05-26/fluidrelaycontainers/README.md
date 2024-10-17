
## `github.com/hashicorp/go-azure-sdk/resource-manager/fluidrelay/2022-05-26/fluidrelaycontainers` Documentation

The `fluidrelaycontainers` SDK allows for interaction with Azure Resource Manager `fluidrelay` (API Version `2022-05-26`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/fluidrelay/2022-05-26/fluidrelaycontainers"
```


### Client Initialization

```go
client := fluidrelaycontainers.NewFluidRelayContainersClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `FluidRelayContainersClient.Delete`

```go
ctx := context.TODO()
id := fluidrelaycontainers.NewFluidRelayContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fluidRelayServerName", "fluidRelayContainerName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FluidRelayContainersClient.Get`

```go
ctx := context.TODO()
id := fluidrelaycontainers.NewFluidRelayContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fluidRelayServerName", "fluidRelayContainerName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FluidRelayContainersClient.ListByFluidRelayServers`

```go
ctx := context.TODO()
id := fluidrelaycontainers.NewFluidRelayServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fluidRelayServerName")

// alternatively `client.ListByFluidRelayServers(ctx, id)` can be used to do batched pagination
items, err := client.ListByFluidRelayServersComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
