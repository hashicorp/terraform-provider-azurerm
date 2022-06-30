
## `github.com/hashicorp/go-azure-sdk/resource-manager/mixedreality/2021-01-01/resource` Documentation

The `resource` SDK allows for interaction with the Azure Resource Manager Service `mixedreality` (API Version `2021-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/mixedreality/2021-01-01/resource"
```


### Client Initialization

```go
client := resource.NewResourceClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ResourceClient.RemoteRenderingAccountsCreate`

```go
ctx := context.TODO()
id := resource.NewRemoteRenderingAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue")

payload := resource.RemoteRenderingAccount{
	// ...
}


read, err := client.RemoteRenderingAccountsCreate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ResourceClient.RemoteRenderingAccountsDelete`

```go
ctx := context.TODO()
id := resource.NewRemoteRenderingAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue")

read, err := client.RemoteRenderingAccountsDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ResourceClient.RemoteRenderingAccountsGet`

```go
ctx := context.TODO()
id := resource.NewRemoteRenderingAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue")

read, err := client.RemoteRenderingAccountsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ResourceClient.RemoteRenderingAccountsListByResourceGroup`

```go
ctx := context.TODO()
id := resource.NewResourceGroupID()

// alternatively `client.RemoteRenderingAccountsListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.RemoteRenderingAccountsListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ResourceClient.RemoteRenderingAccountsListBySubscription`

```go
ctx := context.TODO()
id := resource.NewSubscriptionID()

// alternatively `client.RemoteRenderingAccountsListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.RemoteRenderingAccountsListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ResourceClient.RemoteRenderingAccountsUpdate`

```go
ctx := context.TODO()
id := resource.NewRemoteRenderingAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue")

payload := resource.RemoteRenderingAccount{
	// ...
}


read, err := client.RemoteRenderingAccountsUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ResourceClient.SpatialAnchorsAccountsCreate`

```go
ctx := context.TODO()
id := resource.NewSpatialAnchorsAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue")

payload := resource.SpatialAnchorsAccount{
	// ...
}


read, err := client.SpatialAnchorsAccountsCreate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ResourceClient.SpatialAnchorsAccountsDelete`

```go
ctx := context.TODO()
id := resource.NewSpatialAnchorsAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue")

read, err := client.SpatialAnchorsAccountsDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ResourceClient.SpatialAnchorsAccountsGet`

```go
ctx := context.TODO()
id := resource.NewSpatialAnchorsAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue")

read, err := client.SpatialAnchorsAccountsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ResourceClient.SpatialAnchorsAccountsListByResourceGroup`

```go
ctx := context.TODO()
id := resource.NewResourceGroupID()

// alternatively `client.SpatialAnchorsAccountsListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.SpatialAnchorsAccountsListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ResourceClient.SpatialAnchorsAccountsListBySubscription`

```go
ctx := context.TODO()
id := resource.NewSubscriptionID()

// alternatively `client.SpatialAnchorsAccountsListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.SpatialAnchorsAccountsListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ResourceClient.SpatialAnchorsAccountsUpdate`

```go
ctx := context.TODO()
id := resource.NewSpatialAnchorsAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue")

payload := resource.SpatialAnchorsAccount{
	// ...
}


read, err := client.SpatialAnchorsAccountsUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
