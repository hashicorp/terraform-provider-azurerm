
## `github.com/hashicorp/go-azure-sdk/resource-manager/storagesync/2020-03-01/serverendpointresource` Documentation

The `serverendpointresource` SDK allows for interaction with Azure Resource Manager `storagesync` (API Version `2020-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/storagesync/2020-03-01/serverendpointresource"
```


### Client Initialization

```go
client := serverendpointresource.NewServerEndpointResourceClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ServerEndpointResourceClient.ServerEndpointsCreate`

```go
ctx := context.TODO()
id := serverendpointresource.NewServerEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageSyncServiceName", "syncGroupName", "serverEndpointName")

payload := serverendpointresource.ServerEndpointCreateParameters{
	// ...
}


if err := client.ServerEndpointsCreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ServerEndpointResourceClient.ServerEndpointsDelete`

```go
ctx := context.TODO()
id := serverendpointresource.NewServerEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageSyncServiceName", "syncGroupName", "serverEndpointName")

if err := client.ServerEndpointsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ServerEndpointResourceClient.ServerEndpointsGet`

```go
ctx := context.TODO()
id := serverendpointresource.NewServerEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageSyncServiceName", "syncGroupName", "serverEndpointName")

read, err := client.ServerEndpointsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ServerEndpointResourceClient.ServerEndpointsListBySyncGroup`

```go
ctx := context.TODO()
id := serverendpointresource.NewSyncGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageSyncServiceName", "syncGroupName")

read, err := client.ServerEndpointsListBySyncGroup(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ServerEndpointResourceClient.ServerEndpointsUpdate`

```go
ctx := context.TODO()
id := serverendpointresource.NewServerEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageSyncServiceName", "syncGroupName", "serverEndpointName")

payload := serverendpointresource.ServerEndpointUpdateParameters{
	// ...
}


if err := client.ServerEndpointsUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ServerEndpointResourceClient.ServerEndpointsrecallAction`

```go
ctx := context.TODO()
id := serverendpointresource.NewServerEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageSyncServiceName", "syncGroupName", "serverEndpointName")

payload := serverendpointresource.RecallActionParameters{
	// ...
}


if err := client.ServerEndpointsrecallActionThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
