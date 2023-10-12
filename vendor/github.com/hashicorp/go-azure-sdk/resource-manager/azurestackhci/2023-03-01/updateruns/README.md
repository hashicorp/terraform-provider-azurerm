
## `github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2023-03-01/updateruns` Documentation

The `updateruns` SDK allows for interaction with the Azure Resource Manager Service `azurestackhci` (API Version `2023-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2023-03-01/updateruns"
```


### Client Initialization

```go
client := updateruns.NewUpdateRunsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `UpdateRunsClient.Delete`

```go
ctx := context.TODO()
id := updateruns.NewUpdateRunID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "updateValue", "updateRunValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `UpdateRunsClient.Get`

```go
ctx := context.TODO()
id := updateruns.NewUpdateRunID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "updateValue", "updateRunValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `UpdateRunsClient.List`

```go
ctx := context.TODO()
id := updateruns.NewUpdateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "updateValue")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `UpdateRunsClient.Put`

```go
ctx := context.TODO()
id := updateruns.NewUpdateRunID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "updateValue", "updateRunValue")

payload := updateruns.UpdateRun{
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
