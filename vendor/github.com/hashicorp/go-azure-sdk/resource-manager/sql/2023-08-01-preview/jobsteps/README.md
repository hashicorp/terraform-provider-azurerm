
## `github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/jobsteps` Documentation

The `jobsteps` SDK allows for interaction with Azure Resource Manager `sql` (API Version `2023-08-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/jobsteps"
```


### Client Initialization

```go
client := jobsteps.NewJobStepsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `JobStepsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := jobsteps.NewStepID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverName", "jobAgentName", "jobName", "stepName")

payload := jobsteps.JobStep{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `JobStepsClient.Delete`

```go
ctx := context.TODO()
id := jobsteps.NewStepID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverName", "jobAgentName", "jobName", "stepName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `JobStepsClient.Get`

```go
ctx := context.TODO()
id := jobsteps.NewStepID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverName", "jobAgentName", "jobName", "stepName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `JobStepsClient.GetByVersion`

```go
ctx := context.TODO()
id := jobsteps.NewVersionStepID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverName", "jobAgentName", "jobName", "versionName", "stepName")

read, err := client.GetByVersion(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `JobStepsClient.ListByJob`

```go
ctx := context.TODO()
id := jobsteps.NewJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverName", "jobAgentName", "jobName")

// alternatively `client.ListByJob(ctx, id)` can be used to do batched pagination
items, err := client.ListByJobComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `JobStepsClient.ListByVersion`

```go
ctx := context.TODO()
id := jobsteps.NewVersionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverName", "jobAgentName", "jobName", "versionName")

// alternatively `client.ListByVersion(ctx, id)` can be used to do batched pagination
items, err := client.ListByVersionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
