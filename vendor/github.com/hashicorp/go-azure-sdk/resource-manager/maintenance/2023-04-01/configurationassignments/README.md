
## `github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2023-04-01/configurationassignments` Documentation

The `configurationassignments` SDK allows for interaction with Azure Resource Manager `maintenance` (API Version `2023-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2023-04-01/configurationassignments"
```


### Client Initialization

```go
client := configurationassignments.NewConfigurationAssignmentsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ConfigurationAssignmentsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := configurationassignments.NewScopedConfigurationAssignmentID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "configurationAssignmentName")

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
id := configurationassignments.NewScopedConfigurationAssignmentID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "configurationAssignmentName")

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
id := configurationassignments.NewScopedConfigurationAssignmentID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "configurationAssignmentName")

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
id := configurationassignments.NewScopedConfigurationAssignmentID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "configurationAssignmentName")

read, err := client.DeleteParent(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ConfigurationAssignmentsClient.ForResourceGroupCreateOrUpdate`

```go
ctx := context.TODO()
id := configurationassignments.NewProviderConfigurationAssignmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "configurationAssignmentName")

payload := configurationassignments.ConfigurationAssignment{
	// ...
}


read, err := client.ForResourceGroupCreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ConfigurationAssignmentsClient.ForResourceGroupDelete`

```go
ctx := context.TODO()
id := configurationassignments.NewProviderConfigurationAssignmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "configurationAssignmentName")

read, err := client.ForResourceGroupDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ConfigurationAssignmentsClient.ForResourceGroupGet`

```go
ctx := context.TODO()
id := configurationassignments.NewProviderConfigurationAssignmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "configurationAssignmentName")

read, err := client.ForResourceGroupGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ConfigurationAssignmentsClient.ForResourceGroupUpdate`

```go
ctx := context.TODO()
id := configurationassignments.NewProviderConfigurationAssignmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "configurationAssignmentName")

payload := configurationassignments.ConfigurationAssignment{
	// ...
}


read, err := client.ForResourceGroupUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ConfigurationAssignmentsClient.ForSubscriptionsCreateOrUpdate`

```go
ctx := context.TODO()
id := configurationassignments.NewConfigurationAssignmentID("12345678-1234-9876-4563-123456789012", "configurationAssignmentName")

payload := configurationassignments.ConfigurationAssignment{
	// ...
}


read, err := client.ForSubscriptionsCreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ConfigurationAssignmentsClient.ForSubscriptionsDelete`

```go
ctx := context.TODO()
id := configurationassignments.NewConfigurationAssignmentID("12345678-1234-9876-4563-123456789012", "configurationAssignmentName")

read, err := client.ForSubscriptionsDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ConfigurationAssignmentsClient.ForSubscriptionsGet`

```go
ctx := context.TODO()
id := configurationassignments.NewConfigurationAssignmentID("12345678-1234-9876-4563-123456789012", "configurationAssignmentName")

read, err := client.ForSubscriptionsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ConfigurationAssignmentsClient.ForSubscriptionsUpdate`

```go
ctx := context.TODO()
id := configurationassignments.NewConfigurationAssignmentID("12345678-1234-9876-4563-123456789012", "configurationAssignmentName")

payload := configurationassignments.ConfigurationAssignment{
	// ...
}


read, err := client.ForSubscriptionsUpdate(ctx, id, payload)
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
id := configurationassignments.NewScopedConfigurationAssignmentID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "configurationAssignmentName")

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
id := configurationassignments.NewScopedConfigurationAssignmentID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "configurationAssignmentName")

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
id := commonids.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

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
id := commonids.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

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
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

read, err := client.WithinSubscriptionList(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
