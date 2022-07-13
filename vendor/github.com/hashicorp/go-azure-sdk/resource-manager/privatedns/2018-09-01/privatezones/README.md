
## `github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2018-09-01/privatezones` Documentation

The `privatezones` SDK allows for interaction with the Azure Resource Manager Service `privatedns` (API Version `2018-09-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2018-09-01/privatezones"
```


### Client Initialization

```go
client := privatezones.NewPrivateZonesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `PrivateZonesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := privatezones.NewPrivateDnsZoneID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateZoneValue")

payload := privatezones.PrivateZone{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload, privatezones.DefaultCreateOrUpdateOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `PrivateZonesClient.Delete`

```go
ctx := context.TODO()
id := privatezones.NewPrivateDnsZoneID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateZoneValue")

if err := client.DeleteThenPoll(ctx, id, privatezones.DefaultDeleteOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `PrivateZonesClient.Get`

```go
ctx := context.TODO()
id := privatezones.NewPrivateDnsZoneID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateZoneValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PrivateZonesClient.List`

```go
ctx := context.TODO()
id := privatezones.NewSubscriptionID()

// alternatively `client.List(ctx, id, privatezones.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, privatezones.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `PrivateZonesClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := privatezones.NewResourceGroupID()

// alternatively `client.ListByResourceGroup(ctx, id, privatezones.DefaultListByResourceGroupOperationOptions())` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id, privatezones.DefaultListByResourceGroupOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `PrivateZonesClient.Update`

```go
ctx := context.TODO()
id := privatezones.NewPrivateDnsZoneID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateZoneValue")

payload := privatezones.PrivateZone{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload, privatezones.DefaultUpdateOperationOptions()); err != nil {
	// handle the error
}
```
