
## `github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2023-05-01/netappaccounts` Documentation

The `netappaccounts` SDK allows for interaction with the Azure Resource Manager Service `netapp` (API Version `2023-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2023-05-01/netappaccounts"
```


### Client Initialization

```go
client := netappaccounts.NewNetAppAccountsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `NetAppAccountsClient.AccountsCreateOrUpdate`

```go
ctx := context.TODO()
id := netappaccounts.NewNetAppAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountValue")

payload := netappaccounts.NetAppAccount{
	// ...
}


if err := client.AccountsCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `NetAppAccountsClient.AccountsDelete`

```go
ctx := context.TODO()
id := netappaccounts.NewNetAppAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountValue")

if err := client.AccountsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `NetAppAccountsClient.AccountsGet`

```go
ctx := context.TODO()
id := netappaccounts.NewNetAppAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountValue")

read, err := client.AccountsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NetAppAccountsClient.AccountsList`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.AccountsList(ctx, id)` can be used to do batched pagination
items, err := client.AccountsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `NetAppAccountsClient.AccountsListBySubscription`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.AccountsListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.AccountsListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `NetAppAccountsClient.AccountsRenewCredentials`

```go
ctx := context.TODO()
id := netappaccounts.NewNetAppAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountValue")

if err := client.AccountsRenewCredentialsThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `NetAppAccountsClient.AccountsUpdate`

```go
ctx := context.TODO()
id := netappaccounts.NewNetAppAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountValue")

payload := netappaccounts.NetAppAccountPatch{
	// ...
}


if err := client.AccountsUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
