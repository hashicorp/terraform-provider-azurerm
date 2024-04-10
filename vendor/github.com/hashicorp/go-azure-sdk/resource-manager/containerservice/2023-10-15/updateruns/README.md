
## `github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2023-10-15/updateruns` Documentation

The `updateruns` SDK allows for interaction with the Azure Resource Manager Service `containerservice` (API Version `2023-10-15`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2023-10-15/updateruns"
```


### Client Initialization

```go
client := updateruns.NewUpdateRunsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `UpdateRunsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := updateruns.NewUpdateRunID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fleetValue", "updateRunValue")

payload := updateruns.UpdateRun{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload, updateruns.DefaultCreateOrUpdateOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `UpdateRunsClient.Delete`

```go
ctx := context.TODO()
id := updateruns.NewUpdateRunID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fleetValue", "updateRunValue")

if err := client.DeleteThenPoll(ctx, id, updateruns.DefaultDeleteOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `UpdateRunsClient.Get`

```go
ctx := context.TODO()
id := updateruns.NewUpdateRunID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fleetValue", "updateRunValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `UpdateRunsClient.ListByFleet`

```go
ctx := context.TODO()
id := updateruns.NewFleetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fleetValue")

// alternatively `client.ListByFleet(ctx, id)` can be used to do batched pagination
items, err := client.ListByFleetComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `UpdateRunsClient.Start`

```go
ctx := context.TODO()
id := updateruns.NewUpdateRunID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fleetValue", "updateRunValue")

if err := client.StartThenPoll(ctx, id, updateruns.DefaultStartOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `UpdateRunsClient.Stop`

```go
ctx := context.TODO()
id := updateruns.NewUpdateRunID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fleetValue", "updateRunValue")

if err := client.StopThenPoll(ctx, id, updateruns.DefaultStopOperationOptions()); err != nil {
	// handle the error
}
```
