
## `github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/privatelinkhubs` Documentation

The `privatelinkhubs` SDK allows for interaction with Azure Resource Manager `synapse` (API Version `2021-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/privatelinkhubs"
```


### Client Initialization

```go
client := privatelinkhubs.NewPrivateLinkHubsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `PrivateLinkHubsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := privatelinkhubs.NewPrivateLinkHubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateLinkHubName")

payload := privatelinkhubs.PrivateLinkHub{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PrivateLinkHubsClient.Delete`

```go
ctx := context.TODO()
id := privatelinkhubs.NewPrivateLinkHubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateLinkHubName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `PrivateLinkHubsClient.Get`

```go
ctx := context.TODO()
id := privatelinkhubs.NewPrivateLinkHubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateLinkHubName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PrivateLinkHubsClient.List`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `PrivateLinkHubsClient.ListByResourceGroup`

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


### Example Usage: `PrivateLinkHubsClient.PrivateEndpointConnectionsPrivateLinkHubGet`

```go
ctx := context.TODO()
id := privatelinkhubs.NewPrivateLinkHubPrivateEndpointConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateLinkHubName", "privateEndpointConnectionName")

read, err := client.PrivateEndpointConnectionsPrivateLinkHubGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PrivateLinkHubsClient.PrivateEndpointConnectionsPrivateLinkHubList`

```go
ctx := context.TODO()
id := privatelinkhubs.NewPrivateLinkHubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateLinkHubName")

// alternatively `client.PrivateEndpointConnectionsPrivateLinkHubList(ctx, id)` can be used to do batched pagination
items, err := client.PrivateEndpointConnectionsPrivateLinkHubListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `PrivateLinkHubsClient.Update`

```go
ctx := context.TODO()
id := privatelinkhubs.NewPrivateLinkHubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateLinkHubName")

payload := privatelinkhubs.PrivateLinkHubPatchInfo{
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
