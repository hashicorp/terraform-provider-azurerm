
## `github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2022-06-15/webtestsapis` Documentation

The `webtestsapis` SDK allows for interaction with Azure Resource Manager `applicationinsights` (API Version `2022-06-15`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2022-06-15/webtestsapis"
```


### Client Initialization

```go
client := webtestsapis.NewWebTestsAPIsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `WebTestsAPIsClient.WebTestsCreateOrUpdate`

```go
ctx := context.TODO()
id := webtestsapis.NewWebTestID("12345678-1234-9876-4563-123456789012", "example-resource-group", "webTestName")

payload := webtestsapis.WebTest{
	// ...
}


read, err := client.WebTestsCreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebTestsAPIsClient.WebTestsDelete`

```go
ctx := context.TODO()
id := webtestsapis.NewWebTestID("12345678-1234-9876-4563-123456789012", "example-resource-group", "webTestName")

read, err := client.WebTestsDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebTestsAPIsClient.WebTestsGet`

```go
ctx := context.TODO()
id := webtestsapis.NewWebTestID("12345678-1234-9876-4563-123456789012", "example-resource-group", "webTestName")

read, err := client.WebTestsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebTestsAPIsClient.WebTestsList`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.WebTestsList(ctx, id)` can be used to do batched pagination
items, err := client.WebTestsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebTestsAPIsClient.WebTestsListByComponent`

```go
ctx := context.TODO()
id := webtestsapis.NewComponentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "componentName")

// alternatively `client.WebTestsListByComponent(ctx, id)` can be used to do batched pagination
items, err := client.WebTestsListByComponentComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebTestsAPIsClient.WebTestsListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.WebTestsListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.WebTestsListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebTestsAPIsClient.WebTestsUpdateTags`

```go
ctx := context.TODO()
id := webtestsapis.NewWebTestID("12345678-1234-9876-4563-123456789012", "example-resource-group", "webTestName")

payload := webtestsapis.TagsResource{
	// ...
}


read, err := client.WebTestsUpdateTags(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
