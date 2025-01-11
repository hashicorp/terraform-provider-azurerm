
## `github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/clouds` Documentation

The `clouds` SDK allows for interaction with Azure Resource Manager `systemcentervirtualmachinemanager` (API Version `2023-10-07`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/clouds"
```


### Client Initialization

```go
client := clouds.NewCloudsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `CloudsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := clouds.NewCloudID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cloudName")

payload := clouds.Cloud{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CloudsClient.Delete`

```go
ctx := context.TODO()
id := clouds.NewCloudID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cloudName")

if err := client.DeleteThenPoll(ctx, id, clouds.DefaultDeleteOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `CloudsClient.Get`

```go
ctx := context.TODO()
id := clouds.NewCloudID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cloudName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CloudsClient.ListByResourceGroup`

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


### Example Usage: `CloudsClient.ListBySubscription`

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


### Example Usage: `CloudsClient.Update`

```go
ctx := context.TODO()
id := clouds.NewCloudID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cloudName")

payload := clouds.CloudTagsUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
