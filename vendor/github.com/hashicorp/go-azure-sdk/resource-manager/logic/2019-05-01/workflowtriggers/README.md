
## `github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/workflowtriggers` Documentation

The `workflowtriggers` SDK allows for interaction with Azure Resource Manager `logic` (API Version `2019-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/workflowtriggers"
```


### Client Initialization

```go
client := workflowtriggers.NewWorkflowTriggersClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `WorkflowTriggersClient.Get`

```go
ctx := context.TODO()
id := workflowtriggers.NewTriggerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workflowName", "triggerName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WorkflowTriggersClient.GetSchemaJson`

```go
ctx := context.TODO()
id := workflowtriggers.NewTriggerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workflowName", "triggerName")

read, err := client.GetSchemaJson(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WorkflowTriggersClient.List`

```go
ctx := context.TODO()
id := workflowtriggers.NewWorkflowID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workflowName")

// alternatively `client.List(ctx, id, workflowtriggers.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, workflowtriggers.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WorkflowTriggersClient.ListCallbackURL`

```go
ctx := context.TODO()
id := workflowtriggers.NewTriggerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workflowName", "triggerName")

read, err := client.ListCallbackURL(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WorkflowTriggersClient.Reset`

```go
ctx := context.TODO()
id := workflowtriggers.NewTriggerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workflowName", "triggerName")

read, err := client.Reset(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WorkflowTriggersClient.Run`

```go
ctx := context.TODO()
id := workflowtriggers.NewTriggerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workflowName", "triggerName")

read, err := client.Run(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WorkflowTriggersClient.SetState`

```go
ctx := context.TODO()
id := workflowtriggers.NewTriggerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workflowName", "triggerName")

payload := workflowtriggers.SetTriggerStateActionDefinition{
	// ...
}


read, err := client.SetState(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WorkflowTriggersClient.WorkflowVersionTriggersListCallbackURL`

```go
ctx := context.TODO()
id := workflowtriggers.NewVersionTriggerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workflowName", "versionId", "triggerName")

payload := workflowtriggers.GetCallbackURLParameters{
	// ...
}


read, err := client.WorkflowVersionTriggersListCallbackURL(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
