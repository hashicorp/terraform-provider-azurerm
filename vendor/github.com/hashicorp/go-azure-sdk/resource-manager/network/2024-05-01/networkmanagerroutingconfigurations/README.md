
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-05-01/networkmanagerroutingconfigurations` Documentation

The `networkmanagerroutingconfigurations` SDK allows for interaction with Azure Resource Manager `network` (API Version `2024-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-05-01/networkmanagerroutingconfigurations"
```


### Client Initialization

```go
client := networkmanagerroutingconfigurations.NewNetworkManagerRoutingConfigurationsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `NetworkManagerRoutingConfigurationsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := networkmanagerroutingconfigurations.NewRoutingConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerName", "routingConfigurationName")

payload := networkmanagerroutingconfigurations.NetworkManagerRoutingConfiguration{
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


### Example Usage: `NetworkManagerRoutingConfigurationsClient.Delete`

```go
ctx := context.TODO()
id := networkmanagerroutingconfigurations.NewRoutingConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerName", "routingConfigurationName")

if err := client.DeleteThenPoll(ctx, id, networkmanagerroutingconfigurations.DefaultDeleteOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `NetworkManagerRoutingConfigurationsClient.Get`

```go
ctx := context.TODO()
id := networkmanagerroutingconfigurations.NewRoutingConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerName", "routingConfigurationName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NetworkManagerRoutingConfigurationsClient.List`

```go
ctx := context.TODO()
id := networkmanagerroutingconfigurations.NewNetworkManagerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerName")

// alternatively `client.List(ctx, id, networkmanagerroutingconfigurations.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, networkmanagerroutingconfigurations.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
