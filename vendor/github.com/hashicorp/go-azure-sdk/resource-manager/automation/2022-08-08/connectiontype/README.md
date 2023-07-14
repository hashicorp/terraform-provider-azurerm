
## `github.com/hashicorp/go-azure-sdk/resource-manager/automation/2022-08-08/connectiontype` Documentation

The `connectiontype` SDK allows for interaction with the Azure Resource Manager Service `automation` (API Version `2022-08-08`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/automation/2022-08-08/connectiontype"
```


### Client Initialization

```go
client := connectiontype.NewConnectionTypeClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ConnectionTypeClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := connectiontype.NewConnectionTypeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountValue", "connectionTypeValue")

payload := connectiontype.ConnectionTypeCreateOrUpdateParameters{
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


### Example Usage: `ConnectionTypeClient.Delete`

```go
ctx := context.TODO()
id := connectiontype.NewConnectionTypeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountValue", "connectionTypeValue")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ConnectionTypeClient.Get`

```go
ctx := context.TODO()
id := connectiontype.NewConnectionTypeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountValue", "connectionTypeValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ConnectionTypeClient.ListByAutomationAccount`

```go
ctx := context.TODO()
id := connectiontype.NewAutomationAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountValue")

// alternatively `client.ListByAutomationAccount(ctx, id)` can be used to do batched pagination
items, err := client.ListByAutomationAccountComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
