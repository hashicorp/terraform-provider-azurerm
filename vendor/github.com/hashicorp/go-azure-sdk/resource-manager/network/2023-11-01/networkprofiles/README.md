
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/networkprofiles` Documentation

The `networkprofiles` SDK allows for interaction with Azure Resource Manager `network` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/networkprofiles"
```


### Client Initialization

```go
client := networkprofiles.NewNetworkProfilesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `NetworkProfilesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := networkprofiles.NewNetworkProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkProfileName")

payload := networkprofiles.NetworkProfile{
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


### Example Usage: `NetworkProfilesClient.Delete`

```go
ctx := context.TODO()
id := networkprofiles.NewNetworkProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkProfileName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `NetworkProfilesClient.Get`

```go
ctx := context.TODO()
id := networkprofiles.NewNetworkProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkProfileName")

read, err := client.Get(ctx, id, networkprofiles.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NetworkProfilesClient.List`

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


### Example Usage: `NetworkProfilesClient.ListAll`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListAll(ctx, id)` can be used to do batched pagination
items, err := client.ListAllComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `NetworkProfilesClient.UpdateTags`

```go
ctx := context.TODO()
id := networkprofiles.NewNetworkProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkProfileName")

payload := networkprofiles.TagsObject{
	// ...
}


read, err := client.UpdateTags(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
