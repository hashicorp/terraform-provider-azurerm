
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


### Example Usage: `UpdatesClient.Delete`

```go
ctx := context.TODO()
id := updates.NewUpdateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "updateValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `UpdatesClient.Get`

```go
ctx := context.TODO()
id := updates.NewUpdateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "updateValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `UpdatesClient.List`

```go
ctx := context.TODO()
id := updates.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `UpdatesClient.Post`

```go
ctx := context.TODO()
id := updates.NewUpdateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "updateValue")

if err := client.PostThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `UpdatesClient.Put`

```go
ctx := context.TODO()
id := updates.NewUpdateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "updateValue")

payload := updates.Update{
	// ...
}


read, err := client.Put(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
