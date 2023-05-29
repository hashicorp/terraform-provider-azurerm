
## `github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2022-09-02-preview/fleetmembers` Documentation

The `fleetmembers` SDK allows for interaction with the Azure Resource Manager Service `containerservice` (API Version `2022-09-02-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2022-09-02-preview/fleetmembers"
```


### Client Initialization

```go
client := fleetmembers.NewFleetMembersClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `FleetMembersClient.Create`

```go
ctx := context.TODO()
id := fleetmembers.NewMemberID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fleetValue", "memberValue")

payload := fleetmembers.FleetMember{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload, fleetmembers.DefaultCreateOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `FleetMembersClient.Delete`

```go
ctx := context.TODO()
id := fleetmembers.NewMemberID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fleetValue", "memberValue")

if err := client.DeleteThenPoll(ctx, id, fleetmembers.DefaultDeleteOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `FleetMembersClient.Get`

```go
ctx := context.TODO()
id := fleetmembers.NewMemberID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fleetValue", "memberValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FleetMembersClient.ListByFleet`

```go
ctx := context.TODO()
id := fleetmembers.NewFleetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fleetValue")

// alternatively `client.ListByFleet(ctx, id)` can be used to do batched pagination
items, err := client.ListByFleetComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
