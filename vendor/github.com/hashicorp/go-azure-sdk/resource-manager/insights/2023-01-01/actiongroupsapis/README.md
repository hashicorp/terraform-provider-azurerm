
## `github.com/hashicorp/go-azure-sdk/resource-manager/insights/2023-01-01/actiongroupsapis` Documentation

The `actiongroupsapis` SDK allows for interaction with Azure Resource Manager `insights` (API Version `2023-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/insights/2023-01-01/actiongroupsapis"
```


### Client Initialization

```go
client := actiongroupsapis.NewActionGroupsAPIsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ActionGroupsAPIsClient.ActionGroupsCreateNotificationsAtActionGroupResourceLevel`

```go
ctx := context.TODO()
id := actiongroupsapis.NewActionGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "actionGroupName")

payload := actiongroupsapis.NotificationRequestBody{
	// ...
}


if err := client.ActionGroupsCreateNotificationsAtActionGroupResourceLevelThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ActionGroupsAPIsClient.ActionGroupsCreateOrUpdate`

```go
ctx := context.TODO()
id := actiongroupsapis.NewActionGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "actionGroupName")

payload := actiongroupsapis.ActionGroupResource{
	// ...
}


read, err := client.ActionGroupsCreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ActionGroupsAPIsClient.ActionGroupsDelete`

```go
ctx := context.TODO()
id := actiongroupsapis.NewActionGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "actionGroupName")

read, err := client.ActionGroupsDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ActionGroupsAPIsClient.ActionGroupsEnableReceiver`

```go
ctx := context.TODO()
id := actiongroupsapis.NewActionGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "actionGroupName")

payload := actiongroupsapis.EnableRequest{
	// ...
}


read, err := client.ActionGroupsEnableReceiver(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ActionGroupsAPIsClient.ActionGroupsGet`

```go
ctx := context.TODO()
id := actiongroupsapis.NewActionGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "actionGroupName")

read, err := client.ActionGroupsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ActionGroupsAPIsClient.ActionGroupsGetTestNotificationsAtActionGroupResourceLevel`

```go
ctx := context.TODO()
id := actiongroupsapis.NewNotificationStatusID("12345678-1234-9876-4563-123456789012", "example-resource-group", "actionGroupName", "notificationId")

read, err := client.ActionGroupsGetTestNotificationsAtActionGroupResourceLevel(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ActionGroupsAPIsClient.ActionGroupsListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ActionGroupsListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ActionGroupsListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ActionGroupsAPIsClient.ActionGroupsListBySubscriptionId`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ActionGroupsListBySubscriptionId(ctx, id)` can be used to do batched pagination
items, err := client.ActionGroupsListBySubscriptionIdComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ActionGroupsAPIsClient.ActionGroupsUpdate`

```go
ctx := context.TODO()
id := actiongroupsapis.NewActionGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "actionGroupName")

payload := actiongroupsapis.ActionGroupPatchBody{
	// ...
}


read, err := client.ActionGroupsUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
