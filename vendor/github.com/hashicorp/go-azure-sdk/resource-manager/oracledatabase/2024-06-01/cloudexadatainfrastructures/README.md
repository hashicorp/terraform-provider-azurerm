
## `github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/cloudexadatainfrastructures` Documentation

The `cloudexadatainfrastructures` SDK allows for interaction with Azure Resource Manager `oracledatabase` (API Version `2024-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/cloudexadatainfrastructures"
```


### Client Initialization

```go
client := cloudexadatainfrastructures.NewCloudExadataInfrastructuresClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `CloudExadataInfrastructuresClient.AddStorageCapacity`

```go
ctx := context.TODO()
id := cloudexadatainfrastructures.NewCloudExadataInfrastructureID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cloudExadataInfrastructureName")

if err := client.AddStorageCapacityThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CloudExadataInfrastructuresClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := cloudexadatainfrastructures.NewCloudExadataInfrastructureID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cloudExadataInfrastructureName")

payload := cloudexadatainfrastructures.CloudExadataInfrastructure{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CloudExadataInfrastructuresClient.Delete`

```go
ctx := context.TODO()
id := cloudexadatainfrastructures.NewCloudExadataInfrastructureID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cloudExadataInfrastructureName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CloudExadataInfrastructuresClient.Get`

```go
ctx := context.TODO()
id := cloudexadatainfrastructures.NewCloudExadataInfrastructureID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cloudExadataInfrastructureName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CloudExadataInfrastructuresClient.ListByResourceGroup`

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


### Example Usage: `CloudExadataInfrastructuresClient.ListBySubscription`

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


### Example Usage: `CloudExadataInfrastructuresClient.Update`

```go
ctx := context.TODO()
id := cloudexadatainfrastructures.NewCloudExadataInfrastructureID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cloudExadataInfrastructureName")

payload := cloudexadatainfrastructures.CloudExadataInfrastructureUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
