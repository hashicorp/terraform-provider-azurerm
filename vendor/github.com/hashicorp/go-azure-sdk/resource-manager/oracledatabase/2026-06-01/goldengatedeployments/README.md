
## `github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2026-06-01/goldengatedeployments` Documentation

The `goldengatedeployments` SDK allows for interaction with Azure Resource Manager `oracledatabase` (API Version `2026-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2026-06-01/goldengatedeployments"
```


### Client Initialization

```go
client := goldengatedeployments.NewGoldenGateDeploymentsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `GoldenGateDeploymentsClient.AssignConnection`

```go
ctx := context.TODO()
id := goldengatedeployments.NewGoldenGateDeploymentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "goldenGateDeploymentName")

payload := goldengatedeployments.AssignUnassignConnection{
	// ...
}


if err := client.AssignConnectionThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `GoldenGateDeploymentsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := goldengatedeployments.NewGoldenGateDeploymentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "goldenGateDeploymentName")

payload := goldengatedeployments.GoldenGateDeployment{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `GoldenGateDeploymentsClient.Delete`

```go
ctx := context.TODO()
id := goldengatedeployments.NewGoldenGateDeploymentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "goldenGateDeploymentName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `GoldenGateDeploymentsClient.Get`

```go
ctx := context.TODO()
id := goldengatedeployments.NewGoldenGateDeploymentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "goldenGateDeploymentName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `GoldenGateDeploymentsClient.GetAssignedConnection`

```go
ctx := context.TODO()
id := goldengatedeployments.NewAssignedConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "goldenGateDeploymentName", "assignmentId")

read, err := client.GetAssignedConnection(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `GoldenGateDeploymentsClient.ListAssignedConnectionsByParent`

```go
ctx := context.TODO()
id := goldengatedeployments.NewGoldenGateDeploymentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "goldenGateDeploymentName")

// alternatively `client.ListAssignedConnectionsByParent(ctx, id)` can be used to do batched pagination
items, err := client.ListAssignedConnectionsByParentComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `GoldenGateDeploymentsClient.ListByResourceGroup`

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


### Example Usage: `GoldenGateDeploymentsClient.ListBySubscription`

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


### Example Usage: `GoldenGateDeploymentsClient.UnassignConnection`

```go
ctx := context.TODO()
id := goldengatedeployments.NewGoldenGateDeploymentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "goldenGateDeploymentName")

payload := goldengatedeployments.AssignUnassignConnection{
	// ...
}


if err := client.UnassignConnectionThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `GoldenGateDeploymentsClient.Update`

```go
ctx := context.TODO()
id := goldengatedeployments.NewGoldenGateDeploymentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "goldenGateDeploymentName")

payload := goldengatedeployments.GoldenGateDeploymentUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
