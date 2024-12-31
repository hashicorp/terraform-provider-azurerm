
## `github.com/hashicorp/go-azure-sdk/resource-manager/deviceregistry/2024-11-01/assetendpointprofiles` Documentation

The `assetendpointprofiles` SDK allows for interaction with Azure Resource Manager `deviceregistry` (API Version `2024-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/deviceregistry/2024-11-01/assetendpointprofiles"
```


### Client Initialization

```go
client := assetendpointprofiles.NewAssetEndpointProfilesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AssetEndpointProfilesClient.CreateOrReplace`

```go
ctx := context.TODO()
id := assetendpointprofiles.NewAssetEndpointProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "assetEndpointProfileName")

payload := assetendpointprofiles.AssetEndpointProfile{
	// ...
}


if err := client.CreateOrReplaceThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AssetEndpointProfilesClient.Delete`

```go
ctx := context.TODO()
id := assetendpointprofiles.NewAssetEndpointProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "assetEndpointProfileName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AssetEndpointProfilesClient.Get`

```go
ctx := context.TODO()
id := assetendpointprofiles.NewAssetEndpointProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "assetEndpointProfileName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AssetEndpointProfilesClient.ListByResourceGroup`

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


### Example Usage: `AssetEndpointProfilesClient.ListBySubscription`

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


### Example Usage: `AssetEndpointProfilesClient.Update`

```go
ctx := context.TODO()
id := assetendpointprofiles.NewAssetEndpointProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "assetEndpointProfileName")

payload := assetendpointprofiles.AssetEndpointProfileUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
