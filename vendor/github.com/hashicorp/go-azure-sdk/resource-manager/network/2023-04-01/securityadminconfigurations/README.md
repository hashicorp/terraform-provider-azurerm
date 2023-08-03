
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/securityadminconfigurations` Documentation

The `securityadminconfigurations` SDK allows for interaction with the Azure Resource Manager Service `network` (API Version `2023-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/securityadminconfigurations"
```


### Client Initialization

```go
client := securityadminconfigurations.NewSecurityAdminConfigurationsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SecurityAdminConfigurationsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := securityadminconfigurations.NewSecurityAdminConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerValue", "securityAdminConfigurationValue")

payload := securityadminconfigurations.SecurityAdminConfiguration{
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


### Example Usage: `SecurityAdminConfigurationsClient.Delete`

```go
ctx := context.TODO()
id := securityadminconfigurations.NewSecurityAdminConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerValue", "securityAdminConfigurationValue")

if err := client.DeleteThenPoll(ctx, id, securityadminconfigurations.DefaultDeleteOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `SecurityAdminConfigurationsClient.Get`

```go
ctx := context.TODO()
id := securityadminconfigurations.NewSecurityAdminConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerValue", "securityAdminConfigurationValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SecurityAdminConfigurationsClient.List`

```go
ctx := context.TODO()
id := securityadminconfigurations.NewNetworkManagerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerValue")

// alternatively `client.List(ctx, id, securityadminconfigurations.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, securityadminconfigurations.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
