
## `github.com/hashicorp/go-azure-sdk/resource-manager/resources/2020-05-01/resourcemanagementprivatelink` Documentation

The `resourcemanagementprivatelink` SDK allows for interaction with Azure Resource Manager `resources` (API Version `2020-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/resources/2020-05-01/resourcemanagementprivatelink"
```


### Client Initialization

```go
client := resourcemanagementprivatelink.NewResourceManagementPrivateLinkClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ResourceManagementPrivateLinkClient.Delete`

```go
ctx := context.TODO()
id := resourcemanagementprivatelink.NewResourceManagementPrivateLinkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "resourceManagementPrivateLinkName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ResourceManagementPrivateLinkClient.Get`

```go
ctx := context.TODO()
id := resourcemanagementprivatelink.NewResourceManagementPrivateLinkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "resourceManagementPrivateLinkName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ResourceManagementPrivateLinkClient.List`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

read, err := client.List(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ResourceManagementPrivateLinkClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

read, err := client.ListByResourceGroup(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ResourceManagementPrivateLinkClient.Put`

```go
ctx := context.TODO()
id := resourcemanagementprivatelink.NewResourceManagementPrivateLinkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "resourceManagementPrivateLinkName")

payload := resourcemanagementprivatelink.ResourceManagementPrivateLinkLocation{
	// ...
}


read, err := client.Put(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
