
## `github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/nodecountinformation` Documentation

The `nodecountinformation` SDK allows for interaction with Azure Resource Manager `automation` (API Version `2024-10-23`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/nodecountinformation"
```


### Client Initialization

```go
client := nodecountinformation.NewNodeCountInformationClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `NodeCountInformationClient.Get`

```go
ctx := context.TODO()
id := nodecountinformation.NewCountTypeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "nodeconfiguration")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
