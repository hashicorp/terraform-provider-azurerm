
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/dscpconfigurations` Documentation

The `dscpconfigurations` SDK allows for interaction with Azure Resource Manager `network` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/dscpconfigurations"
```


### Client Initialization

```go
client := dscpconfigurations.NewDscpConfigurationsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DscpConfigurationsClient.DscpConfigurationList`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.DscpConfigurationList(ctx, id)` can be used to do batched pagination
items, err := client.DscpConfigurationListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DscpConfigurationsClient.DscpConfigurationListAll`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.DscpConfigurationListAll(ctx, id)` can be used to do batched pagination
items, err := client.DscpConfigurationListAllComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
