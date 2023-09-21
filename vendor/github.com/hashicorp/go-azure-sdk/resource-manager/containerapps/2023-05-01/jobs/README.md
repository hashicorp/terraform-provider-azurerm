
## `github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2023-05-01/jobs` Documentation

The `jobs` SDK allows for interaction with the Azure Resource Manager Service `containerapps` (API Version `2023-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2023-05-01/jobs"
```


### Client Initialization

```go
client := jobs.NewJobsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `JobsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := jobs.NewJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "jobValue")

payload := jobs.Job{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `JobsClient.Delete`

```go
ctx := context.TODO()
id := jobs.NewJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "jobValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `JobsClient.ExecutionsList`

```go
ctx := context.TODO()
id := jobs.NewJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "jobValue")

// alternatively `client.ExecutionsList(ctx, id, jobs.DefaultExecutionsListOperationOptions())` can be used to do batched pagination
items, err := client.ExecutionsListComplete(ctx, id, jobs.DefaultExecutionsListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `JobsClient.Get`

```go
ctx := context.TODO()
id := jobs.NewJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "jobValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `JobsClient.JobExecution`

```go
ctx := context.TODO()
id := jobs.NewExecutionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "jobValue", "executionValue")

read, err := client.JobExecution(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `JobsClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := jobs.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `JobsClient.ListBySubscription`

```go
ctx := context.TODO()
id := jobs.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `JobsClient.ListSecrets`

```go
ctx := context.TODO()
id := jobs.NewJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "jobValue")

read, err := client.ListSecrets(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `JobsClient.Start`

```go
ctx := context.TODO()
id := jobs.NewJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "jobValue")

payload := jobs.JobExecutionTemplate{
	// ...
}


if err := client.StartThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `JobsClient.StopExecution`

```go
ctx := context.TODO()
id := jobs.NewExecutionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "jobValue", "executionValue")

if err := client.StopExecutionThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `JobsClient.StopMultipleExecutions`

```go
ctx := context.TODO()
id := jobs.NewJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "jobValue")

if err := client.StopMultipleExecutionsThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `JobsClient.Update`

```go
ctx := context.TODO()
id := jobs.NewJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "jobValue")

payload := jobs.JobPatchProperties{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
