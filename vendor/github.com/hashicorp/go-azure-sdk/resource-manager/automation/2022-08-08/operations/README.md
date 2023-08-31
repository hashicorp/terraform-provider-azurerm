
## `github.com/hashicorp/go-azure-sdk/resource-manager/automation/2022-08-08/operations` Documentation

The `operations` SDK allows for interaction with the Azure Resource Manager Service `automation` (API Version `2022-08-08`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/automation/2022-08-08/operations"
```


### Client Initialization

```go
client := operations.NewOperationsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `OperationsClient.ConvertGraphRunbookContent`

```go
ctx := context.TODO()
id := operations.NewAutomationAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountValue")

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
