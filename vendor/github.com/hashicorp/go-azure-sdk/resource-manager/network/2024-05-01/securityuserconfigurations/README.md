
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-05-01/securityuserconfigurations` Documentation

The `securityuserconfigurations` SDK allows for interaction with Azure Resource Manager `network` (API Version `2024-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-05-01/securityuserconfigurations"
```


### Client Initialization

```go
client := securityuserconfigurations.NewSecurityUserConfigurationsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SecurityUserConfigurationsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := securityuserconfigurations.NewSecurityUserConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerName", "securityUserConfigurationName")

payload := securityuserconfigurations.SecurityUserConfiguration{
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


### Example Usage: `SecurityUserConfigurationsClient.Delete`

```go
ctx := context.TODO()
id := securityuserconfigurations.NewSecurityUserConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerName", "securityUserConfigurationName")

if err := client.DeleteThenPoll(ctx, id, securityuserconfigurations.DefaultDeleteOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `SecurityUserConfigurationsClient.Get`

```go
ctx := context.TODO()
id := securityuserconfigurations.NewSecurityUserConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerName", "securityUserConfigurationName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SecurityUserConfigurationsClient.List`

```go
ctx := context.TODO()
id := securityuserconfigurations.NewNetworkManagerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerName")

// alternatively `client.List(ctx, id, securityuserconfigurations.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, securityuserconfigurations.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
