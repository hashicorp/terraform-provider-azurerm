
## `github.com/hashicorp/go-azure-sdk/resource-manager/storagesync/2020-03-01/registeredserverresource` Documentation

The `registeredserverresource` SDK allows for interaction with Azure Resource Manager `storagesync` (API Version `2020-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/storagesync/2020-03-01/registeredserverresource"
```


### Client Initialization

```go
client := registeredserverresource.NewRegisteredServerResourceClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `RegisteredServerResourceClient.RegisteredServersCreate`

```go
ctx := context.TODO()
id := registeredserverresource.NewRegisteredServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageSyncServiceName", "serverId")

payload := registeredserverresource.RegisteredServerCreateParameters{
	// ...
}


if err := client.RegisteredServersCreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `RegisteredServerResourceClient.RegisteredServersDelete`

```go
ctx := context.TODO()
id := registeredserverresource.NewRegisteredServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageSyncServiceName", "serverId")

if err := client.RegisteredServersDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `RegisteredServerResourceClient.RegisteredServersGet`

```go
ctx := context.TODO()
id := registeredserverresource.NewRegisteredServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageSyncServiceName", "serverId")

read, err := client.RegisteredServersGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RegisteredServerResourceClient.RegisteredServersListByStorageSyncService`

```go
ctx := context.TODO()
id := registeredserverresource.NewStorageSyncServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageSyncServiceName")

read, err := client.RegisteredServersListByStorageSyncService(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RegisteredServerResourceClient.RegisteredServerstriggerRollover`

```go
ctx := context.TODO()
id := registeredserverresource.NewRegisteredServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageSyncServiceName", "serverId")

payload := registeredserverresource.TriggerRolloverRequest{
	// ...
}


if err := client.RegisteredServerstriggerRolloverThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
