
## `github.com/hashicorp/go-azure-sdk/resource-manager/automanage/2022-05-04/configurationprofilehcrpassignments` Documentation

The `configurationprofilehcrpassignments` SDK allows for interaction with Azure Resource Manager `automanage` (API Version `2022-05-04`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/automanage/2022-05-04/configurationprofilehcrpassignments"
```


### Client Initialization

```go
client := configurationprofilehcrpassignments.NewConfigurationProfileHCRPAssignmentsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ConfigurationProfileHCRPAssignmentsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := configurationprofilehcrpassignments.NewProviders2ConfigurationProfileAssignmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "machineName", "configurationProfileAssignmentName")

payload := configurationprofilehcrpassignments.ConfigurationProfileAssignment{
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


### Example Usage: `ConfigurationProfileHCRPAssignmentsClient.Delete`

```go
ctx := context.TODO()
id := configurationprofilehcrpassignments.NewProviders2ConfigurationProfileAssignmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "machineName", "configurationProfileAssignmentName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ConfigurationProfileHCRPAssignmentsClient.Get`

```go
ctx := context.TODO()
id := configurationprofilehcrpassignments.NewProviders2ConfigurationProfileAssignmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "machineName", "configurationProfileAssignmentName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
