
## `github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/jobagents` Documentation

The `jobagents` SDK allows for interaction with Azure Resource Manager `sql` (API Version `2023-08-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/jobagents"
```


### Client Initialization

```go
client := jobagents.NewJobAgentsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `JobAgentsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := jobagents.NewJobAgentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverName", "jobAgentName")

payload := jobagents.JobAgent{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `JobAgentsClient.Delete`

```go
ctx := context.TODO()
id := jobagents.NewJobAgentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverName", "jobAgentName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `JobAgentsClient.Get`

```go
ctx := context.TODO()
id := jobagents.NewJobAgentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverName", "jobAgentName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `JobAgentsClient.ListByServer`

```go
ctx := context.TODO()
id := commonids.NewSqlServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverName")

// alternatively `client.ListByServer(ctx, id)` can be used to do batched pagination
items, err := client.ListByServerComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `JobAgentsClient.Update`

```go
ctx := context.TODO()
id := jobagents.NewJobAgentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverName", "jobAgentName")

payload := jobagents.JobAgentUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
