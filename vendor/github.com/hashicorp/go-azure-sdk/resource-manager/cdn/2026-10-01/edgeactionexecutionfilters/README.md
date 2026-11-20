
## `github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/edgeactionexecutionfilters` Documentation

The `edgeactionexecutionfilters` SDK allows for interaction with Azure Resource Manager `cdn` (API Version `2026-10-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/edgeactionexecutionfilters"
```


### Client Initialization

```go
client := edgeactionexecutionfilters.NewEdgeActionExecutionFiltersClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `EdgeActionExecutionFiltersClient.Create`

```go
ctx := context.TODO()
id := edgeactionexecutionfilters.NewExecutionFilterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "edgeActionName", "executionFilterName")

payload := edgeactionexecutionfilters.EdgeActionExecutionFilter{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `EdgeActionExecutionFiltersClient.Delete`

```go
ctx := context.TODO()
id := edgeactionexecutionfilters.NewExecutionFilterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "edgeActionName", "executionFilterName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `EdgeActionExecutionFiltersClient.Get`

```go
ctx := context.TODO()
id := edgeactionexecutionfilters.NewExecutionFilterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "edgeActionName", "executionFilterName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EdgeActionExecutionFiltersClient.ListByEdgeAction`

```go
ctx := context.TODO()
id := edgeactionexecutionfilters.NewEdgeActionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "edgeActionName")

// alternatively `client.ListByEdgeAction(ctx, id)` can be used to do batched pagination
items, err := client.ListByEdgeActionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `EdgeActionExecutionFiltersClient.Update`

```go
ctx := context.TODO()
id := edgeactionexecutionfilters.NewExecutionFilterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "edgeActionName", "executionFilterName")

payload := edgeactionexecutionfilters.EdgeActionExecutionFilterUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
