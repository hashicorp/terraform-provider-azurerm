
## `github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2023-03-01/updates` Documentation

The `updates` SDK allows for interaction with the Azure Resource Manager Service `azurestackhci` (API Version `2023-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2023-03-01/updates"
```


### Client Initialization

```go
client := updates.NewUpdatesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `UpdatesClient.UpdatesDelete`

```go
ctx := context.TODO()
id := updates.NewUpdateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "updateValue")

if err := client.UpdatesDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `UpdatesClient.UpdatesGet`

```go
ctx := context.TODO()
id := updates.NewUpdateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "updateValue")

read, err := client.UpdatesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `UpdatesClient.UpdatesList`

```go
ctx := context.TODO()
id := updates.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue")

// alternatively `client.UpdatesList(ctx, id)` can be used to do batched pagination
items, err := client.UpdatesListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `UpdatesClient.UpdatesPost`

```go
ctx := context.TODO()
id := updates.NewUpdateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "updateValue")

if err := client.UpdatesPostThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `UpdatesClient.UpdatesPut`

```go
ctx := context.TODO()
id := updates.NewUpdateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "updateValue")

payload := updates.Update{
	// ...
}


read, err := client.UpdatesPut(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
