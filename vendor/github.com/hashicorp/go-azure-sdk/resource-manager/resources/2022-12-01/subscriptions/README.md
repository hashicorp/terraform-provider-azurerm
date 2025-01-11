
## `github.com/hashicorp/go-azure-sdk/resource-manager/resources/2022-12-01/subscriptions` Documentation

The `subscriptions` SDK allows for interaction with Azure Resource Manager `resources` (API Version `2022-12-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/resources/2022-12-01/subscriptions"
```


### Client Initialization

```go
client := subscriptions.NewSubscriptionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SubscriptionsClient.CheckZonePeers`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

payload := subscriptions.CheckZonePeersRequest{
	// ...
}


read, err := client.CheckZonePeers(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SubscriptionsClient.Get`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SubscriptionsClient.List`

```go
ctx := context.TODO()


// alternatively `client.List(ctx)` can be used to do batched pagination
items, err := client.ListComplete(ctx)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `SubscriptionsClient.ListLocations`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

read, err := client.ListLocations(ctx, id, subscriptions.DefaultListLocationsOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
