
## `github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/sourcecontrolsyncjob` Documentation

The `sourcecontrolsyncjob` SDK allows for interaction with Azure Resource Manager `automation` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/sourcecontrolsyncjob"
```


### Client Initialization

```go
client := sourcecontrolsyncjob.NewSourceControlSyncJobClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SourceControlSyncJobClient.Create`

```go
ctx := context.TODO()
id := sourcecontrolsyncjob.NewSourceControlSyncJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "sourceControlName", "sourceControlSyncJobId")

payload := sourcecontrolsyncjob.SourceControlSyncJobCreateParameters{
	// ...
}


read, err := client.Create(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SourceControlSyncJobClient.Get`

```go
ctx := context.TODO()
id := sourcecontrolsyncjob.NewSourceControlSyncJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "sourceControlName", "sourceControlSyncJobId")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SourceControlSyncJobClient.ListByAutomationAccount`

```go
ctx := context.TODO()
id := sourcecontrolsyncjob.NewSourceControlID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "sourceControlName")

// alternatively `client.ListByAutomationAccount(ctx, id, sourcecontrolsyncjob.DefaultListByAutomationAccountOperationOptions())` can be used to do batched pagination
items, err := client.ListByAutomationAccountComplete(ctx, id, sourcecontrolsyncjob.DefaultListByAutomationAccountOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
