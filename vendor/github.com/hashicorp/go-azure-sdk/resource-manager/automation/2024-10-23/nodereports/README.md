
## `github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/nodereports` Documentation

The `nodereports` SDK allows for interaction with Azure Resource Manager `automation` (API Version `2024-10-23`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/nodereports"
```


### Client Initialization

```go
client := nodereports.NewNodeReportsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `NodeReportsClient.Get`

```go
ctx := context.TODO()
id := nodereports.NewReportID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "nodeId", "reportId")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NodeReportsClient.GetContent`

```go
ctx := context.TODO()
id := nodereports.NewReportID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "nodeId", "reportId")

read, err := client.GetContent(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NodeReportsClient.ListByNode`

```go
ctx := context.TODO()
id := nodereports.NewNodeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "nodeId")

// alternatively `client.ListByNode(ctx, id, nodereports.DefaultListByNodeOperationOptions())` can be used to do batched pagination
items, err := client.ListByNodeComplete(ctx, id, nodereports.DefaultListByNodeOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
