
## `github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/jobexecutions` Documentation

The `jobexecutions` SDK allows for interaction with Azure Resource Manager `sql` (API Version `2023-08-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/jobexecutions"
```


### Client Initialization

```go
client := jobexecutions.NewJobExecutionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `JobExecutionsClient.Cancel`

```go
ctx := context.TODO()
id := jobexecutions.NewExecutionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverName", "jobAgentName", "jobName", "jobExecutionId")

read, err := client.Cancel(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `JobExecutionsClient.Create`

```go
ctx := context.TODO()
id := jobexecutions.NewJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverName", "jobAgentName", "jobName")

if err := client.CreateThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `JobExecutionsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := jobexecutions.NewExecutionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverName", "jobAgentName", "jobName", "jobExecutionId")

if err := client.CreateOrUpdateThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `JobExecutionsClient.Get`

```go
ctx := context.TODO()
id := jobexecutions.NewExecutionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverName", "jobAgentName", "jobName", "jobExecutionId")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `JobExecutionsClient.ListByAgent`

```go
ctx := context.TODO()
id := jobexecutions.NewJobAgentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverName", "jobAgentName")

// alternatively `client.ListByAgent(ctx, id, jobexecutions.DefaultListByAgentOperationOptions())` can be used to do batched pagination
items, err := client.ListByAgentComplete(ctx, id, jobexecutions.DefaultListByAgentOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `JobExecutionsClient.ListByJob`

```go
ctx := context.TODO()
id := jobexecutions.NewJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverName", "jobAgentName", "jobName")

// alternatively `client.ListByJob(ctx, id, jobexecutions.DefaultListByJobOperationOptions())` can be used to do batched pagination
items, err := client.ListByJobComplete(ctx, id, jobexecutions.DefaultListByJobOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
