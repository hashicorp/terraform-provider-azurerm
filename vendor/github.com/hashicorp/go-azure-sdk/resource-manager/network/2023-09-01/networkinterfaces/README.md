
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/networkinterfaces` Documentation

The `networkinterfaces` SDK allows for interaction with Azure Resource Manager `network` (API Version `2023-09-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/networkinterfaces"
```


### Client Initialization

```go
client := networkinterfaces.NewNetworkInterfacesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `NetworkInterfacesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := commonids.NewNetworkInterfaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkInterfaceName")

payload := networkinterfaces.NetworkInterface{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `NetworkInterfacesClient.Delete`

```go
ctx := context.TODO()
id := commonids.NewNetworkInterfaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkInterfaceName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `NetworkInterfacesClient.Get`

```go
ctx := context.TODO()
id := commonids.NewNetworkInterfaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkInterfaceName")

read, err := client.Get(ctx, id, networkinterfaces.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NetworkInterfacesClient.GetCloudServiceNetworkInterface`

```go
ctx := context.TODO()
id := networkinterfaces.NewRoleInstanceNetworkInterfaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cloudServiceName", "roleInstanceName", "networkInterfaceName")

read, err := client.GetCloudServiceNetworkInterface(ctx, id, networkinterfaces.DefaultGetCloudServiceNetworkInterfaceOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NetworkInterfacesClient.GetEffectiveRouteTable`

```go
ctx := context.TODO()
id := commonids.NewNetworkInterfaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkInterfaceName")

// alternatively `client.GetEffectiveRouteTable(ctx, id)` can be used to do batched pagination
items, err := client.GetEffectiveRouteTableComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `NetworkInterfacesClient.GetVirtualMachineScaleSetIPConfiguration`

```go
ctx := context.TODO()
id := commonids.NewVirtualMachineScaleSetIPConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineScaleSetName", "virtualMachineName", "networkInterfaceName", "ipConfigurationName")

