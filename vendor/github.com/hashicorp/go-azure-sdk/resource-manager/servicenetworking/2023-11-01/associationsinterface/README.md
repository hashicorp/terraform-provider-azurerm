
## `github.com/hashicorp/go-azure-sdk/resource-manager/servicenetworking/2023-11-01/associationsinterface` Documentation

The `associationsinterface` SDK allows for interaction with Azure Resource Manager `servicenetworking` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/servicenetworking/2023-11-01/associationsinterface"
```


### Client Initialization

```go
client := associationsinterface.NewAssociationsInterfaceClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AssociationsInterfaceClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := associationsinterface.NewAssociationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "trafficControllerName", "associationName")

payload := associationsinterface.Association{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AssociationsInterfaceClient.Delete`

```go
ctx := context.TODO()
id := associationsinterface.NewAssociationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "trafficControllerName", "associationName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AssociationsInterfaceClient.Get`

```go
ctx := context.TODO()
id := associationsinterface.NewAssociationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "trafficControllerName", "associationName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AssociationsInterfaceClient.ListByTrafficController`

```go
ctx := context.TODO()
id := associationsinterface.NewTrafficControllerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "trafficControllerName")

// alternatively `client.ListByTrafficController(ctx, id)` can be used to do batched pagination
items, err := client.ListByTrafficControllerComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AssociationsInterfaceClient.Update`

```go
ctx := context.TODO()
id := associationsinterface.NewAssociationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "trafficControllerName", "associationName")

payload := associationsinterface.AssociationUpdate{
	// ...
}


read, err := client.Update(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
