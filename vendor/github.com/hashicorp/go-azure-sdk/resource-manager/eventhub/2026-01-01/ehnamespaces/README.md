
## `github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2026-01-01/ehnamespaces` Documentation

The `ehnamespaces` SDK allows for interaction with Azure Resource Manager `eventhub` (API Version `2026-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2026-01-01/ehnamespaces"
```


### Client Initialization

```go
client := ehnamespaces.NewEHNamespacesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `EHNamespacesClient.DisasterRecoveryConfigsCheckNameAvailability`

```go
ctx := context.TODO()
id := ehnamespaces.NewNamespaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName")

payload := ehnamespaces.CheckNameAvailabilityParameter{
	// ...
}


read, err := client.DisasterRecoveryConfigsCheckNameAvailability(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EHNamespacesClient.NamespacesCreateOrUpdate`

```go
ctx := context.TODO()
id := ehnamespaces.NewNamespaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName")

payload := ehnamespaces.EHNamespace{
	// ...
}


if err := client.NamespacesCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `EHNamespacesClient.NamespacesDelete`

```go
ctx := context.TODO()
id := ehnamespaces.NewNamespaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName")

if err := client.NamespacesDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `EHNamespacesClient.NamespacesFailover`

```go
ctx := context.TODO()
id := ehnamespaces.NewNamespaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName")

payload := ehnamespaces.FailOver{
	// ...
}


if err := client.NamespacesFailoverThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `EHNamespacesClient.NamespacesGet`

```go
ctx := context.TODO()
id := ehnamespaces.NewNamespaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName")

read, err := client.NamespacesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EHNamespacesClient.NamespacesList`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.NamespacesList(ctx, id)` can be used to do batched pagination
items, err := client.NamespacesListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `EHNamespacesClient.NamespacesListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.NamespacesListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.NamespacesListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `EHNamespacesClient.NamespacesUpdate`

```go
ctx := context.TODO()
id := ehnamespaces.NewNamespaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName")

payload := ehnamespaces.EHNamespace{
	// ...
}


read, err := client.NamespacesUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EHNamespacesClient.PrivateLinkResourcesGet`

```go
ctx := context.TODO()
id := ehnamespaces.NewNamespaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName")

// alternatively `client.PrivateLinkResourcesGet(ctx, id)` can be used to do batched pagination
items, err := client.PrivateLinkResourcesGetComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
