
## `github.com/hashicorp/go-azure-sdk/resource-manager/subscription/2021-10-01/subscriptions` Documentation

The `subscriptions` SDK allows for interaction with Azure Resource Manager `subscription` (API Version `2021-10-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/subscription/2021-10-01/subscriptions"
```


### Client Initialization

```go
client := subscriptions.NewSubscriptionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SubscriptionsClient.AliasCreate`

```go
ctx := context.TODO()
id := subscriptions.NewAliasID("aliasName")

payload := subscriptions.PutAliasRequest{
	// ...
}


if err := client.AliasCreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `SubscriptionsClient.AliasDelete`

```go
ctx := context.TODO()
id := subscriptions.NewAliasID("aliasName")

read, err := client.AliasDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SubscriptionsClient.AliasGet`

```go
ctx := context.TODO()
id := subscriptions.NewAliasID("aliasName")

read, err := client.AliasGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SubscriptionsClient.AliasList`

```go
ctx := context.TODO()


// alternatively `client.AliasList(ctx)` can be used to do batched pagination
items, err := client.AliasListComplete(ctx)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `SubscriptionsClient.BillingAccountGetPolicy`

```go
ctx := context.TODO()
id := subscriptions.NewBillingAccountID("billingAccountId")

read, err := client.BillingAccountGetPolicy(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SubscriptionsClient.SubscriptionAcceptOwnership`

```go
ctx := context.TODO()
id := subscriptions.NewProviderSubscriptionID("12345678-1234-9876-4563-123456789012")

payload := subscriptions.AcceptOwnershipRequest{
	// ...
}


if err := client.SubscriptionAcceptOwnershipThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `SubscriptionsClient.SubscriptionAcceptOwnershipStatus`

```go
ctx := context.TODO()
id := subscriptions.NewProviderSubscriptionID("12345678-1234-9876-4563-123456789012")

read, err := client.SubscriptionAcceptOwnershipStatus(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SubscriptionsClient.SubscriptionCancel`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

read, err := client.SubscriptionCancel(ctx, id, subscriptions.DefaultSubscriptionCancelOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SubscriptionsClient.SubscriptionEnable`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

read, err := client.SubscriptionEnable(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SubscriptionsClient.SubscriptionOperationGet`

```go
ctx := context.TODO()
id := subscriptions.NewSubscriptionOperationID("operationId")

read, err := client.SubscriptionOperationGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SubscriptionsClient.SubscriptionPolicyAddUpdatePolicyForTenant`

```go
ctx := context.TODO()

payload := subscriptions.PutTenantPolicyRequestProperties{
	// ...
}


read, err := client.SubscriptionPolicyAddUpdatePolicyForTenant(ctx, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SubscriptionsClient.SubscriptionPolicyGetPolicyForTenant`

```go
ctx := context.TODO()


read, err := client.SubscriptionPolicyGetPolicyForTenant(ctx)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SubscriptionsClient.SubscriptionPolicyListPolicyForTenant`

```go
ctx := context.TODO()


// alternatively `client.SubscriptionPolicyListPolicyForTenant(ctx)` can be used to do batched pagination
items, err := client.SubscriptionPolicyListPolicyForTenantComplete(ctx)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `SubscriptionsClient.SubscriptionRename`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

payload := subscriptions.SubscriptionName{
	// ...
}


read, err := client.SubscriptionRename(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
