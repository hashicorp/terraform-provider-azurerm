
## `github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/integrationaccounts` Documentation

The `integrationaccounts` SDK allows for interaction with the Azure Resource Manager Service `logic` (API Version `2019-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/integrationaccounts"
```


### Client Initialization

```go
client := integrationaccounts.NewIntegrationAccountsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `IntegrationAccountsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := integrationaccounts.NewIntegrationAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "integrationAccountValue")

payload := integrationaccounts.IntegrationAccount{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IntegrationAccountsClient.Delete`

```go
ctx := context.TODO()
id := integrationaccounts.NewIntegrationAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "integrationAccountValue")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IntegrationAccountsClient.Get`

```go
ctx := context.TODO()
id := integrationaccounts.NewIntegrationAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "integrationAccountValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IntegrationAccountsClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := integrationaccounts.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id, integrationaccounts.DefaultListByResourceGroupOperationOptions())` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id, integrationaccounts.DefaultListByResourceGroupOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `IntegrationAccountsClient.ListBySubscription`

```go
ctx := context.TODO()
id := integrationaccounts.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id, integrationaccounts.DefaultListBySubscriptionOperationOptions())` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id, integrationaccounts.DefaultListBySubscriptionOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `IntegrationAccountsClient.ListCallbackUrl`

```go
ctx := context.TODO()
id := integrationaccounts.NewIntegrationAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "integrationAccountValue")

payload := integrationaccounts.GetCallbackUrlParameters{
	// ...
}


read, err := client.ListCallbackUrl(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IntegrationAccountsClient.ListKeyVaultKeys`

```go
ctx := context.TODO()
id := integrationaccounts.NewIntegrationAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "integrationAccountValue")

payload := integrationaccounts.ListKeyVaultKeysDefinition{
	// ...
}


read, err := client.ListKeyVaultKeys(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IntegrationAccountsClient.LogTrackingEvents`

```go
ctx := context.TODO()
id := integrationaccounts.NewIntegrationAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "integrationAccountValue")

payload := integrationaccounts.TrackingEventsDefinition{
	// ...
}


read, err := client.LogTrackingEvents(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IntegrationAccountsClient.RegenerateAccessKey`

```go
ctx := context.TODO()
id := integrationaccounts.NewIntegrationAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "integrationAccountValue")

payload := integrationaccounts.RegenerateActionParameter{
	// ...
}


read, err := client.RegenerateAccessKey(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IntegrationAccountsClient.Update`

```go
ctx := context.TODO()
id := integrationaccounts.NewIntegrationAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "integrationAccountValue")

payload := integrationaccounts.IntegrationAccount{
	// ...
}


read, err := client.Update(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
