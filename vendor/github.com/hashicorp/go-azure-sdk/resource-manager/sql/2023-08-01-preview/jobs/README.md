
## `github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/jobs` Documentation

The `jobs` SDK allows for interaction with Azure Resource Manager `sql` (API Version `2023-08-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/jobs"
```


### Client Initialization

```go
client := jobs.NewJobsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `JobsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := jobs.NewJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverName", "jobAgentName", "jobName")

payload := jobs.Job{
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


### Example Usage: `JobsClient.Delete`

```go
ctx := context.TODO()
id := jobs.NewJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverName", "jobAgentName", "jobName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `JobsClient.Get`

```go
ctx := context.TODO()
id := jobs.NewJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverName", "jobAgentName", "jobName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `JobsClient.ListByAgent`

```go
ctx := context.TODO()
id := jobs.NewJobAgentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverName", "jobAgentName")

// alternatively `client.ListByAgent(ctx, id)` can be used to do batched pagination
items, err := client.ListByAgentComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
