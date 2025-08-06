
## `github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/connection` Documentation

The `connection` SDK allows for interaction with Azure Resource Manager `automation` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/connection"
```


### Client Initialization

```go
client := connection.NewConnectionClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ConnectionClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := connection.NewConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "connectionName")

payload := connection.ConnectionCreateOrUpdateParameters{
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


### Example Usage: `ConnectionClient.Delete`

```go
ctx := context.TODO()
id := connection.NewConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "connectionName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ConnectionClient.Get`

```go
ctx := context.TODO()
id := connection.NewConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "connectionName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ConnectionClient.ListByAutomationAccount`

```go
ctx := context.TODO()
id := connection.NewAutomationAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName")

// alternatively `client.ListByAutomationAccount(ctx, id)` can be used to do batched pagination
items, err := client.ListByAutomationAccountComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ConnectionClient.Update`

```go
ctx := context.TODO()
id := connection.NewConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "connectionName")

payload := connection.ConnectionUpdateParameters{
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
