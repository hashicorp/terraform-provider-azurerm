
## `github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/namespaces` Documentation

The `namespaces` SDK allows for interaction with Azure Resource Manager `eventgrid` (API Version `2025-02-15`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/namespaces"
```


### Client Initialization

```go
client := namespaces.NewNamespacesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `NamespacesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := namespaces.NewNamespaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName")

payload := namespaces.Namespace{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `NamespacesClient.Delete`

```go
ctx := context.TODO()
id := namespaces.NewNamespaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `NamespacesClient.Get`

```go
ctx := context.TODO()
id := namespaces.NewNamespaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NamespacesClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id, namespaces.DefaultListByResourceGroupOperationOptions())` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id, namespaces.DefaultListByResourceGroupOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `NamespacesClient.ListBySubscription`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id, namespaces.DefaultListBySubscriptionOperationOptions())` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id, namespaces.DefaultListBySubscriptionOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `NamespacesClient.ListSharedAccessKeys`

```go
ctx := context.TODO()
id := namespaces.NewNamespaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName")

read, err := client.ListSharedAccessKeys(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NamespacesClient.RegenerateKey`

```go
ctx := context.TODO()
id := namespaces.NewNamespaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName")

payload := namespaces.NamespaceRegenerateKeyRequest{
	// ...
}


if err := client.RegenerateKeyThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `NamespacesClient.Update`

```go
ctx := context.TODO()
id := namespaces.NewNamespaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName")

payload := namespaces.NamespaceUpdateParameters{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `NamespacesClient.ValidateCustomDomainOwnership`

```go
ctx := context.TODO()
id := namespaces.NewNamespaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName")

if err := client.ValidateCustomDomainOwnershipThenPoll(ctx, id); err != nil {
	// handle the error
}
```
