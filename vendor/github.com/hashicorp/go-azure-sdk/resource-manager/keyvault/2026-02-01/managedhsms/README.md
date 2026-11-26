
## `github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2026-02-01/managedhsms` Documentation

The `managedhsms` SDK allows for interaction with Azure Resource Manager `keyvault` (API Version `2026-02-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2026-02-01/managedhsms"
```


### Client Initialization

```go
client := managedhsms.NewManagedHsmsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ManagedHsmsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := managedhsms.NewManagedHSMID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedHSMName")

payload := managedhsms.ManagedHsm{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ManagedHsmsClient.Delete`

```go
ctx := context.TODO()
id := managedhsms.NewManagedHSMID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedHSMName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ManagedHsmsClient.Get`

```go
ctx := context.TODO()
id := managedhsms.NewManagedHSMID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedHSMName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagedHsmsClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id, managedhsms.DefaultListByResourceGroupOperationOptions())` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id, managedhsms.DefaultListByResourceGroupOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ManagedHsmsClient.ListBySubscription`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id, managedhsms.DefaultListBySubscriptionOperationOptions())` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id, managedhsms.DefaultListBySubscriptionOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ManagedHsmsClient.MHSMPrivateLinkResourcesListByMHSMResource`

```go
ctx := context.TODO()
id := managedhsms.NewManagedHSMID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedHSMName")

read, err := client.MHSMPrivateLinkResourcesListByMHSMResource(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagedHsmsClient.MHSMRegionsListByResource`

```go
ctx := context.TODO()
id := managedhsms.NewManagedHSMID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedHSMName")

// alternatively `client.MHSMRegionsListByResource(ctx, id)` can be used to do batched pagination
items, err := client.MHSMRegionsListByResourceComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ManagedHsmsClient.Update`

```go
ctx := context.TODO()
id := managedhsms.NewManagedHSMID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedHSMName")

payload := managedhsms.ManagedHsm{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
