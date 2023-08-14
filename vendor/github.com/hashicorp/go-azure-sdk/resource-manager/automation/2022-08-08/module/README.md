
## `github.com/hashicorp/go-azure-sdk/resource-manager/automation/2022-08-08/module` Documentation

The `module` SDK allows for interaction with the Azure Resource Manager Service `automation` (API Version `2022-08-08`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/automation/2022-08-08/module"
```


### Client Initialization

```go
client := module.NewModuleClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ModuleClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := module.NewModuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountValue", "moduleValue")

payload := module.ModuleCreateOrUpdateParameters{
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


### Example Usage: `ModuleClient.Delete`

```go
ctx := context.TODO()
id := module.NewModuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountValue", "moduleValue")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ModuleClient.Get`

```go
ctx := context.TODO()
id := module.NewModuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountValue", "moduleValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ModuleClient.ListByAutomationAccount`

```go
ctx := context.TODO()
id := module.NewAutomationAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountValue")

// alternatively `client.ListByAutomationAccount(ctx, id)` can be used to do batched pagination
items, err := client.ListByAutomationAccountComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ModuleClient.Update`

```go
ctx := context.TODO()
id := module.NewModuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountValue", "moduleValue")

payload := module.ModuleUpdateParameters{
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
