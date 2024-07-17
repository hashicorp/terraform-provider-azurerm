
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2023-05-01-preview/subscription` Documentation

The `subscription` SDK allows for interaction with the Azure Resource Manager Service `apimanagement` (API Version `2023-05-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2023-05-01-preview/subscription"
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


### Example Usage: `SubscriptionClient.WorkspaceSubscriptionCreateOrUpdate`

```go
ctx := context.TODO()
id := subscription.NewWorkspaceSubscriptions2ID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "workspaceIdValue", "subscriptionValue")

payload := subscription.SubscriptionCreateParameters{
	// ...
}


read, err := client.WorkspaceSubscriptionCreateOrUpdate(ctx, id, payload, subscription.DefaultWorkspaceSubscriptionCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SubscriptionClient.WorkspaceSubscriptionDelete`

```go
ctx := context.TODO()
id := subscription.NewWorkspaceSubscriptions2ID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "workspaceIdValue", "subscriptionValue")

read, err := client.WorkspaceSubscriptionDelete(ctx, id, subscription.DefaultWorkspaceSubscriptionDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SubscriptionClient.WorkspaceSubscriptionGet`

```go
ctx := context.TODO()
id := subscription.NewWorkspaceSubscriptions2ID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "workspaceIdValue", "subscriptionValue")

read, err := client.WorkspaceSubscriptionGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SubscriptionClient.WorkspaceSubscriptionGetEntityTag`

```go
ctx := context.TODO()
id := subscription.NewWorkspaceSubscriptions2ID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "workspaceIdValue", "subscriptionValue")

read, err := client.WorkspaceSubscriptionGetEntityTag(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SubscriptionClient.WorkspaceSubscriptionList`

```go
ctx := context.TODO()
id := subscription.NewWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "workspaceIdValue")

// alternatively `client.WorkspaceSubscriptionList(ctx, id, subscription.DefaultWorkspaceSubscriptionListOperationOptions())` can be used to do batched pagination
items, err := client.WorkspaceSubscriptionListComplete(ctx, id, subscription.DefaultWorkspaceSubscriptionListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `SubscriptionClient.WorkspaceSubscriptionListSecrets`

```go
ctx := context.TODO()
id := subscription.NewWorkspaceSubscriptions2ID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "workspaceIdValue", "subscriptionValue")

read, err := client.WorkspaceSubscriptionListSecrets(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SubscriptionClient.WorkspaceSubscriptionRegeneratePrimaryKey`

```go
ctx := context.TODO()
id := subscription.NewWorkspaceSubscriptions2ID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "workspaceIdValue", "subscriptionValue")

read, err := client.WorkspaceSubscriptionRegeneratePrimaryKey(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SubscriptionClient.WorkspaceSubscriptionRegenerateSecondaryKey`

```go
ctx := context.TODO()
id := subscription.NewWorkspaceSubscriptions2ID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "workspaceIdValue", "subscriptionValue")

read, err := client.WorkspaceSubscriptionRegenerateSecondaryKey(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SubscriptionClient.WorkspaceSubscriptionUpdate`

```go
ctx := context.TODO()
id := subscription.NewWorkspaceSubscriptions2ID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "workspaceIdValue", "subscriptionValue")

payload := subscription.SubscriptionUpdateParameters{
	// ...
}


read, err := client.WorkspaceSubscriptionUpdate(ctx, id, payload, subscription.DefaultWorkspaceSubscriptionUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
