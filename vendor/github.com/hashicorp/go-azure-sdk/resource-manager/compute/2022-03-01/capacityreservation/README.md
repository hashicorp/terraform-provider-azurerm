
## `github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/capacityreservation` Documentation

The `capacityreservation` SDK allows for interaction with Azure Resource Manager `compute` (API Version `2022-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/capacityreservation"
```


### Client Initialization

```go
client := capacityreservation.NewCapacityReservationClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `CapacityReservationClient.ListByCapacityReservationGroup`

```go
ctx := context.TODO()
id := capacityreservation.NewCapacityReservationGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "capacityReservationGroupName")

// alternatively `client.ListByCapacityReservationGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByCapacityReservationGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
