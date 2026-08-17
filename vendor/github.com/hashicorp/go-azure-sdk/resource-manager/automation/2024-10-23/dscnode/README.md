
## `github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/dscnode` Documentation

The `dscnode` SDK allows for interaction with Azure Resource Manager `automation` (API Version `2024-10-23`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/dscnode"
```


### Client Initialization

```go
client := dscnode.NewDscNodeClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DscNodeClient.Delete`

```go
ctx := context.TODO()
id := dscnode.NewNodeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "nodeId")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DscNodeClient.Get`

```go
ctx := context.TODO()
id := dscnode.NewNodeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "nodeId")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DscNodeClient.ListByAutomationAccount`

```go
ctx := context.TODO()
id := dscnode.NewAutomationAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName")

// alternatively `client.ListByAutomationAccount(ctx, id, dscnode.DefaultListByAutomationAccountOperationOptions())` can be used to do batched pagination
items, err := client.ListByAutomationAccountComplete(ctx, id, dscnode.DefaultListByAutomationAccountOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DscNodeClient.Update`

```go
ctx := context.TODO()
id := dscnode.NewNodeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "nodeId")

payload := dscnode.DscNodeUpdateParameters{
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
