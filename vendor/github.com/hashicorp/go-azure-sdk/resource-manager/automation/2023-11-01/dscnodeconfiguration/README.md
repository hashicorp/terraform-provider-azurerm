
## `github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/dscnodeconfiguration` Documentation

The `dscnodeconfiguration` SDK allows for interaction with Azure Resource Manager `automation` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/dscnodeconfiguration"
```


### Client Initialization

```go
client := dscnodeconfiguration.NewDscNodeConfigurationClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DscNodeConfigurationClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := dscnodeconfiguration.NewNodeConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "nodeConfigurationName")

payload := dscnodeconfiguration.DscNodeConfigurationCreateOrUpdateParameters{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `DscNodeConfigurationClient.Delete`

```go
ctx := context.TODO()
id := dscnodeconfiguration.NewNodeConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "nodeConfigurationName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DscNodeConfigurationClient.Get`

```go
ctx := context.TODO()
id := dscnodeconfiguration.NewNodeConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "nodeConfigurationName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DscNodeConfigurationClient.ListByAutomationAccount`

```go
ctx := context.TODO()
id := dscnodeconfiguration.NewAutomationAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName")

// alternatively `client.ListByAutomationAccount(ctx, id, dscnodeconfiguration.DefaultListByAutomationAccountOperationOptions())` can be used to do batched pagination
items, err := client.ListByAutomationAccountComplete(ctx, id, dscnodeconfiguration.DefaultListByAutomationAccountOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
