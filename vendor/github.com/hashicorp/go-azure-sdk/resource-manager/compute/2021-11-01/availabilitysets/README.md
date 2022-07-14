
## `github.com/hashicorp/go-azure-sdk/resource-manager/compute/2021-11-01/availabilitysets` Documentation

The `availabilitysets` SDK allows for interaction with the Azure Resource Manager Service `compute` (API Version `2021-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/compute/2021-11-01/availabilitysets"
```


### Client Initialization

```go
client := availabilitysets.NewAvailabilitySetsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AvailabilitySetsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := availabilitysets.NewAvailabilitySetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "availabilitySetValue")

payload := availabilitysets.AvailabilitySet{
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


### Example Usage: `AvailabilitySetsClient.Delete`

```go
ctx := context.TODO()
id := availabilitysets.NewAvailabilitySetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "availabilitySetValue")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AvailabilitySetsClient.Get`

```go
ctx := context.TODO()
id := availabilitysets.NewAvailabilitySetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "availabilitySetValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AvailabilitySetsClient.List`

```go
ctx := context.TODO()
id := availabilitysets.NewResourceGroupID()

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AvailabilitySetsClient.ListAvailableSizes`

```go
ctx := context.TODO()
id := availabilitysets.NewAvailabilitySetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "availabilitySetValue")

read, err := client.ListAvailableSizes(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AvailabilitySetsClient.ListBySubscription`

```go
ctx := context.TODO()
id := availabilitysets.NewSubscriptionID()

// alternatively `client.ListBySubscription(ctx, id, availabilitysets.DefaultListBySubscriptionOperationOptions())` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id, availabilitysets.DefaultListBySubscriptionOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AvailabilitySetsClient.Update`

```go
ctx := context.TODO()
id := availabilitysets.NewAvailabilitySetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "availabilitySetValue")

payload := availabilitysets.AvailabilitySetUpdate{
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
