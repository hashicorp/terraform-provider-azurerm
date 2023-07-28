
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/cloudservicepublicipaddresses` Documentation

The `cloudservicepublicipaddresses` SDK allows for interaction with the Azure Resource Manager Service `network` (API Version `2023-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/cloudservicepublicipaddresses"
```


### Client Initialization

```go
client := cloudservicepublicipaddresses.NewCloudServicePublicIPAddressesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `CloudServicePublicIPAddressesClient.PublicIPAddressesGetCloudServicePublicIPAddress`

```go
ctx := context.TODO()
id := cloudservicepublicipaddresses.NewCloudServicesPublicIPAddressID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cloudServiceValue", "roleInstanceValue", "networkInterfaceValue", "ipConfigurationValue", "publicIPAddressValue")

read, err := client.PublicIPAddressesGetCloudServicePublicIPAddress(ctx, id, cloudservicepublicipaddresses.DefaultPublicIPAddressesGetCloudServicePublicIPAddressOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CloudServicePublicIPAddressesClient.PublicIPAddressesListCloudServicePublicIPAddresses`

```go
ctx := context.TODO()
id := cloudservicepublicipaddresses.NewProviderCloudServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cloudServiceValue")

// alternatively `client.PublicIPAddressesListCloudServicePublicIPAddresses(ctx, id)` can be used to do batched pagination
items, err := client.PublicIPAddressesListCloudServicePublicIPAddressesComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `CloudServicePublicIPAddressesClient.PublicIPAddressesListCloudServiceRoleInstancePublicIPAddresses`

```go
ctx := context.TODO()
id := cloudservicepublicipaddresses.NewCloudServicesIPConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cloudServiceValue", "roleInstanceValue", "networkInterfaceValue", "ipConfigurationValue")

// alternatively `client.PublicIPAddressesListCloudServiceRoleInstancePublicIPAddresses(ctx, id)` can be used to do batched pagination
items, err := client.PublicIPAddressesListCloudServiceRoleInstancePublicIPAddressesComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
