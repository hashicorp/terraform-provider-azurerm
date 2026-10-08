
## `github.com/hashicorp/go-azure-sdk/resource-manager/costmanagement/2025-03-01/scheduledactions` Documentation

The `scheduledactions` SDK allows for interaction with Azure Resource Manager `costmanagement` (API Version `2025-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/costmanagement/2025-03-01/scheduledactions"
```


### Client Initialization

```go
client := scheduledactions.NewScheduledActionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ScheduledActionsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := scheduledactions.NewScheduledActionID("scheduledActionName")

payload := scheduledactions.ScheduledAction{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload, scheduledactions.DefaultCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ScheduledActionsClient.Delete`

```go
ctx := context.TODO()
id := scheduledactions.NewScheduledActionID("scheduledActionName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ScheduledActionsClient.Get`

```go
ctx := context.TODO()
id := scheduledactions.NewScheduledActionID("scheduledActionName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ScheduledActionsClient.List`

```go
ctx := context.TODO()


// alternatively `client.List(ctx, scheduledactions.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, scheduledactions.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ScheduledActionsClient.Run`

```go
ctx := context.TODO()
id := scheduledactions.NewScheduledActionID("scheduledActionName")

read, err := client.Run(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
