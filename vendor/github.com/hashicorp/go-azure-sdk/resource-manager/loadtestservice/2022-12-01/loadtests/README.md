
## `github.com/hashicorp/go-azure-sdk/resource-manager/loadtestservice/2022-12-01/loadtests` Documentation

The `loadtests` SDK allows for interaction with the Azure Resource Manager Service `loadtestservice` (API Version `2022-12-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/loadtestservice/2022-12-01/loadtests"
```


### Client Initialization

```go
client := loadtests.NewLoadTestsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `LoadTestsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := loadtests.NewLoadTestID("12345678-1234-9876-4563-123456789012", "example-resource-group", "loadTestValue")

payload := loadtests.LoadTestResource{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `LoadTestsClient.Delete`

```go
ctx := context.TODO()
id := loadtests.NewLoadTestID("12345678-1234-9876-4563-123456789012", "example-resource-group", "loadTestValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `LoadTestsClient.Get`

```go
ctx := context.TODO()
id := loadtests.NewLoadTestID("12345678-1234-9876-4563-123456789012", "example-resource-group", "loadTestValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `LoadTestsClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := loadtests.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `LoadTestsClient.ListBySubscription`

```go
ctx := context.TODO()
id := loadtests.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `LoadTestsClient.Update`

```go
ctx := context.TODO()
id := loadtests.NewLoadTestID("12345678-1234-9876-4563-123456789012", "example-resource-group", "loadTestValue")

payload := loadtests.LoadTestResourcePatchRequestBody{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
