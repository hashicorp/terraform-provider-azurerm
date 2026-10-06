
## `github.com/hashicorp/go-azure-sdk/resource-manager/automanage/2022-05-04/configurationprofilehciassignments` Documentation

The `configurationprofilehciassignments` SDK allows for interaction with Azure Resource Manager `automanage` (API Version `2022-05-04`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/automanage/2022-05-04/configurationprofilehciassignments"
```


### Client Initialization

```go
client := configurationprofilehciassignments.NewConfigurationProfileHCIAssignmentsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ConfigurationProfileHCIAssignmentsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := configurationprofilehciassignments.NewConfigurationProfileAssignmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName", "configurationProfileAssignmentName")

payload := configurationprofilehciassignments.ConfigurationProfileAssignment{
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


### Example Usage: `ConfigurationProfileHCIAssignmentsClient.Delete`

```go
ctx := context.TODO()
id := configurationprofilehciassignments.NewConfigurationProfileAssignmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName", "configurationProfileAssignmentName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ConfigurationProfileHCIAssignmentsClient.Get`

```go
ctx := context.TODO()
id := configurationprofilehciassignments.NewConfigurationProfileAssignmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName", "configurationProfileAssignmentName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
