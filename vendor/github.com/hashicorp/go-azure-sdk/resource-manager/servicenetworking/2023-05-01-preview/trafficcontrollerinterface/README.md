
## `github.com/hashicorp/go-azure-sdk/resource-manager/servicenetworking/2023-05-01-preview/trafficcontrollerinterface` Documentation

The `trafficcontrollerinterface` SDK allows for interaction with the Azure Resource Manager Service `servicenetworking` (API Version `2023-05-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/servicenetworking/2023-05-01-preview/trafficcontrollerinterface"
```


### Client Initialization

```go
client := trafficcontrollerinterface.NewTrafficControllerInterfaceClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `TrafficControllerInterfaceClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := trafficcontrollerinterface.NewTrafficControllerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "trafficControllerValue")

payload := trafficcontrollerinterface.TrafficController{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `TrafficControllerInterfaceClient.Delete`

```go
ctx := context.TODO()
id := trafficcontrollerinterface.NewTrafficControllerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "trafficControllerValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `TrafficControllerInterfaceClient.Get`

```go
ctx := context.TODO()
id := trafficcontrollerinterface.NewTrafficControllerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "trafficControllerValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TrafficControllerInterfaceClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `TrafficControllerInterfaceClient.ListBySubscription`

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


### Example Usage: `TrafficControllerInterfaceClient.Update`

```go
ctx := context.TODO()
id := trafficcontrollerinterface.NewTrafficControllerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "trafficControllerValue")

payload := trafficcontrollerinterface.TrafficControllerUpdate{
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
