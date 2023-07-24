
## `github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2022-07-01-preview/configurationassignments` Documentation

The `configurationassignments` SDK allows for interaction with the Azure Resource Manager Service `maintenance` (API Version `2022-07-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2022-07-01-preview/configurationassignments"
```


### Client Initialization

```go
client := configurationassignments.NewConfigurationAssignmentsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ConfigurationAssignmentsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := configurationassignments.NewScopedConfigurationAssignmentID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "configurationAssignmentValue")

payload := configurationassignments.ConfigurationAssignment{
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


### Example Usage: `ConfigurationAssignmentsClient.CreateOrUpdateParent`

```go
ctx := context.TODO()
id := configurationassignments.NewScopedConfigurationAssignmentID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "configurationAssignmentValue")

payload := configurationassignments.ConfigurationAssignment{
	// ...
}


read, err := client.CreateOrUpdateParent(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ConfigurationAssignmentsClient.Delete`

```go
ctx := context.TODO()
id := configurationassignments.NewScopedConfigurationAssignmentID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "configurationAssignmentValue")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ConfigurationAssignmentsClient.DeleteParent`

```go
ctx := context.TODO()
id := configurationassignments.NewScopedConfigurationAssignmentID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "configurationAssignmentValue")

read, err := client.DeleteParent(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ConfigurationAssignmentsClient.Get`

```go
ctx := context.TODO()
id := configurationassignments.NewScopedConfigurationAssignmentID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "configurationAssignmentValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ConfigurationAssignmentsClient.GetParent`

```go
ctx := context.TODO()
id := configurationassignments.NewScopedConfigurationAssignmentID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "configurationAssignmentValue")

read, err := client.GetParent(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ConfigurationAssignmentsClient.List`

```go
ctx := context.TODO()
id := configurationassignments.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

read, err := client.List(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ConfigurationAssignmentsClient.ListParent`

```go
ctx := context.TODO()
id := configurationassignments.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

read, err := client.ListParent(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ConfigurationAssignmentsClient.WithinSubscriptionList`

```go
ctx := context.TODO()
id := configurationassignments.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

read, err := client.WithinSubscriptionList(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