read, err := client.GetVirtualMachineScaleSetIPConfiguration(ctx, id, networkinterfaces.DefaultGetVirtualMachineScaleSetIPConfigurationOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NetworkInterfacesClient.GetVirtualMachineScaleSetNetworkInterface`

```go
ctx := context.TODO()
id := commonids.NewVirtualMachineScaleSetNetworkInterfaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineScaleSetName", "virtualMachineName", "networkInterfaceName")

read, err := client.GetVirtualMachineScaleSetNetworkInterface(ctx, id, networkinterfaces.DefaultGetVirtualMachineScaleSetNetworkInterfaceOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NetworkInterfacesClient.List`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `NetworkInterfacesClient.ListAll`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListAll(ctx, id)` can be used to do batched pagination
items, err := client.ListAllComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `NetworkInterfacesClient.ListCloudServiceNetworkInterfaces`

```go
ctx := context.TODO()
id := networkinterfaces.NewProviderCloudServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cloudServiceName")

// alternatively `client.ListCloudServiceNetworkInterfaces(ctx, id)` can be used to do batched pagination
items, err := client.ListCloudServiceNetworkInterfacesComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `NetworkInterfacesClient.ListCloudServiceRoleInstanceNetworkInterfaces`

```go
ctx := context.TODO()
id := networkinterfaces.NewRoleInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cloudServiceName", "roleInstanceName")

// alternatively `client.ListCloudServiceRoleInstanceNetworkInterfaces(ctx, id)` can be used to do batched pagination
items, err := client.ListCloudServiceRoleInstanceNetworkInterfacesComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `NetworkInterfacesClient.ListEffectiveNetworkSecurityGroups`

```go
ctx := context.TODO()
id := commonids.NewNetworkInterfaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkInterfaceName")

// alternatively `client.ListEffectiveNetworkSecurityGroups(ctx, id)` can be used to do batched pagination
items, err := client.ListEffectiveNetworkSecurityGroupsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `NetworkInterfacesClient.ListVirtualMachineScaleSetIPConfigurations`

```go
ctx := context.TODO()
id := commonids.NewVirtualMachineScaleSetNetworkInterfaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineScaleSetName", "virtualMachineName", "networkInterfaceName")

// alternatively `client.ListVirtualMachineScaleSetIPConfigurations(ctx, id, networkinterfaces.DefaultListVirtualMachineScaleSetIPConfigurationsOperationOptions())` can be used to do batched pagination
items, err := client.ListVirtualMachineScaleSetIPConfigurationsComplete(ctx, id, networkinterfaces.DefaultListVirtualMachineScaleSetIPConfigurationsOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `NetworkInterfacesClient.ListVirtualMachineScaleSetNetworkInterfaces`

```go
ctx := context.TODO()
id := networkinterfaces.NewVirtualMachineScaleSetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineScaleSetName")

// alternatively `client.ListVirtualMachineScaleSetNetworkInterfaces(ctx, id)` can be used to do batched pagination
items, err := client.ListVirtualMachineScaleSetNetworkInterfacesComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `NetworkInterfacesClient.ListVirtualMachineScaleSetVMNetworkInterfaces`

```go
ctx := context.TODO()
id := networkinterfaces.NewVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineScaleSetName", "virtualMachineName")

// alternatively `client.ListVirtualMachineScaleSetVMNetworkInterfaces(ctx, id)` can be used to do batched pagination
items, err := client.ListVirtualMachineScaleSetVMNetworkInterfacesComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `NetworkInterfacesClient.NetworkInterfaceIPConfigurationsGet`

```go
ctx := context.TODO()
id := commonids.NewNetworkInterfaceIPConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkInterfaceName", "ipConfigurationName")

read, err := client.NetworkInterfaceIPConfigurationsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NetworkInterfacesClient.NetworkInterfaceIPConfigurationsList`

```go
ctx := context.TODO()
id := commonids.NewNetworkInterfaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkInterfaceName")

// alternatively `client.NetworkInterfaceIPConfigurationsList(ctx, id)` can be used to do batched pagination
items, err := client.NetworkInterfaceIPConfigurationsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `NetworkInterfacesClient.NetworkInterfaceLoadBalancersList`

```go
ctx := context.TODO()
id := commonids.NewNetworkInterfaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkInterfaceName")

// alternatively `client.NetworkInterfaceLoadBalancersList(ctx, id)` can be used to do batched pagination
items, err := client.NetworkInterfaceLoadBalancersListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `NetworkInterfacesClient.NetworkInterfaceTapConfigurationsCreateOrUpdate`

```go
ctx := context.TODO()
id := networkinterfaces.NewTapConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkInterfaceName", "tapConfigurationName")

payload := networkinterfaces.NetworkInterfaceTapConfiguration{
	// ...
}


if err := client.NetworkInterfaceTapConfigurationsCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `NetworkInterfacesClient.NetworkInterfaceTapConfigurationsDelete`

```go
ctx := context.TODO()
id := networkinterfaces.NewTapConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkInterfaceName", "tapConfigurationName")

if err := client.NetworkInterfaceTapConfigurationsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `NetworkInterfacesClient.NetworkInterfaceTapConfigurationsGet`

```go
ctx := context.TODO()
id := networkinterfaces.NewTapConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkInterfaceName", "tapConfigurationName")

read, err := client.NetworkInterfaceTapConfigurationsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NetworkInterfacesClient.NetworkInterfaceTapConfigurationsList`

```go
ctx := context.TODO()
id := commonids.NewNetworkInterfaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkInterfaceName")

// alternatively `client.NetworkInterfaceTapConfigurationsList(ctx, id)` can be used to do batched pagination
items, err := client.NetworkInterfaceTapConfigurationsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `NetworkInterfacesClient.UpdateTags`

```go
ctx := context.TODO()
id := commonids.NewNetworkInterfaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkInterfaceName")

payload := networkinterfaces.TagsObject{
	// ...
}


read, err := client.UpdateTags(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
