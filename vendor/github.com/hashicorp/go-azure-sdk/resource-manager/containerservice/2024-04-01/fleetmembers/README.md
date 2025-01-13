
## `github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2024-04-01/fleetmembers` Documentation

The `fleetmembers` SDK allows for interaction with Azure Resource Manager `containerservice` (API Version `2024-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2024-04-01/fleetmembers"
```


### Client Initialization

```go
client := fleetmembers.NewFleetMembersClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `FleetMembersClient.Create`

```go
ctx := context.TODO()
id := fleetmembers.NewMemberID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fleetName", "memberName")

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
id := fleetmembers.NewMemberID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fleetName", "memberName")

if err := client.DeleteThenPoll(ctx, id, fleetmembers.DefaultDeleteOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `FleetMembersClient.Get`

```go
ctx := context.TODO()
id := fleetmembers.NewMemberID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fleetName", "memberName")

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
id := fleetmembers.NewFleetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fleetName")

// alternatively `client.ListByFleet(ctx, id)` can be used to do batched pagination
items, err := client.ListByFleetComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `FleetMembersClient.Update`

```go
ctx := context.TODO()
id := fleetmembers.NewMemberID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fleetName", "memberName")

payload := fleetmembers.FleetMemberUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload, fleetmembers.DefaultUpdateOperationOptions()); err != nil {
	// handle the error
}
```
