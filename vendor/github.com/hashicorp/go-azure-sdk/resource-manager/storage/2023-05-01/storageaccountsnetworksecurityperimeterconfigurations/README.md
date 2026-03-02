
## `github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-05-01/storageaccountsnetworksecurityperimeterconfigurations` Documentation

The `storageaccountsnetworksecurityperimeterconfigurations` SDK allows for interaction with Azure Resource Manager `storage` (API Version `2023-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-05-01/storageaccountsnetworksecurityperimeterconfigurations"
```


### Client Initialization

```go
client := storageaccountsnetworksecurityperimeterconfigurations.NewStorageAccountsNetworkSecurityPerimeterConfigurationsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `StorageAccountsNetworkSecurityPerimeterConfigurationsClient.NetworkSecurityPerimeterConfigurationsGet`

```go
ctx := context.TODO()
id := storageaccountsnetworksecurityperimeterconfigurations.NewNetworkSecurityPerimeterConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountName", "networkSecurityPerimeterConfigurationName")

read, err := client.NetworkSecurityPerimeterConfigurationsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StorageAccountsNetworkSecurityPerimeterConfigurationsClient.NetworkSecurityPerimeterConfigurationsList`

```go
ctx := context.TODO()
id := commonids.NewStorageAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountName")

// alternatively `client.NetworkSecurityPerimeterConfigurationsList(ctx, id)` can be used to do batched pagination
items, err := client.NetworkSecurityPerimeterConfigurationsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `StorageAccountsNetworkSecurityPerimeterConfigurationsClient.NetworkSecurityPerimeterConfigurationsReconcile`

```go
ctx := context.TODO()
id := storageaccountsnetworksecurityperimeterconfigurations.NewNetworkSecurityPerimeterConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountName", "networkSecurityPerimeterConfigurationName")

if err := client.NetworkSecurityPerimeterConfigurationsReconcileThenPoll(ctx, id); err != nil {
	// handle the error
}
```
