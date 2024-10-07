
## `github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/jobstream` Documentation

The `jobstream` SDK allows for interaction with Azure Resource Manager `automation` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/jobstream"
```


### Client Initialization

```go
client := jobstream.NewJobStreamClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `JobStreamClient.Get`

```go
ctx := context.TODO()
id := jobstream.NewStreamID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "jobName", "jobStreamId")

read, err := client.Get(ctx, id, jobstream.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `JobStreamClient.ListByJob`

```go
ctx := context.TODO()
id := jobstream.NewJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "jobName")

// alternatively `client.ListByJob(ctx, id, jobstream.DefaultListByJobOperationOptions())` can be used to do batched pagination
items, err := client.ListByJobComplete(ctx, id, jobstream.DefaultListByJobOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
