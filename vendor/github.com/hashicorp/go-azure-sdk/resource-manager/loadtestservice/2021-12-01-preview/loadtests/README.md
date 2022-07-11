
## `github.com/hashicorp/go-azure-sdk/resource-manager/loadtestservice/2021-12-01-preview/loadtests` Documentation

The `loadtests` SDK allows for interaction with the Azure Resource Manager Service `loadtestservice` (API Version `2021-12-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/loadtestservice/2021-12-01-preview/loadtests"
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


read, err := client.CreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
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
id := loadtests.NewResourceGroupID()

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
id := loadtests.NewSubscriptionID()

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


read, err := client.Update(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
