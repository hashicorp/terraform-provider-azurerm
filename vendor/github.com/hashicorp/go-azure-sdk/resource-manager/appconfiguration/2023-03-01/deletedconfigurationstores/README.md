
## `github.com/hashicorp/go-azure-sdk/resource-manager/appconfiguration/2023-03-01/deletedconfigurationstores` Documentation

The `deletedconfigurationstores` SDK allows for interaction with Azure Resource Manager `appconfiguration` (API Version `2023-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/appconfiguration/2023-03-01/deletedconfigurationstores"
```


### Client Initialization

```go
client := deletedconfigurationstores.NewDeletedConfigurationStoresClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DeletedConfigurationStoresClient.ConfigurationStoresGetDeleted`

```go
ctx := context.TODO()
id := deletedconfigurationstores.NewDeletedConfigurationStoreID("12345678-1234-9876-4563-123456789012", "locationName", "deletedConfigurationStoreName")

read, err := client.ConfigurationStoresGetDeleted(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DeletedConfigurationStoresClient.ConfigurationStoresListDeleted`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ConfigurationStoresListDeleted(ctx, id)` can be used to do batched pagination
items, err := client.ConfigurationStoresListDeletedComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DeletedConfigurationStoresClient.ConfigurationStoresPurgeDeleted`

```go
ctx := context.TODO()
id := deletedconfigurationstores.NewDeletedConfigurationStoreID("12345678-1234-9876-4563-123456789012", "locationName", "deletedConfigurationStoreName")

if err := client.ConfigurationStoresPurgeDeletedThenPoll(ctx, id); err != nil {
	// handle the error
}
```
