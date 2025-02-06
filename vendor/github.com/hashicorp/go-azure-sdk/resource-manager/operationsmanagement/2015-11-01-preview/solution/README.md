
## `github.com/hashicorp/go-azure-sdk/resource-manager/operationsmanagement/2015-11-01-preview/solution` Documentation

The `solution` SDK allows for interaction with Azure Resource Manager `operationsmanagement` (API Version `2015-11-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/operationsmanagement/2015-11-01-preview/solution"
```


### Client Initialization

```go
client := solution.NewSolutionClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SolutionClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := solution.NewSolutionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "solutionName")

payload := solution.Solution{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `SolutionClient.Delete`

```go
ctx := context.TODO()
id := solution.NewSolutionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "solutionName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `SolutionClient.Get`

```go
ctx := context.TODO()
id := solution.NewSolutionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "solutionName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SolutionClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

read, err := client.ListByResourceGroup(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SolutionClient.ListBySubscription`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

read, err := client.ListBySubscription(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SolutionClient.Update`

```go
ctx := context.TODO()
id := solution.NewSolutionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "solutionName")

payload := solution.SolutionPatch{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
