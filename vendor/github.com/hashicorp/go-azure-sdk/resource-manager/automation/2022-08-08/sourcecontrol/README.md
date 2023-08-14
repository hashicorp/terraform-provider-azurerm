
## `github.com/hashicorp/go-azure-sdk/resource-manager/automation/2022-08-08/sourcecontrol` Documentation

The `sourcecontrol` SDK allows for interaction with the Azure Resource Manager Service `automation` (API Version `2022-08-08`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/automation/2022-08-08/sourcecontrol"
```


### Client Initialization

```go
client := sourcecontrol.NewSourceControlClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SourceControlClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := sourcecontrol.NewSourceControlID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountValue", "sourceControlValue")

payload := sourcecontrol.SourceControlCreateOrUpdateParameters{
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


### Example Usage: `SourceControlClient.Delete`

```go
ctx := context.TODO()
id := sourcecontrol.NewSourceControlID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountValue", "sourceControlValue")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SourceControlClient.Get`

```go
ctx := context.TODO()
id := sourcecontrol.NewSourceControlID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountValue", "sourceControlValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SourceControlClient.ListByAutomationAccount`

```go
ctx := context.TODO()
id := sourcecontrol.NewAutomationAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountValue")

// alternatively `client.ListByAutomationAccount(ctx, id, sourcecontrol.DefaultListByAutomationAccountOperationOptions())` can be used to do batched pagination
items, err := client.ListByAutomationAccountComplete(ctx, id, sourcecontrol.DefaultListByAutomationAccountOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `SourceControlClient.Update`

```go
ctx := context.TODO()
id := sourcecontrol.NewSourceControlID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountValue", "sourceControlValue")

payload := sourcecontrol.SourceControlUpdateParameters{
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
