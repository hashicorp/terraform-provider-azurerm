
## `github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2023-10-15/fleetupdatestrategies` Documentation

The `fleetupdatestrategies` SDK allows for interaction with the Azure Resource Manager Service `containerservice` (API Version `2023-10-15`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2023-10-15/fleetupdatestrategies"
```


### Client Initialization

```go
client := fleetupdatestrategies.NewFleetUpdateStrategiesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `FleetUpdateStrategiesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := fleetupdatestrategies.NewUpdateStrategyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fleetValue", "updateStrategyValue")

payload := fleetupdatestrategies.FleetUpdateStrategy{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload, fleetupdatestrategies.DefaultCreateOrUpdateOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `FleetUpdateStrategiesClient.Delete`

```go
ctx := context.TODO()
id := fleetupdatestrategies.NewUpdateStrategyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fleetValue", "updateStrategyValue")

if err := client.DeleteThenPoll(ctx, id, fleetupdatestrategies.DefaultDeleteOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `FleetUpdateStrategiesClient.Get`

```go
ctx := context.TODO()
id := fleetupdatestrategies.NewUpdateStrategyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fleetValue", "updateStrategyValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FleetUpdateStrategiesClient.ListByFleet`

```go
ctx := context.TODO()
id := fleetupdatestrategies.NewFleetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fleetValue")

// alternatively `client.ListByFleet(ctx, id)` can be used to do batched pagination
items, err := client.ListByFleetComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
