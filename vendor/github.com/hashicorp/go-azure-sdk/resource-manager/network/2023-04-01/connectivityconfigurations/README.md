
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/connectivityconfigurations` Documentation

The `connectivityconfigurations` SDK allows for interaction with the Azure Resource Manager Service `network` (API Version `2023-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/connectivityconfigurations"
```


### Client Initialization

```go
client := connectivityconfigurations.NewConnectivityConfigurationsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ConnectivityConfigurationsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := connectivityconfigurations.NewConnectivityConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerValue", "connectivityConfigurationValue")

payload := connectivityconfigurations.ConnectivityConfiguration{
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


### Example Usage: `ConnectivityConfigurationsClient.Delete`

```go
ctx := context.TODO()
id := connectivityconfigurations.NewConnectivityConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerValue", "connectivityConfigurationValue")

if err := client.DeleteThenPoll(ctx, id, connectivityconfigurations.DefaultDeleteOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `ConnectivityConfigurationsClient.Get`

```go
ctx := context.TODO()
id := connectivityconfigurations.NewConnectivityConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerValue", "connectivityConfigurationValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ConnectivityConfigurationsClient.List`

```go
ctx := context.TODO()
id := connectivityconfigurations.NewNetworkManagerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerValue")

// alternatively `client.List(ctx, id, connectivityconfigurations.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, connectivityconfigurations.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
