
## `github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/jobstepexecutions` Documentation

The `jobstepexecutions` SDK allows for interaction with Azure Resource Manager `sql` (API Version `2023-08-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/jobstepexecutions"
```


### Client Initialization

```go
client := jobstepexecutions.NewJobStepExecutionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `JobStepExecutionsClient.Get`

```go
ctx := context.TODO()
id := jobstepexecutions.NewExecutionStepID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverName", "jobAgentName", "jobName", "jobExecutionId", "stepName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `JobStepExecutionsClient.ListByJobExecution`

```go
ctx := context.TODO()
id := jobstepexecutions.NewExecutionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverName", "jobAgentName", "jobName", "jobExecutionId")

// alternatively `client.ListByJobExecution(ctx, id, jobstepexecutions.DefaultListByJobExecutionOperationOptions())` can be used to do batched pagination
items, err := client.ListByJobExecutionComplete(ctx, id, jobstepexecutions.DefaultListByJobExecutionOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
