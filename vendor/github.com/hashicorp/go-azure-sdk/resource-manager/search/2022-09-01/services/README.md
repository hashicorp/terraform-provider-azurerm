
## `github.com/hashicorp/go-azure-sdk/resource-manager/search/2022-09-01/services` Documentation

The `services` SDK allows for interaction with the Azure Resource Manager Service `search` (API Version `2022-09-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/search/2022-09-01/services"
```


### Client Initialization

```go
client := services.NewServicesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ServicesClient.CheckNameAvailability`

```go
ctx := context.TODO()
id := services.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

payload := services.CheckNameAvailabilityInput{
	// ...
}


read, err := client.CheckNameAvailability(ctx, id, payload, services.DefaultCheckNameAvailabilityOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ServicesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := services.NewSearchServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "searchServiceValue")

payload := services.SearchService{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload, services.DefaultCreateOrUpdateOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `ServicesClient.Delete`

```go
ctx := context.TODO()
id := services.NewSearchServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "searchServiceValue")

read, err := client.Delete(ctx, id, services.DefaultDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ServicesClient.Get`

```go
ctx := context.TODO()
id := services.NewSearchServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "searchServiceValue")

read, err := client.Get(ctx, id, services.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ServicesClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := services.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id, services.DefaultListByResourceGroupOperationOptions())` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id, services.DefaultListByResourceGroupOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ServicesClient.ListBySubscription`

```go
ctx := context.TODO()
id := services.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id, services.DefaultListBySubscriptionOperationOptions())` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id, services.DefaultListBySubscriptionOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ServicesClient.Update`

```go
ctx := context.TODO()
id := services.NewSearchServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "searchServiceValue")

payload := services.SearchServiceUpdate{
	// ...
}


read, err := client.Update(ctx, id, payload, services.DefaultUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
