
## `github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/listkeys` Documentation

The `listkeys` SDK allows for interaction with Azure Resource Manager `automation` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/listkeys"
```


### Client Initialization

```go
client := listkeys.NewListKeysClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ListKeysClient.KeysListByAutomationAccount`

```go
ctx := context.TODO()
id := listkeys.NewAutomationAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName")

read, err := client.KeysListByAutomationAccount(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
