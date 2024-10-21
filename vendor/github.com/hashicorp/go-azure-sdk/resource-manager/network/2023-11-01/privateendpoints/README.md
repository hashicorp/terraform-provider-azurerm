
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/privateendpoints` Documentation

The `privateendpoints` SDK allows for interaction with Azure Resource Manager `network` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/privateendpoints"
```


### Client Initialization

```go
client := privateendpoints.NewPrivateEndpointsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `PrivateEndpointsClient.AvailablePrivateEndpointTypesList`

```go
ctx := context.TODO()
id := privateendpoints.NewLocationID("12345678-1234-9876-4563-123456789012", "locationName")

// alternatively `client.AvailablePrivateEndpointTypesList(ctx, id)` can be used to do batched pagination
items, err := client.AvailablePrivateEndpointTypesListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `PrivateEndpointsClient.AvailablePrivateEndpointTypesListByResourceGroup`

```go
ctx := context.TODO()
id := privateendpoints.NewProviderLocationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "locationName")

// alternatively `client.AvailablePrivateEndpointTypesListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.AvailablePrivateEndpointTypesListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `PrivateEndpointsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := privateendpoints.NewPrivateEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateEndpointName")

payload := privateendpoints.PrivateEndpoint{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `PrivateEndpointsClient.Delete`

```go
ctx := context.TODO()
id := privateendpoints.NewPrivateEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateEndpointName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `PrivateEndpointsClient.Get`

```go
ctx := context.TODO()
id := privateendpoints.NewPrivateEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateEndpointName")

read, err := client.Get(ctx, id, privateendpoints.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PrivateEndpointsClient.List`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `PrivateEndpointsClient.ListBySubscription`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
