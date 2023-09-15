
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/subscription` Documentation

The `subscription` SDK allows for interaction with the Azure Resource Manager Service `apimanagement` (API Version `2021-08-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/subscription"
```


### Client Initialization

```go
client := subscription.NewSubscriptionClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SubscriptionClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := subscription.NewSubscriptions2ID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "subscriptionValue")

payload := subscription.SubscriptionCreateParameters{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload, subscription.DefaultCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SubscriptionClient.Delete`

```go
ctx := context.TODO()
id := subscription.NewSubscriptions2ID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "subscriptionValue")

read, err := client.Delete(ctx, id, subscription.DefaultDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SubscriptionClient.Get`

```go
ctx := context.TODO()
id := subscription.NewSubscriptions2ID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "subscriptionValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SubscriptionClient.GetEntityTag`

```go
ctx := context.TODO()
id := subscription.NewSubscriptions2ID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "subscriptionValue")

read, err := client.GetEntityTag(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SubscriptionClient.List`

```go
ctx := context.TODO()
id := subscription.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue")

// alternatively `client.List(ctx, id, subscription.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, subscription.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `SubscriptionClient.ListSecrets`

```go
ctx := context.TODO()
id := subscription.NewSubscriptions2ID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "subscriptionValue")

read, err := client.ListSecrets(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SubscriptionClient.RegeneratePrimaryKey`

```go
ctx := context.TODO()
id := subscription.NewSubscriptions2ID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "subscriptionValue")

read, err := client.RegeneratePrimaryKey(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SubscriptionClient.RegenerateSecondaryKey`

```go
ctx := context.TODO()
id := subscription.NewSubscriptions2ID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "subscriptionValue")

read, err := client.RegenerateSecondaryKey(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SubscriptionClient.Update`

```go
ctx := context.TODO()
id := subscription.NewSubscriptions2ID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "subscriptionValue")

payload := subscription.SubscriptionUpdateParameters{
	// ...
}


read, err := client.Update(ctx, id, payload, subscription.DefaultUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SubscriptionClient.UserSubscriptionGet`

```go
ctx := context.TODO()
id := subscription.NewUserSubscriptions2ID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "userIdValue", "subscriptionValue")

read, err := client.UserSubscriptionGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
