
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/ipgroups` Documentation

The `ipgroups` SDK allows for interaction with the Azure Resource Manager Service `network` (API Version `2023-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/ipgroups"
```


### Client Initialization

```go
client := ipgroups.NewIPGroupsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `IPGroupsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := ipgroups.NewIPGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "ipGroupValue")

payload := ipgroups.IPGroup{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `IPGroupsClient.Delete`

```go
ctx := context.TODO()
id := ipgroups.NewIPGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "ipGroupValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `IPGroupsClient.Get`

```go
ctx := context.TODO()
id := ipgroups.NewIPGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "ipGroupValue")

read, err := client.Get(ctx, id, ipgroups.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IPGroupsClient.List`

```go
ctx := context.TODO()
id := ipgroups.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `IPGroupsClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := ipgroups.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `IPGroupsClient.UpdateGroups`

```go
ctx := context.TODO()
id := ipgroups.NewIPGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "ipGroupValue")

payload := ipgroups.TagsObject{
	// ...
}


read, err := client.UpdateGroups(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
