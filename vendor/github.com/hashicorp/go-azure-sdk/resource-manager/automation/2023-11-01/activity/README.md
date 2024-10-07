
## `github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/activity` Documentation

The `activity` SDK allows for interaction with Azure Resource Manager `automation` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/activity"
```


### Client Initialization

```go
client := activity.NewActivityClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ActivityClient.Get`

```go
ctx := context.TODO()
id := activity.NewActivityID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "moduleName", "activityName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ActivityClient.ListByModule`

```go
ctx := context.TODO()
id := activity.NewModuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "moduleName")

// alternatively `client.ListByModule(ctx, id)` can be used to do batched pagination
items, err := client.ListByModuleComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
