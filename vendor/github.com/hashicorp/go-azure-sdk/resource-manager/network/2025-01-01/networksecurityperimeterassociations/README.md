
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/networksecurityperimeterassociations` Documentation

The `networksecurityperimeterassociations` SDK allows for interaction with Azure Resource Manager `network` (API Version `2025-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/networksecurityperimeterassociations"
```


### Client Initialization

```go
client := networksecurityperimeterassociations.NewNetworkSecurityPerimeterAssociationsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `NetworkSecurityPerimeterAssociationsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := networksecurityperimeterassociations.NewResourceAssociationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkSecurityPerimeterName", "resourceAssociationName")

payload := networksecurityperimeterassociations.NspAssociation{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `NetworkSecurityPerimeterAssociationsClient.Delete`

```go
ctx := context.TODO()
id := networksecurityperimeterassociations.NewResourceAssociationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkSecurityPerimeterName", "resourceAssociationName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `NetworkSecurityPerimeterAssociationsClient.Get`

```go
ctx := context.TODO()
id := networksecurityperimeterassociations.NewResourceAssociationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkSecurityPerimeterName", "resourceAssociationName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NetworkSecurityPerimeterAssociationsClient.List`

```go
ctx := context.TODO()
id := networksecurityperimeterassociations.NewNetworkSecurityPerimeterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkSecurityPerimeterName")

// alternatively `client.List(ctx, id, networksecurityperimeterassociations.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, networksecurityperimeterassociations.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `NetworkSecurityPerimeterAssociationsClient.Reconcile`

```go
ctx := context.TODO()
id := networksecurityperimeterassociations.NewResourceAssociationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkSecurityPerimeterName", "resourceAssociationName")
var payload interface{}

read, err := client.Reconcile(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
