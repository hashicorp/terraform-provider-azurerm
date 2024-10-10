
## `github.com/hashicorp/go-azure-sdk/resource-manager/advisor/2023-01-01/suppressions` Documentation

The `suppressions` SDK allows for interaction with Azure Resource Manager `advisor` (API Version `2023-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/advisor/2023-01-01/suppressions"
```


### Client Initialization

```go
client := suppressions.NewSuppressionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SuppressionsClient.Create`

```go
ctx := context.TODO()
id := suppressions.NewScopedSuppressionID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "recommendationId", "suppressionName")

payload := suppressions.SuppressionContract{
	// ...
}


read, err := client.Create(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SuppressionsClient.Delete`

```go
ctx := context.TODO()
id := suppressions.NewScopedSuppressionID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "recommendationId", "suppressionName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SuppressionsClient.Get`

```go
ctx := context.TODO()
id := suppressions.NewScopedSuppressionID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "recommendationId", "suppressionName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SuppressionsClient.List`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.List(ctx, id, suppressions.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, suppressions.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
