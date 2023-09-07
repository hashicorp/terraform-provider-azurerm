
## `github.com/hashicorp/go-azure-sdk/resource-manager/automation/2022-08-08/job` Documentation

The `job` SDK allows for interaction with the Azure Resource Manager Service `automation` (API Version `2022-08-08`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/automation/2022-08-08/job"
```


### Client Initialization

```go
client := job.NewJobClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `JobClient.Create`

```go
ctx := context.TODO()
id := job.NewJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountValue", "jobValue")

payload := job.JobCreateParameters{
	// ...
}


read, err := client.Create(ctx, id, payload, job.DefaultCreateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `JobClient.Get`

```go
ctx := context.TODO()
id := job.NewJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountValue", "jobValue")

read, err := client.Get(ctx, id, job.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `JobClient.GetOutput`

```go
ctx := context.TODO()
id := job.NewJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountValue", "jobValue")

read, err := client.GetOutput(ctx, id, job.DefaultGetOutputOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `JobClient.GetRunbookContent`

```go
ctx := context.TODO()
id := job.NewJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountValue", "jobValue")

read, err := client.GetRunbookContent(ctx, id, job.DefaultGetRunbookContentOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `JobClient.ListByAutomationAccount`

```go
ctx := context.TODO()
id := job.NewAutomationAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountValue")

// alternatively `client.ListByAutomationAccount(ctx, id, job.DefaultListByAutomationAccountOperationOptions())` can be used to do batched pagination
items, err := client.ListByAutomationAccountComplete(ctx, id, job.DefaultListByAutomationAccountOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `JobClient.Resume`

```go
ctx := context.TODO()
id := job.NewJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountValue", "jobValue")

read, err := client.Resume(ctx, id, job.DefaultResumeOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `JobClient.Stop`

```go
ctx := context.TODO()
id := job.NewJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountValue", "jobValue")

read, err := client.Stop(ctx, id, job.DefaultStopOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `JobClient.Suspend`

```go
ctx := context.TODO()
id := job.NewJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountValue", "jobValue")

read, err := client.Suspend(ctx, id, job.DefaultSuspendOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
