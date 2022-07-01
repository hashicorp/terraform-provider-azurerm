
## `github.com/hashicorp/go-azure-sdk/resource-manager/appconfiguration/2022-05-01/configurationstores` Documentation

The `configurationstores` SDK allows for interaction with the Azure Resource Manager Service `appconfiguration` (API Version `2022-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/appconfiguration/2022-05-01/configurationstores"
```


### Client Initialization

```go
client := configurationstores.NewConfigurationStoresClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
if err != nil {
	// handle the error
}
```


### Example Usage: `ConfigurationStoresClient.Create`

```go
ctx := context.TODO()
id := configurationstores.NewConfigurationStoreID("12345678-1234-9876-4563-123456789012", "example-resource-group", "configStoreValue")

payload := configurationstores.ConfigurationStore{
	// ...
}

future, err := client.Create(ctx, id, payload)
if err != nil {
	// handle the error
}
if err := future.Poller.PollUntilDone(); err != nil {
	// handle the error
}
```


### Example Usage: `ConfigurationStoresClient.Delete`

```go
ctx := context.TODO()
id := configurationstores.NewConfigurationStoreID("12345678-1234-9876-4563-123456789012", "example-resource-group", "configStoreValue")
future, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if err := future.Poller.PollUntilDone(); err != nil {
	// handle the error
}
```


### Example Usage: `ConfigurationStoresClient.Get`

```go
ctx := context.TODO()
id := configurationstores.NewConfigurationStoreID("12345678-1234-9876-4563-123456789012", "example-resource-group", "configStoreValue")
read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ConfigurationStoresClient.List`

```go
ctx := context.TODO()
id := configurationstores.NewSubscriptionID()
// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ConfigurationStoresClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := configurationstores.NewResourceGroupID()
// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ConfigurationStoresClient.ListKeys`

```go
ctx := context.TODO()
id := configurationstores.NewConfigurationStoreID("12345678-1234-9876-4563-123456789012", "example-resource-group", "configStoreValue")
// alternatively `client.ListKeys(ctx, id)` can be used to do batched pagination
items, err := client.ListKeysComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ConfigurationStoresClient.RegenerateKey`

```go
ctx := context.TODO()
id := configurationstores.NewConfigurationStoreID("12345678-1234-9876-4563-123456789012", "example-resource-group", "configStoreValue")

payload := configurationstores.RegenerateKeyParameters{
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


### Example Usage: `ConfigurationStoresClient.Update`

```go
ctx := context.TODO()
id := configurationstores.NewConfigurationStoreID("12345678-1234-9876-4563-123456789012", "example-resource-group", "configStoreValue")

payload := configurationstores.ConfigurationStoreUpdateParameters{
	// ...
}

future, err := client.Update(ctx, id, payload)
if err != nil {
	// handle the error
}
if err := future.Poller.PollUntilDone(); err != nil {
	// handle the error
}
```
