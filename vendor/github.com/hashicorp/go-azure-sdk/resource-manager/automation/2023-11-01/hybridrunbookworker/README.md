
## `github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/hybridrunbookworker` Documentation

The `hybridrunbookworker` SDK allows for interaction with Azure Resource Manager `automation` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/hybridrunbookworker"
```


### Client Initialization

```go
client := hybridrunbookworker.NewHybridRunbookWorkerClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `HybridRunbookWorkerClient.Create`

```go
ctx := context.TODO()
id := hybridrunbookworker.NewHybridRunbookWorkerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "hybridRunbookWorkerGroupName", "hybridRunbookWorkerId")

payload := hybridrunbookworker.HybridRunbookWorkerCreateParameters{
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


### Example Usage: `HybridRunbookWorkerClient.Delete`

```go
ctx := context.TODO()
id := hybridrunbookworker.NewHybridRunbookWorkerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "hybridRunbookWorkerGroupName", "hybridRunbookWorkerId")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `HybridRunbookWorkerClient.Get`

```go
ctx := context.TODO()
id := hybridrunbookworker.NewHybridRunbookWorkerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "hybridRunbookWorkerGroupName", "hybridRunbookWorkerId")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `HybridRunbookWorkerClient.ListByHybridRunbookWorkerGroup`

```go
ctx := context.TODO()
id := hybridrunbookworker.NewHybridRunbookWorkerGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "hybridRunbookWorkerGroupName")

// alternatively `client.ListByHybridRunbookWorkerGroup(ctx, id, hybridrunbookworker.DefaultListByHybridRunbookWorkerGroupOperationOptions())` can be used to do batched pagination
items, err := client.ListByHybridRunbookWorkerGroupComplete(ctx, id, hybridrunbookworker.DefaultListByHybridRunbookWorkerGroupOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `HybridRunbookWorkerClient.Move`

```go
ctx := context.TODO()
id := hybridrunbookworker.NewHybridRunbookWorkerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "hybridRunbookWorkerGroupName", "hybridRunbookWorkerId")

payload := hybridrunbookworker.HybridRunbookWorkerMoveParameters{
	// ...
}


read, err := client.Move(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
