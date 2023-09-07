
## `github.com/hashicorp/go-azure-sdk/resource-manager/automation/2022-08-08/statistics` Documentation

The `statistics` SDK allows for interaction with the Azure Resource Manager Service `automation` (API Version `2022-08-08`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/automation/2022-08-08/statistics"
```


### Client Initialization

```go
client := statistics.NewStatisticsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `StatisticsClient.ListByAutomationAccount`

```go
ctx := context.TODO()
id := statistics.NewAutomationAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountValue")

read, err := client.ListByAutomationAccount(ctx, id, statistics.DefaultListByAutomationAccountOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
