
## `github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/listdeletedrunbooks` Documentation

The `listdeletedrunbooks` SDK allows for interaction with Azure Resource Manager `automation` (API Version `2024-10-23`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/listdeletedrunbooks"
```


### Client Initialization

```go
client := listdeletedrunbooks.NewListDeletedRunbooksClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ListDeletedRunbooksClient.AutomationAccountListDeletedRunbooks`

```go
ctx := context.TODO()
id := listdeletedrunbooks.NewAutomationAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName")

// alternatively `client.AutomationAccountListDeletedRunbooks(ctx, id)` can be used to do batched pagination
items, err := client.AutomationAccountListDeletedRunbooksComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
