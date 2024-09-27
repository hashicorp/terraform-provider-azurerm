
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/networkmanagereffectiveconnectivityconfiguration` Documentation

The `networkmanagereffectiveconnectivityconfiguration` SDK allows for interaction with Azure Resource Manager `network` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/networkmanagereffectiveconnectivityconfiguration"
```


### Client Initialization

```go
client := networkmanagereffectiveconnectivityconfiguration.NewNetworkManagerEffectiveConnectivityConfigurationClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `NetworkManagerEffectiveConnectivityConfigurationClient.ListNetworkManagerEffectiveConnectivityConfigurations`

```go
ctx := context.TODO()
id := commonids.NewVirtualNetworkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualNetworkName")

payload := networkmanagereffectiveconnectivityconfiguration.QueryRequestOptions{
	// ...
}


read, err := client.ListNetworkManagerEffectiveConnectivityConfigurations(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
