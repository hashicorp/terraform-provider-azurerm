
## `github.com/hashicorp/go-azure-sdk/resource-manager/automation/2022-08-08/dscconfiguration` Documentation

The `dscconfiguration` SDK allows for interaction with the Azure Resource Manager Service `automation` (API Version `2022-08-08`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/automation/2022-08-08/dscconfiguration"
```


### Client Initialization

```go
client := dscconfiguration.NewDscConfigurationClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DscConfigurationClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := dscconfiguration.NewConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountValue", "configurationValue")

payload := dscconfiguration.DscConfigurationCreateOrUpdateParameters{
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


### Example Usage: `DscConfigurationClient.Delete`

```go
ctx := context.TODO()
id := dscconfiguration.NewConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountValue", "configurationValue")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DscConfigurationClient.Get`

```go
ctx := context.TODO()
id := dscconfiguration.NewConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountValue", "configurationValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DscConfigurationClient.GetContent`

```go
ctx := context.TODO()
id := dscconfiguration.NewConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountValue", "configurationValue")

read, err := client.GetContent(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DscConfigurationClient.ListByAutomationAccount`

```go
ctx := context.TODO()
id := dscconfiguration.NewAutomationAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountValue")

// alternatively `client.ListByAutomationAccount(ctx, id, dscconfiguration.DefaultListByAutomationAccountOperationOptions())` can be used to do batched pagination
items, err := client.ListByAutomationAccountComplete(ctx, id, dscconfiguration.DefaultListByAutomationAccountOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DscConfigurationClient.Update`

```go
ctx := context.TODO()
id := dscconfiguration.NewConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountValue", "configurationValue")

payload := dscconfiguration.DscConfigurationUpdateParameters{
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
