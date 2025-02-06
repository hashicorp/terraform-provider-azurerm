
## `github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/operations` Documentation

The `operations` SDK allows for interaction with Azure Resource Manager `automation` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/operations"
```


### Client Initialization

```go
client := operations.NewOperationsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `OperationsClient.ConvertGraphRunbookContent`

```go
ctx := context.TODO()
id := operations.NewAutomationAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName")

payload := operations.GraphicalRunbookContent{
	// ...
}


read, err := client.ConvertGraphRunbookContent(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
