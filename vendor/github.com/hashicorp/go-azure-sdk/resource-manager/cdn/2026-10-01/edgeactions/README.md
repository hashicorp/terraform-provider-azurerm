
## `github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/edgeactions` Documentation

The `edgeactions` SDK allows for interaction with Azure Resource Manager `cdn` (API Version `2026-10-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/edgeactions"
```


### Client Initialization

```go
client := edgeactions.NewEdgeActionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `EdgeActionsClient.Create`

```go
ctx := context.TODO()
id := edgeactions.NewEdgeActionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "edgeActionName")

payload := edgeactions.EdgeAction{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `EdgeActionsClient.Delete`

```go
ctx := context.TODO()
id := edgeactions.NewEdgeActionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "edgeActionName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `EdgeActionsClient.Get`

```go
ctx := context.TODO()
id := edgeactions.NewEdgeActionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "edgeActionName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EdgeActionsClient.ListByResourceGroup`

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


### Example Usage: `EdgeActionsClient.ListBySubscription`

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


### Example Usage: `EdgeActionsClient.Update`

```go
ctx := context.TODO()
id := edgeactions.NewEdgeActionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "edgeActionName")

payload := edgeactions.EdgeActionUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
