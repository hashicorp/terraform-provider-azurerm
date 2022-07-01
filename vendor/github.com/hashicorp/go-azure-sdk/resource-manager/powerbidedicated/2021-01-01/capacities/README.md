
## `github.com/hashicorp/go-azure-sdk/resource-manager/powerbidedicated/2021-01-01/capacities` Documentation

The `capacities` SDK allows for interaction with the Azure Resource Manager Service `powerbidedicated` (API Version `2021-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/powerbidedicated/2021-01-01/capacities"
```


### Client Initialization

```go
client := capacities.NewCapacitiesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `CapacitiesClient.CheckNameAvailability`

```go
ctx := context.TODO()
id := capacities.NewLocationID("12345678-1234-9876-4563-123456789012", "locationValue")

payload := capacities.CheckCapacityNameAvailabilityParameters{
	// ...
}


read, err := client.CheckNameAvailability(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CapacitiesClient.Create`

```go
ctx := context.TODO()
id := capacities.NewCapacitiesID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dedicatedCapacityValue")

payload := capacities.DedicatedCapacity{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CapacitiesClient.Delete`

```go
ctx := context.TODO()
id := capacities.NewCapacitiesID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dedicatedCapacityValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CapacitiesClient.GetDetails`

```go
ctx := context.TODO()
id := capacities.NewCapacitiesID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dedicatedCapacityValue")

read, err := client.GetDetails(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CapacitiesClient.List`

```go
ctx := context.TODO()
id := capacities.NewSubscriptionID()

read, err := client.List(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CapacitiesClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := capacities.NewResourceGroupID()

read, err := client.ListByResourceGroup(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CapacitiesClient.ListSkusForCapacity`

```go
ctx := context.TODO()
id := capacities.NewCapacitiesID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dedicatedCapacityValue")

read, err := client.ListSkusForCapacity(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CapacitiesClient.Resume`

```go
ctx := context.TODO()
id := capacities.NewCapacitiesID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dedicatedCapacityValue")

if err := client.ResumeThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CapacitiesClient.Suspend`

```go
ctx := context.TODO()
id := capacities.NewCapacitiesID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dedicatedCapacityValue")

if err := client.SuspendThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CapacitiesClient.Update`

```go
ctx := context.TODO()
id := capacities.NewCapacitiesID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dedicatedCapacityValue")

payload := capacities.DedicatedCapacityUpdateParameters{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
