
## `github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-01-01/jobs` Documentation

The `jobs` SDK allows for interaction with Azure Resource Manager `containerapps` (API Version `2025-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-01-01/jobs"
```


### Client Initialization

```go
client := jobs.NewJobsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `JobsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := jobs.NewJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "jobName")

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
id := jobs.NewJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "jobName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `JobsClient.ExecutionsList`

```go
ctx := context.TODO()
id := jobs.NewJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "jobName")

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
id := jobs.NewJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "jobName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `JobsClient.GetDetector`

```go
ctx := context.TODO()
id := jobs.NewDetectorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "jobName", "detectorName")

read, err := client.GetDetector(ctx, id)
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
id := jobs.NewExecutionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "jobName", "executionName")

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
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

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
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `JobsClient.ListDetectors`

```go
ctx := context.TODO()
id := jobs.NewJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "jobName")

// alternatively `client.ListDetectors(ctx, id)` can be used to do batched pagination
items, err := client.ListDetectorsComplete(ctx, id)
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
id := jobs.NewJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "jobName")

read, err := client.ListSecrets(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `JobsClient.ProxyGet`

```go
ctx := context.TODO()
id := jobs.NewDetectorPropertyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "jobName", "detectorPropertyName")

read, err := client.ProxyGet(ctx, id)
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
id := jobs.NewJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "jobName")

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
id := jobs.NewExecutionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "jobName", "executionName")

if err := client.StopExecutionThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `JobsClient.StopMultipleExecutions`

```go
ctx := context.TODO()
id := jobs.NewJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "jobName")

// alternatively `client.StopMultipleExecutions(ctx, id)` can be used to do batched pagination
items, err := client.StopMultipleExecutionsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `JobsClient.Update`

```go
ctx := context.TODO()
id := jobs.NewJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "jobName")

payload := jobs.JobPatchProperties{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
