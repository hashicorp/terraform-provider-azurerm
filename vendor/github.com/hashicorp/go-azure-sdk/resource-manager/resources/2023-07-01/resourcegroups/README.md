
## `github.com/hashicorp/go-azure-sdk/resource-manager/resources/2023-07-01/resourcegroups` Documentation

The `resourcegroups` SDK allows for interaction with Azure Resource Manager `resources` (API Version `2023-07-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/resources/2023-07-01/resourcegroups"
```


### Client Initialization

```go
client := resourcegroups.NewResourceGroupsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ResourceGroupsClient.CheckExistence`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

read, err := client.CheckExistence(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ResourceGroupsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

payload := resourcegroups.ResourceGroup{
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


### Example Usage: `ResourceGroupsClient.Delete`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

if err := client.DeleteThenPoll(ctx, id, resourcegroups.DefaultDeleteOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `ResourceGroupsClient.ExportTemplate`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

payload := resourcegroups.ExportTemplateRequest{
	// ...
}


if err := client.ExportTemplateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ResourceGroupsClient.Get`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ResourceGroupsClient.List`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.List(ctx, id, resourcegroups.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, resourcegroups.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ResourceGroupsClient.ResourcesListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ResourcesListByResourceGroup(ctx, id, resourcegroups.DefaultResourcesListByResourceGroupOperationOptions())` can be used to do batched pagination
items, err := client.ResourcesListByResourceGroupComplete(ctx, id, resourcegroups.DefaultResourcesListByResourceGroupOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ResourceGroupsClient.Update`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

payload := resourcegroups.ResourceGroupPatchable{
	// ...
}


read, err := client.Update(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
