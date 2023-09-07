
## `github.com/hashicorp/go-azure-sdk/resource-manager/automation/2022-08-08/sourcecontrolsyncjobstreams` Documentation

The `sourcecontrolsyncjobstreams` SDK allows for interaction with the Azure Resource Manager Service `automation` (API Version `2022-08-08`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/automation/2022-08-08/sourcecontrolsyncjobstreams"
```


### Client Initialization

```go
client := sourcecontrolsyncjobstreams.NewSourceControlSyncJobStreamsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SourceControlSyncJobStreamsClient.Get`

```go
ctx := context.TODO()
id := sourcecontrolsyncjobstreams.NewSourceControlSyncJobStreamID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountValue", "sourceControlValue", "sourceControlSyncJobIdValue", "streamIdValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SourceControlSyncJobStreamsClient.ListBySyncJob`

```go
ctx := context.TODO()
id := sourcecontrolsyncjobstreams.NewSourceControlSyncJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountValue", "sourceControlValue", "sourceControlSyncJobIdValue")

// alternatively `client.ListBySyncJob(ctx, id, sourcecontrolsyncjobstreams.DefaultListBySyncJobOperationOptions())` can be used to do batched pagination
items, err := client.ListBySyncJobComplete(ctx, id, sourcecontrolsyncjobstreams.DefaultListBySyncJobOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
