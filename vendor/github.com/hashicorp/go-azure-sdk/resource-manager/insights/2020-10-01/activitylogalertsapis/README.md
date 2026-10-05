
## `github.com/hashicorp/go-azure-sdk/resource-manager/insights/2020-10-01/activitylogalertsapis` Documentation

The `activitylogalertsapis` SDK allows for interaction with Azure Resource Manager `insights` (API Version `2020-10-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/insights/2020-10-01/activitylogalertsapis"
```


### Client Initialization

```go
client := activitylogalertsapis.NewActivityLogAlertsAPIsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ActivityLogAlertsAPIsClient.ActivityLogAlertsCreateOrUpdate`

```go
ctx := context.TODO()
id := activitylogalertsapis.NewActivityLogAlertID("12345678-1234-9876-4563-123456789012", "example-resource-group", "activityLogAlertName")

payload := activitylogalertsapis.ActivityLogAlertResource{
	// ...
}


read, err := client.ActivityLogAlertsCreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ActivityLogAlertsAPIsClient.ActivityLogAlertsDelete`

```go
ctx := context.TODO()
id := activitylogalertsapis.NewActivityLogAlertID("12345678-1234-9876-4563-123456789012", "example-resource-group", "activityLogAlertName")

read, err := client.ActivityLogAlertsDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ActivityLogAlertsAPIsClient.ActivityLogAlertsGet`

```go
ctx := context.TODO()
id := activitylogalertsapis.NewActivityLogAlertID("12345678-1234-9876-4563-123456789012", "example-resource-group", "activityLogAlertName")

read, err := client.ActivityLogAlertsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ActivityLogAlertsAPIsClient.ActivityLogAlertsListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ActivityLogAlertsListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ActivityLogAlertsListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ActivityLogAlertsAPIsClient.ActivityLogAlertsListBySubscriptionId`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ActivityLogAlertsListBySubscriptionId(ctx, id)` can be used to do batched pagination
items, err := client.ActivityLogAlertsListBySubscriptionIdComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ActivityLogAlertsAPIsClient.ActivityLogAlertsUpdate`

```go
ctx := context.TODO()
id := activitylogalertsapis.NewActivityLogAlertID("12345678-1234-9876-4563-123456789012", "example-resource-group", "activityLogAlertName")

payload := activitylogalertsapis.AlertRulePatchObject{
	// ...
}


read, err := client.ActivityLogAlertsUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
