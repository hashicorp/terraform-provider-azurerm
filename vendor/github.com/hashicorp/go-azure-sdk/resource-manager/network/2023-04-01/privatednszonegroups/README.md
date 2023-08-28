
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/privatednszonegroups` Documentation

The `privatednszonegroups` SDK allows for interaction with the Azure Resource Manager Service `network` (API Version `2023-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/privatednszonegroups"
```


### Client Initialization

```go
client := privatednszonegroups.NewPrivateDnsZoneGroupsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `PrivateDnsZoneGroupsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := privatednszonegroups.NewPrivateDnsZoneGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateEndpointValue", "privateDnsZoneGroupValue")

payload := privatednszonegroups.PrivateDnsZoneGroup{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `PrivateDnsZoneGroupsClient.Delete`

```go
ctx := context.TODO()
id := privatednszonegroups.NewPrivateDnsZoneGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateEndpointValue", "privateDnsZoneGroupValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `PrivateDnsZoneGroupsClient.Get`

```go
ctx := context.TODO()
id := privatednszonegroups.NewPrivateDnsZoneGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateEndpointValue", "privateDnsZoneGroupValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PrivateDnsZoneGroupsClient.List`

```go
ctx := context.TODO()
id := privatednszonegroups.NewPrivateEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateEndpointValue")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
