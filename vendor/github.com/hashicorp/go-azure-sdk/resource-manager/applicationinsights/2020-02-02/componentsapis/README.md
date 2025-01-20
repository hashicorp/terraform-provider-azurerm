
## `github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2020-02-02/componentsapis` Documentation

The `componentsapis` SDK allows for interaction with Azure Resource Manager `applicationinsights` (API Version `2020-02-02`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2020-02-02/componentsapis"
```


### Client Initialization

```go
client := componentsapis.NewComponentsAPIsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ComponentsAPIsClient.ComponentsCreateOrUpdate`

```go
ctx := context.TODO()
id := componentsapis.NewComponentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "componentName")

payload := componentsapis.ApplicationInsightsComponent{
	// ...
}


read, err := client.ComponentsCreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ComponentsAPIsClient.ComponentsDelete`

```go
ctx := context.TODO()
id := componentsapis.NewComponentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "componentName")

read, err := client.ComponentsDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ComponentsAPIsClient.ComponentsGet`

```go
ctx := context.TODO()
id := componentsapis.NewComponentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "componentName")

read, err := client.ComponentsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ComponentsAPIsClient.ComponentsGetPurgeStatus`

```go
ctx := context.TODO()
id := componentsapis.NewOperationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "componentName", "purgeId")

read, err := client.ComponentsGetPurgeStatus(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ComponentsAPIsClient.ComponentsList`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ComponentsList(ctx, id)` can be used to do batched pagination
items, err := client.ComponentsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ComponentsAPIsClient.ComponentsListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ComponentsListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ComponentsListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ComponentsAPIsClient.ComponentsPurge`

```go
ctx := context.TODO()
id := componentsapis.NewComponentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "componentName")

payload := componentsapis.ComponentPurgeBody{
	// ...
}


read, err := client.ComponentsPurge(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ComponentsAPIsClient.ComponentsUpdateTags`

```go
ctx := context.TODO()
id := componentsapis.NewComponentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "componentName")

payload := componentsapis.TagsResource{
	// ...
}


read, err := client.ComponentsUpdateTags(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
