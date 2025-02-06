
## `github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/hybridrunbookworkergroup` Documentation

The `hybridrunbookworkergroup` SDK allows for interaction with Azure Resource Manager `automation` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/hybridrunbookworkergroup"
```


### Client Initialization

```go
client := hybridrunbookworkergroup.NewHybridRunbookWorkerGroupClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `HybridRunbookWorkerGroupClient.Create`

```go
ctx := context.TODO()
id := hybridrunbookworkergroup.NewHybridRunbookWorkerGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "hybridRunbookWorkerGroupName")

payload := hybridrunbookworkergroup.HybridRunbookWorkerGroupCreateOrUpdateParameters{
	// ...
}


read, err := client.Create(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `HybridRunbookWorkerGroupClient.Delete`

```go
ctx := context.TODO()
id := hybridrunbookworkergroup.NewHybridRunbookWorkerGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "hybridRunbookWorkerGroupName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `HybridRunbookWorkerGroupClient.Get`

```go
ctx := context.TODO()
id := hybridrunbookworkergroup.NewHybridRunbookWorkerGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "hybridRunbookWorkerGroupName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `HybridRunbookWorkerGroupClient.ListByAutomationAccount`

```go
ctx := context.TODO()
id := hybridrunbookworkergroup.NewAutomationAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName")

// alternatively `client.ListByAutomationAccount(ctx, id, hybridrunbookworkergroup.DefaultListByAutomationAccountOperationOptions())` can be used to do batched pagination
items, err := client.ListByAutomationAccountComplete(ctx, id, hybridrunbookworkergroup.DefaultListByAutomationAccountOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `HybridRunbookWorkerGroupClient.Update`

```go
ctx := context.TODO()
id := hybridrunbookworkergroup.NewHybridRunbookWorkerGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "hybridRunbookWorkerGroupName")

payload := hybridrunbookworkergroup.HybridRunbookWorkerGroupCreateOrUpdateParameters{
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
