
## `github.com/hashicorp/go-azure-sdk/resource-manager/resources/2020-05-01/managementlocks` Documentation

The `managementlocks` SDK allows for interaction with the Azure Resource Manager Service `resources` (API Version `2020-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/resources/2020-05-01/managementlocks"
```


### Client Initialization

```go
client := managementlocks.NewManagementLocksClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ManagementLocksClient.CreateOrUpdateAtResourceGroupLevel`

```go
ctx := context.TODO()
id := managementlocks.NewProviderLockID("12345678-1234-9876-4563-123456789012", "example-resource-group", "lockValue")

payload := managementlocks.ManagementLockObject{
	// ...
}


read, err := client.CreateOrUpdateAtResourceGroupLevel(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagementLocksClient.CreateOrUpdateAtResourceLevel`

```go
ctx := context.TODO()
id := managementlocks.NewScopedLockID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "lockValue")

payload := managementlocks.ManagementLockObject{
	// ...
}


read, err := client.CreateOrUpdateAtResourceLevel(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagementLocksClient.CreateOrUpdateAtSubscriptionLevel`

```go
ctx := context.TODO()
id := managementlocks.NewLockID("12345678-1234-9876-4563-123456789012", "lockValue")

payload := managementlocks.ManagementLockObject{
	// ...
}


read, err := client.CreateOrUpdateAtSubscriptionLevel(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagementLocksClient.CreateOrUpdateByScope`

```go
ctx := context.TODO()
id := managementlocks.NewScopedLockID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "lockValue")

payload := managementlocks.ManagementLockObject{
	// ...
}


read, err := client.CreateOrUpdateByScope(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagementLocksClient.DeleteAtResourceGroupLevel`

```go
ctx := context.TODO()
id := managementlocks.NewProviderLockID("12345678-1234-9876-4563-123456789012", "example-resource-group", "lockValue")

read, err := client.DeleteAtResourceGroupLevel(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagementLocksClient.DeleteAtResourceLevel`

```go
ctx := context.TODO()
id := managementlocks.NewScopedLockID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "lockValue")

read, err := client.DeleteAtResourceLevel(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagementLocksClient.DeleteAtSubscriptionLevel`

```go
ctx := context.TODO()
id := managementlocks.NewLockID("12345678-1234-9876-4563-123456789012", "lockValue")

read, err := client.DeleteAtSubscriptionLevel(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagementLocksClient.DeleteByScope`

```go
ctx := context.TODO()
id := managementlocks.NewScopedLockID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "lockValue")

read, err := client.DeleteByScope(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagementLocksClient.GetAtResourceGroupLevel`

```go
ctx := context.TODO()
id := managementlocks.NewProviderLockID("12345678-1234-9876-4563-123456789012", "example-resource-group", "lockValue")

read, err := client.GetAtResourceGroupLevel(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagementLocksClient.GetAtResourceLevel`

```go
ctx := context.TODO()
id := managementlocks.NewScopedLockID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "lockValue")

read, err := client.GetAtResourceLevel(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagementLocksClient.GetAtSubscriptionLevel`

```go
ctx := context.TODO()
id := managementlocks.NewLockID("12345678-1234-9876-4563-123456789012", "lockValue")

read, err := client.GetAtSubscriptionLevel(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagementLocksClient.GetByScope`

```go
ctx := context.TODO()
id := managementlocks.NewScopedLockID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "lockValue")

read, err := client.GetByScope(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagementLocksClient.ListAtResourceGroupLevel`

```go
ctx := context.TODO()
id := managementlocks.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListAtResourceGroupLevel(ctx, id, managementlocks.DefaultListAtResourceGroupLevelOperationOptions())` can be used to do batched pagination
items, err := client.ListAtResourceGroupLevelComplete(ctx, id, managementlocks.DefaultListAtResourceGroupLevelOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ManagementLocksClient.ListAtResourceLevel`

```go
ctx := context.TODO()
id := managementlocks.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

// alternatively `client.ListAtResourceLevel(ctx, id, managementlocks.DefaultListAtResourceLevelOperationOptions())` can be used to do batched pagination
items, err := client.ListAtResourceLevelComplete(ctx, id, managementlocks.DefaultListAtResourceLevelOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ManagementLocksClient.ListAtSubscriptionLevel`

```go
ctx := context.TODO()
id := managementlocks.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListAtSubscriptionLevel(ctx, id, managementlocks.DefaultListAtSubscriptionLevelOperationOptions())` can be used to do batched pagination
items, err := client.ListAtSubscriptionLevelComplete(ctx, id, managementlocks.DefaultListAtSubscriptionLevelOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ManagementLocksClient.ListByScope`

```go
ctx := context.TODO()
id := managementlocks.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

// alternatively `client.ListByScope(ctx, id, managementlocks.DefaultListByScopeOperationOptions())` can be used to do batched pagination
items, err := client.ListByScopeComplete(ctx, id, managementlocks.DefaultListByScopeOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
