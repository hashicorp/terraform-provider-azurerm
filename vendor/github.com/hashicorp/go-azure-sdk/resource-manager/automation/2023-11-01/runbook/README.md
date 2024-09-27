
## `github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/runbook` Documentation

The `runbook` SDK allows for interaction with Azure Resource Manager `automation` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/runbook"
```


### Client Initialization

```go
client := runbook.NewRunbookClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `RunbookClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := runbook.NewRunbookID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "runbookName")

payload := runbook.RunbookCreateOrUpdateParameters{
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


### Example Usage: `RunbookClient.Delete`

```go
ctx := context.TODO()
id := runbook.NewRunbookID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "runbookName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RunbookClient.Get`

```go
ctx := context.TODO()
id := runbook.NewRunbookID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "runbookName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RunbookClient.GetContent`

```go
ctx := context.TODO()
id := runbook.NewRunbookID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "runbookName")

read, err := client.GetContent(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RunbookClient.ListByAutomationAccount`

```go
ctx := context.TODO()
id := runbook.NewAutomationAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName")

// alternatively `client.ListByAutomationAccount(ctx, id)` can be used to do batched pagination
items, err := client.ListByAutomationAccountComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `RunbookClient.Publish`

```go
ctx := context.TODO()
id := runbook.NewRunbookID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "runbookName")

if err := client.PublishThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `RunbookClient.Update`

```go
ctx := context.TODO()
id := runbook.NewRunbookID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "runbookName")

payload := runbook.RunbookUpdateParameters{
	// ...
}


read, err := client.Update(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
