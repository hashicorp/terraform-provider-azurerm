
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/networksecuritygroups` Documentation

The `networksecuritygroups` SDK allows for interaction with Azure Resource Manager `network` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/networksecuritygroups"
```


### Client Initialization

```go
client := networksecuritygroups.NewNetworkSecurityGroupsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `NetworkSecurityGroupsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := networksecuritygroups.NewNetworkSecurityGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkSecurityGroupName")

payload := networksecuritygroups.NetworkSecurityGroup{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `NetworkSecurityGroupsClient.Delete`

```go
ctx := context.TODO()
id := networksecuritygroups.NewNetworkSecurityGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkSecurityGroupName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `NetworkSecurityGroupsClient.Get`

```go
ctx := context.TODO()
id := networksecuritygroups.NewNetworkSecurityGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkSecurityGroupName")

read, err := client.Get(ctx, id, networksecuritygroups.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NetworkSecurityGroupsClient.List`

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


### Example Usage: `NetworkSecurityGroupsClient.ListAll`

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


### Example Usage: `NetworkSecurityGroupsClient.UpdateTags`

```go
ctx := context.TODO()
id := networksecuritygroups.NewNetworkSecurityGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkSecurityGroupName")

payload := networksecuritygroups.TagsObject{
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
