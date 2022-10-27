
## `github.com/hashicorp/go-azure-sdk/resource-manager/customproviders/2018-09-01-preview/customresourceprovider` Documentation

The `customresourceprovider` SDK allows for interaction with the Azure Resource Manager Service `customproviders` (API Version `2018-09-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/customproviders/2018-09-01-preview/customresourceprovider"
```


### Client Initialization

```go
client := customresourceprovider.NewCustomResourceProviderClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `CustomResourceProviderClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := customresourceprovider.NewResourceProviderID("12345678-1234-9876-4563-123456789012", "example-resource-group", "resourceProviderValue")

payload := customresourceprovider.CustomRPManifest{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CustomResourceProviderClient.Delete`

```go
ctx := context.TODO()
id := customresourceprovider.NewResourceProviderID("12345678-1234-9876-4563-123456789012", "example-resource-group", "resourceProviderValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CustomResourceProviderClient.Get`

```go
ctx := context.TODO()
id := customresourceprovider.NewResourceProviderID("12345678-1234-9876-4563-123456789012", "example-resource-group", "resourceProviderValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CustomResourceProviderClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := customresourceprovider.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `CustomResourceProviderClient.ListBySubscription`

```go
ctx := context.TODO()
id := customresourceprovider.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `CustomResourceProviderClient.Update`

```go
ctx := context.TODO()
id := customresourceprovider.NewResourceProviderID("12345678-1234-9876-4563-123456789012", "example-resource-group", "resourceProviderValue")

payload := customresourceprovider.ResourceProvidersUpdate{
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
