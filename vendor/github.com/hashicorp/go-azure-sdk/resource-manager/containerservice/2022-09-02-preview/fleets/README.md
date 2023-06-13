
## `github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2022-09-02-preview/fleets` Documentation

The `fleets` SDK allows for interaction with the Azure Resource Manager Service `containerservice` (API Version `2022-09-02-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2022-09-02-preview/fleets"
```


### Client Initialization

```go
client := fleets.NewFleetsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `FleetsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := fleets.NewFleetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fleetValue")

payload := fleets.Fleet{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload, fleets.DefaultCreateOrUpdateOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `FleetsClient.Delete`

```go
ctx := context.TODO()
id := fleets.NewFleetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fleetValue")

if err := client.DeleteThenPoll(ctx, id, fleets.DefaultDeleteOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `FleetsClient.Get`

```go
ctx := context.TODO()
id := fleets.NewFleetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fleetValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FleetsClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := fleets.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `FleetsClient.ListBySubscription`

```go
ctx := context.TODO()
id := fleets.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `FleetsClient.ListCredentials`

```go
ctx := context.TODO()
id := fleets.NewFleetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fleetValue")

read, err := client.ListCredentials(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FleetsClient.Update`

```go
ctx := context.TODO()
id := fleets.NewFleetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fleetValue")

payload := fleets.FleetPatch{
	// ...
}


read, err := client.Update(ctx, id, payload, fleets.DefaultUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
