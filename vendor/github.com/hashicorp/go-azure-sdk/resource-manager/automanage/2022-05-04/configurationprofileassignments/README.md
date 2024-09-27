
## `github.com/hashicorp/go-azure-sdk/resource-manager/automanage/2022-05-04/configurationprofileassignments` Documentation

The `configurationprofileassignments` SDK allows for interaction with Azure Resource Manager `automanage` (API Version `2022-05-04`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/automanage/2022-05-04/configurationprofileassignments"
```


### Client Initialization

```go
client := configurationprofileassignments.NewConfigurationProfileAssignmentsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ConfigurationProfileAssignmentsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := configurationprofileassignments.NewVirtualMachineProviders2ConfigurationProfileAssignmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineName", "configurationProfileAssignmentName")

payload := configurationprofileassignments.ConfigurationProfileAssignment{
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


### Example Usage: `ConfigurationProfileAssignmentsClient.Delete`

```go
ctx := context.TODO()
id := configurationprofileassignments.NewVirtualMachineProviders2ConfigurationProfileAssignmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineName", "configurationProfileAssignmentName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ConfigurationProfileAssignmentsClient.Get`

```go
ctx := context.TODO()
id := configurationprofileassignments.NewVirtualMachineProviders2ConfigurationProfileAssignmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineName", "configurationProfileAssignmentName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ConfigurationProfileAssignmentsClient.List`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

read, err := client.List(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ConfigurationProfileAssignmentsClient.ListByClusterName`

```go
ctx := context.TODO()
id := configurationprofileassignments.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName")

read, err := client.ListByClusterName(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ConfigurationProfileAssignmentsClient.ListByMachineName`

```go
ctx := context.TODO()
id := configurationprofileassignments.NewMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "machineName")

read, err := client.ListByMachineName(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ConfigurationProfileAssignmentsClient.ListBySubscription`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

read, err := client.ListBySubscription(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ConfigurationProfileAssignmentsClient.ListByVirtualMachines`

```go
ctx := context.TODO()
id := configurationprofileassignments.NewVirtualMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineName")

read, err := client.ListByVirtualMachines(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
