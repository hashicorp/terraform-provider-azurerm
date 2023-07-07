
## `github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2023-05-02/managedprivateendpoints` Documentation

The `managedprivateendpoints` SDK allows for interaction with the Azure Resource Manager Service `kusto` (API Version `2023-05-02`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2023-05-02/managedprivateendpoints"
```


### Client Initialization

```go
client := managedprivateendpoints.NewManagedPrivateEndpointsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ManagedPrivateEndpointsClient.CheckNameAvailability`

```go
ctx := context.TODO()
id := managedprivateendpoints.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue")

payload := managedprivateendpoints.ManagedPrivateEndpointsCheckNameRequest{
	// ...
}


read, err := client.CheckNameAvailability(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagedPrivateEndpointsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := managedprivateendpoints.NewManagedPrivateEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "managedPrivateEndpointValue")

payload := managedprivateendpoints.ManagedPrivateEndpoint{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ManagedPrivateEndpointsClient.Delete`

```go
ctx := context.TODO()
id := managedprivateendpoints.NewManagedPrivateEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "managedPrivateEndpointValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ManagedPrivateEndpointsClient.Get`

```go
ctx := context.TODO()
id := managedprivateendpoints.NewManagedPrivateEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "managedPrivateEndpointValue")

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
id := managedprivateendpoints.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue")

read, err := client.List(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagedPrivateEndpointsClient.Update`

```go
ctx := context.TODO()
id := managedprivateendpoints.NewManagedPrivateEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "managedPrivateEndpointValue")

payload := managedprivateendpoints.ManagedPrivateEndpoint{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
