
## `github.com/hashicorp/go-azure-sdk/data-plane/synapse/2021-06-01-preview/managedprivateendpoints` Documentation

The `managedprivateendpoints` SDK allows for interaction with <unknown source data type> `synapse` (API Version `2021-06-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/data-plane/synapse/2021-06-01-preview/managedprivateendpoints"
```


### Client Initialization

```go
client := managedprivateendpoints.NewManagedPrivateEndpointsClientWithBaseURI("")
client.Client.Authorizer = authorizer
```


### Example Usage: `ManagedPrivateEndpointsClient.Create`

```go
ctx := context.TODO()
id := managedprivateendpoints.NewManagedPrivateEndpointID("managedVirtualNetworkName", "managedPrivateEndpointName")

payload := managedprivateendpoints.ManagedPrivateEndpoint{
	// ...
}


read, err := client.Create(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagedPrivateEndpointsClient.Delete`

```go
ctx := context.TODO()
id := managedprivateendpoints.NewManagedPrivateEndpointID("managedVirtualNetworkName", "managedPrivateEndpointName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagedPrivateEndpointsClient.Get`

```go
ctx := context.TODO()
id := managedprivateendpoints.NewManagedPrivateEndpointID("managedVirtualNetworkName", "managedPrivateEndpointName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagedPrivateEndpointsClient.List`

```go
ctx := context.TODO()
id := managedprivateendpoints.NewManagedVirtualNetworkID("managedVirtualNetworkName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
