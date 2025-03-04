
## `github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/workflows` Documentation

The `workflows` SDK allows for interaction with Azure Resource Manager `logic` (API Version `2019-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/workflows"
```


### Client Initialization

```go
client := workflows.NewWorkflowsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `WorkflowsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := workflows.NewWorkflowID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workflowName")

payload := workflows.Workflow{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WorkflowsClient.Delete`

```go
ctx := context.TODO()
id := workflows.NewWorkflowID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workflowName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WorkflowsClient.Disable`

```go
ctx := context.TODO()
id := workflows.NewWorkflowID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workflowName")

read, err := client.Disable(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WorkflowsClient.Enable`

```go
ctx := context.TODO()
id := workflows.NewWorkflowID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workflowName")

read, err := client.Enable(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WorkflowsClient.GenerateUpgradedDefinition`

```go
ctx := context.TODO()
id := workflows.NewWorkflowID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workflowName")

payload := workflows.GenerateUpgradedDefinitionParameters{
	// ...
}


read, err := client.GenerateUpgradedDefinition(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WorkflowsClient.Get`

```go
ctx := context.TODO()
id := workflows.NewWorkflowID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workflowName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WorkflowsClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id, workflows.DefaultListByResourceGroupOperationOptions())` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id, workflows.DefaultListByResourceGroupOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WorkflowsClient.ListBySubscription`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id, workflows.DefaultListBySubscriptionOperationOptions())` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id, workflows.DefaultListBySubscriptionOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WorkflowsClient.ListCallbackURL`

```go
ctx := context.TODO()
id := workflows.NewWorkflowID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workflowName")

payload := workflows.GetCallbackURLParameters{
	// ...
}


read, err := client.ListCallbackURL(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WorkflowsClient.ListSwagger`

```go
ctx := context.TODO()
id := workflows.NewWorkflowID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workflowName")

read, err := client.ListSwagger(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WorkflowsClient.Move`

```go
ctx := context.TODO()
id := workflows.NewWorkflowID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workflowName")

payload := workflows.WorkflowReference{
	// ...
}


if err := client.MoveThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `WorkflowsClient.RegenerateAccessKey`

```go
ctx := context.TODO()
id := workflows.NewWorkflowID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workflowName")

payload := workflows.RegenerateActionParameter{
	// ...
}


read, err := client.RegenerateAccessKey(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WorkflowsClient.Update`

```go
ctx := context.TODO()
id := workflows.NewWorkflowID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workflowName")

read, err := client.Update(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WorkflowsClient.ValidateByLocation`

```go
ctx := context.TODO()
id := workflows.NewLocationWorkflowID("12345678-1234-9876-4563-123456789012", "example-resource-group", "locationName", "workflowName")

payload := workflows.Workflow{
	// ...
}


read, err := client.ValidateByLocation(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WorkflowsClient.ValidateByResourceGroup`

```go
ctx := context.TODO()
id := workflows.NewWorkflowID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workflowName")

payload := workflows.Workflow{
	// ...
}


read, err := client.ValidateByResourceGroup(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
