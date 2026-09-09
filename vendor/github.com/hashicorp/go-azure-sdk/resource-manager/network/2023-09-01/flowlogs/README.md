
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/flowlogs` Documentation

The `flowlogs` SDK allows for interaction with Azure Resource Manager `network` (API Version `2023-09-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/flowlogs"
```


### Client Initialization

```go
client := flowlogs.NewFlowLogsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `FlowLogsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := flowlogs.NewFlowLogID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkWatcherName", "flowLogName")

payload := flowlogs.FlowLog{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `FlowLogsClient.Delete`

```go
ctx := context.TODO()
id := flowlogs.NewFlowLogID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkWatcherName", "flowLogName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `FlowLogsClient.Get`

```go
ctx := context.TODO()
id := flowlogs.NewFlowLogID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkWatcherName", "flowLogName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FlowLogsClient.List`

```go
ctx := context.TODO()
id := flowlogs.NewNetworkWatcherID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkWatcherName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `FlowLogsClient.UpdateTags`

```go
ctx := context.TODO()
id := flowlogs.NewFlowLogID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkWatcherName", "flowLogName")

payload := flowlogs.TagsObject{
	// ...
}


read, err := client.UpdateTags(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
