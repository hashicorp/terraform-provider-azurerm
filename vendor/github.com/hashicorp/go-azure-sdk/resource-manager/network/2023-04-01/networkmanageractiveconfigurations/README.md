
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/networkmanageractiveconfigurations` Documentation

The `networkmanageractiveconfigurations` SDK allows for interaction with the Azure Resource Manager Service `network` (API Version `2023-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/networkmanageractiveconfigurations"
```


### Client Initialization

```go
client := networkmanageractiveconfigurations.NewNetworkManagerActiveConfigurationsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `NetworkManagerActiveConfigurationsClient.ListActiveSecurityAdminRules`

```go
ctx := context.TODO()
id := networkmanageractiveconfigurations.NewNetworkManagerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerValue")

payload := networkmanageractiveconfigurations.ActiveConfigurationParameter{
	// ...
}


read, err := client.ListActiveSecurityAdminRules(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
