
## `github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/capacityreservations` Documentation

The `capacityreservations` SDK allows for interaction with Azure Resource Manager `compute` (API Version `2022-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/capacityreservations"
```


### Client Initialization

```go
client := capacityreservations.NewCapacityReservationsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `CapacityReservationsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := capacityreservations.NewCapacityReservationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "capacityReservationGroupName", "capacityReservationName")

payload := capacityreservations.CapacityReservation{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CapacityReservationsClient.Delete`

```go
ctx := context.TODO()
id := capacityreservations.NewCapacityReservationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "capacityReservationGroupName", "capacityReservationName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CapacityReservationsClient.Get`

```go
ctx := context.TODO()
id := capacityreservations.NewCapacityReservationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "capacityReservationGroupName", "capacityReservationName")

read, err := client.Get(ctx, id, capacityreservations.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CapacityReservationsClient.Update`

```go
ctx := context.TODO()
id := capacityreservations.NewCapacityReservationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "capacityReservationGroupName", "capacityReservationName")

payload := capacityreservations.CapacityReservationUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
