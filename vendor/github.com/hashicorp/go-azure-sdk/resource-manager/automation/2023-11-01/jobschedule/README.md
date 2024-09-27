
## `github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/jobschedule` Documentation

The `jobschedule` SDK allows for interaction with Azure Resource Manager `automation` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/jobschedule"
```


### Client Initialization

```go
client := jobschedule.NewJobScheduleClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `JobScheduleClient.Create`

```go
ctx := context.TODO()
id := jobschedule.NewJobScheduleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "jobScheduleId")

payload := jobschedule.JobScheduleCreateParameters{
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


### Example Usage: `JobScheduleClient.Delete`

```go
ctx := context.TODO()
id := jobschedule.NewJobScheduleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "jobScheduleId")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `JobScheduleClient.Get`

```go
ctx := context.TODO()
id := jobschedule.NewJobScheduleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "jobScheduleId")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `JobScheduleClient.ListByAutomationAccount`

```go
ctx := context.TODO()
id := jobschedule.NewAutomationAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName")

// alternatively `client.ListByAutomationAccount(ctx, id, jobschedule.DefaultListByAutomationAccountOperationOptions())` can be used to do batched pagination
items, err := client.ListByAutomationAccountComplete(ctx, id, jobschedule.DefaultListByAutomationAccountOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
