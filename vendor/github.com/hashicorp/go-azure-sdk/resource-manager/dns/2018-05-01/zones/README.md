
## `github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01/zones` Documentation

The `zones` SDK allows for interaction with the Azure Resource Manager Service `dns` (API Version `2018-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01/zones"
```


### Client Initialization

```go
client := zones.NewZonesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ZonesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := zones.NewDnsZoneID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dnsZoneValue")

payload := zones.Zone{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload, zones.DefaultCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ZonesClient.Delete`

```go
ctx := context.TODO()
id := zones.NewDnsZoneID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dnsZoneValue")

if err := client.DeleteThenPoll(ctx, id, zones.DefaultDeleteOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `ZonesClient.Get`

```go
ctx := context.TODO()
id := zones.NewDnsZoneID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dnsZoneValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ZonesClient.List`

```go
ctx := context.TODO()
id := zones.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.List(ctx, id, zones.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, zones.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ZonesClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := zones.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id, zones.DefaultListByResourceGroupOperationOptions())` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id, zones.DefaultListByResourceGroupOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ZonesClient.Update`

```go
ctx := context.TODO()
id := zones.NewDnsZoneID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dnsZoneValue")

payload := zones.ZoneUpdate{
	// ...
}


read, err := client.Update(ctx, id, payload, zones.DefaultUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
