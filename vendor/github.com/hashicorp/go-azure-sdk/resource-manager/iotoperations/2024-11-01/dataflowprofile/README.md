
## `github.com/hashicorp/go-azure-sdk/resource-manager/iotoperations/2024-11-01/dataflowprofile` Documentation

The `dataflowprofile` SDK allows for interaction with Azure Resource Manager `iotoperations` (API Version `2024-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/iotoperations/2024-11-01/dataflowprofile"
```


### Client Initialization

```go
client := dataflowprofile.NewDataflowProfileClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DataflowProfileClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := dataflowprofile.NewDataflowProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "instanceName", "dataflowProfileName")

payload := dataflowprofile.DataflowProfileResource{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `DataflowProfileClient.Delete`

```go
ctx := context.TODO()
id := dataflowprofile.NewDataflowProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "instanceName", "dataflowProfileName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `DataflowProfileClient.Get`

```go
ctx := context.TODO()
id := dataflowprofile.NewDataflowProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "instanceName", "dataflowProfileName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DataflowProfileClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := dataflowprofile.NewInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "instanceName")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
