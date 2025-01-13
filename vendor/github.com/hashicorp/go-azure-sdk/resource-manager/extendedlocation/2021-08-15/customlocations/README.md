
## `github.com/hashicorp/go-azure-sdk/resource-manager/extendedlocation/2021-08-15/customlocations` Documentation

The `customlocations` SDK allows for interaction with Azure Resource Manager `extendedlocation` (API Version `2021-08-15`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/extendedlocation/2021-08-15/customlocations"
```


### Client Initialization

```go
client := customlocations.NewCustomLocationsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `CustomLocationsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := customlocations.NewCustomLocationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "customLocationName")

payload := customlocations.CustomLocation{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CustomLocationsClient.Delete`

```go
ctx := context.TODO()
id := customlocations.NewCustomLocationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "customLocationName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CustomLocationsClient.Get`

```go
ctx := context.TODO()
id := customlocations.NewCustomLocationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "customLocationName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CustomLocationsClient.ListByResourceGroup`

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


### Example Usage: `CustomLocationsClient.ListBySubscription`

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


### Example Usage: `CustomLocationsClient.ListEnabledResourceTypes`

```go
ctx := context.TODO()
id := customlocations.NewCustomLocationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "customLocationName")

// alternatively `client.ListEnabledResourceTypes(ctx, id)` can be used to do batched pagination
items, err := client.ListEnabledResourceTypesComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `CustomLocationsClient.Update`

```go
ctx := context.TODO()
id := customlocations.NewCustomLocationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "customLocationName")

payload := customlocations.PatchableCustomLocations{
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
