
## `github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/jobtargetgroups` Documentation

The `jobtargetgroups` SDK allows for interaction with Azure Resource Manager `sql` (API Version `2023-08-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/jobtargetgroups"
```


### Client Initialization

```go
client := jobtargetgroups.NewJobTargetGroupsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `JobTargetGroupsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := jobtargetgroups.NewTargetGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverName", "jobAgentName", "targetGroupName")

payload := jobtargetgroups.JobTargetGroup{
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


### Example Usage: `JobTargetGroupsClient.Delete`

```go
ctx := context.TODO()
id := jobtargetgroups.NewTargetGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverName", "jobAgentName", "targetGroupName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `JobTargetGroupsClient.Get`

```go
ctx := context.TODO()
id := jobtargetgroups.NewTargetGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverName", "jobAgentName", "targetGroupName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `JobTargetGroupsClient.ListByAgent`

```go
ctx := context.TODO()
id := jobtargetgroups.NewJobAgentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverName", "jobAgentName")

// alternatively `client.ListByAgent(ctx, id)` can be used to do batched pagination
items, err := client.ListByAgentComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
