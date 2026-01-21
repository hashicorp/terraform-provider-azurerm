
## `github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2025-10-15/fleets` Documentation

The `fleets` SDK allows for interaction with Azure Resource Manager `cosmosdb` (API Version `2025-10-15`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2025-10-15/fleets"
```


### Client Initialization

```go
client := fleets.NewFleetsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `FleetsClient.FleetCreate`

```go
ctx := context.TODO()
id := fleets.NewFleetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fleetName")

payload := fleets.FleetResource{
	// ...
}


read, err := client.FleetCreate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FleetsClient.FleetDelete`

```go
ctx := context.TODO()
id := fleets.NewFleetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fleetName")

if err := client.FleetDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `FleetsClient.FleetGet`

```go
ctx := context.TODO()
id := fleets.NewFleetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fleetName")

read, err := client.FleetGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FleetsClient.FleetList`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.FleetList(ctx, id)` can be used to do batched pagination
items, err := client.FleetListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `FleetsClient.FleetListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.FleetListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.FleetListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `FleetsClient.FleetUpdate`

```go
ctx := context.TODO()
id := fleets.NewFleetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fleetName")

payload := fleets.FleetResourceUpdate{
	// ...
}


read, err := client.FleetUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FleetsClient.FleetspaceAccountCreate`

```go
ctx := context.TODO()
id := fleets.NewFleetspaceAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fleetName", "fleetspaceName", "fleetspaceAccountName")

payload := fleets.FleetspaceAccountResource{
	// ...
}


if err := client.FleetspaceAccountCreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `FleetsClient.FleetspaceAccountDelete`

```go
ctx := context.TODO()
id := fleets.NewFleetspaceAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fleetName", "fleetspaceName", "fleetspaceAccountName")

if err := client.FleetspaceAccountDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `FleetsClient.FleetspaceAccountGet`

```go
ctx := context.TODO()
id := fleets.NewFleetspaceAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fleetName", "fleetspaceName", "fleetspaceAccountName")

read, err := client.FleetspaceAccountGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FleetsClient.FleetspaceAccountList`

```go
ctx := context.TODO()
id := fleets.NewFleetspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fleetName", "fleetspaceName")

// alternatively `client.FleetspaceAccountList(ctx, id)` can be used to do batched pagination
items, err := client.FleetspaceAccountListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `FleetsClient.FleetspaceCreate`

```go
ctx := context.TODO()
id := fleets.NewFleetspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fleetName", "fleetspaceName")

payload := fleets.FleetspaceResource{
	// ...
}


if err := client.FleetspaceCreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `FleetsClient.FleetspaceDelete`

```go
ctx := context.TODO()
id := fleets.NewFleetspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fleetName", "fleetspaceName")

if err := client.FleetspaceDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `FleetsClient.FleetspaceGet`

```go
ctx := context.TODO()
id := fleets.NewFleetspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fleetName", "fleetspaceName")

read, err := client.FleetspaceGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FleetsClient.FleetspaceList`

```go
ctx := context.TODO()
id := fleets.NewFleetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fleetName")

// alternatively `client.FleetspaceList(ctx, id)` can be used to do batched pagination
items, err := client.FleetspaceListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `FleetsClient.FleetspaceUpdate`

```go
ctx := context.TODO()
id := fleets.NewFleetspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fleetName", "fleetspaceName")

payload := fleets.FleetspaceUpdate{
	// ...
}


if err := client.FleetspaceUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
