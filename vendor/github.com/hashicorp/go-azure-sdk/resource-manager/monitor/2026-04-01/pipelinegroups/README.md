
## `github.com/hashicorp/go-azure-sdk/resource-manager/monitor/2026-04-01/pipelinegroups` Documentation

The `pipelinegroups` SDK allows for interaction with Azure Resource Manager `monitor` (API Version `2026-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/monitor/2026-04-01/pipelinegroups"
```


### Client Initialization

```go
client := pipelinegroups.NewPipelineGroupsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `PipelineGroupsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := pipelinegroups.NewPipelineGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "pipelineGroupName")

payload := pipelinegroups.PipelineGroup{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `PipelineGroupsClient.Delete`

```go
ctx := context.TODO()
id := pipelinegroups.NewPipelineGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "pipelineGroupName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `PipelineGroupsClient.Get`

```go
ctx := context.TODO()
id := pipelinegroups.NewPipelineGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "pipelineGroupName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PipelineGroupsClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `PipelineGroupsClient.ListBySubscription`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `PipelineGroupsClient.Update`

```go
ctx := context.TODO()
id := pipelinegroups.NewPipelineGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "pipelineGroupName")

payload := pipelinegroups.PipelineGroupUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
