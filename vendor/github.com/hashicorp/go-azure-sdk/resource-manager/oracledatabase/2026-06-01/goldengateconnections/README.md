
## `github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2026-06-01/goldengateconnections` Documentation

The `goldengateconnections` SDK allows for interaction with Azure Resource Manager `oracledatabase` (API Version `2026-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2026-06-01/goldengateconnections"
```


### Client Initialization

```go
client := goldengateconnections.NewGoldenGateConnectionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `GoldenGateConnectionsClient.AssignDeployment`

```go
ctx := context.TODO()
id := goldengateconnections.NewGoldenGateConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "goldenGateConnectionName")

payload := goldengateconnections.AssignUnassignDeployment{
	// ...
}


if err := client.AssignDeploymentThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `GoldenGateConnectionsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := goldengateconnections.NewGoldenGateConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "goldenGateConnectionName")

payload := goldengateconnections.GoldenGateConnection{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `GoldenGateConnectionsClient.Delete`

```go
ctx := context.TODO()
id := goldengateconnections.NewGoldenGateConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "goldenGateConnectionName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `GoldenGateConnectionsClient.Get`

```go
ctx := context.TODO()
id := goldengateconnections.NewGoldenGateConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "goldenGateConnectionName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `GoldenGateConnectionsClient.GetAssignedDeployment`

```go
ctx := context.TODO()
id := goldengateconnections.NewAssignedDeploymentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "goldenGateConnectionName", "assignmentId")

read, err := client.GetAssignedDeployment(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `GoldenGateConnectionsClient.ListAssignedDeploymentsByParent`

```go
ctx := context.TODO()
id := goldengateconnections.NewGoldenGateConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "goldenGateConnectionName")

// alternatively `client.ListAssignedDeploymentsByParent(ctx, id)` can be used to do batched pagination
items, err := client.ListAssignedDeploymentsByParentComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `GoldenGateConnectionsClient.ListByResourceGroup`

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


### Example Usage: `GoldenGateConnectionsClient.ListBySubscription`

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


### Example Usage: `GoldenGateConnectionsClient.UnassignDeployment`

```go
ctx := context.TODO()
id := goldengateconnections.NewGoldenGateConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "goldenGateConnectionName")

payload := goldengateconnections.AssignUnassignDeployment{
	// ...
}


if err := client.UnassignDeploymentThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `GoldenGateConnectionsClient.Update`

```go
ctx := context.TODO()
id := goldengateconnections.NewGoldenGateConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "goldenGateConnectionName")

payload := goldengateconnections.GoldenGateConnectionUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
