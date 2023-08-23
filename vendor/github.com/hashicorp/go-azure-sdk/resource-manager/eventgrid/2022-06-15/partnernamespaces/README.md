
## `github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2022-06-15/partnernamespaces` Documentation

The `partnernamespaces` SDK allows for interaction with the Azure Resource Manager Service `eventgrid` (API Version `2022-06-15`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2022-06-15/partnernamespaces"
```


### Client Initialization

```go
client := partnernamespaces.NewPartnerNamespacesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `PartnerNamespacesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := partnernamespaces.NewPartnerNamespaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "partnerNamespaceValue")

payload := partnernamespaces.PartnerNamespace{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `PartnerNamespacesClient.Delete`

```go
ctx := context.TODO()
id := partnernamespaces.NewPartnerNamespaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "partnerNamespaceValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `PartnerNamespacesClient.Get`

```go
ctx := context.TODO()
id := partnernamespaces.NewPartnerNamespaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "partnerNamespaceValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PartnerNamespacesClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := partnernamespaces.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id, partnernamespaces.DefaultListByResourceGroupOperationOptions())` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id, partnernamespaces.DefaultListByResourceGroupOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `PartnerNamespacesClient.ListBySubscription`

```go
ctx := context.TODO()
id := partnernamespaces.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id, partnernamespaces.DefaultListBySubscriptionOperationOptions())` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id, partnernamespaces.DefaultListBySubscriptionOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `PartnerNamespacesClient.ListSharedAccessKeys`

```go
ctx := context.TODO()
id := partnernamespaces.NewPartnerNamespaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "partnerNamespaceValue")

read, err := client.ListSharedAccessKeys(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PartnerNamespacesClient.RegenerateKey`

```go
ctx := context.TODO()
id := partnernamespaces.NewPartnerNamespaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "partnerNamespaceValue")

payload := partnernamespaces.PartnerNamespaceRegenerateKeyRequest{
	// ...
}


read, err := client.RegenerateKey(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PartnerNamespacesClient.Update`

```go
ctx := context.TODO()
id := partnernamespaces.NewPartnerNamespaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "partnerNamespaceValue")

payload := partnernamespaces.PartnerNamespaceUpdateParameters{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
