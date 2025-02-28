
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/notification` Documentation

The `notification` SDK allows for interaction with Azure Resource Manager `apimanagement` (API Version `2024-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/notification"
```


### Client Initialization

```go
client := notification.NewNotificationClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `NotificationClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := notification.NewNotificationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "AccountClosedPublisher")

read, err := client.CreateOrUpdate(ctx, id, notification.DefaultCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NotificationClient.Get`

```go
ctx := context.TODO()
id := notification.NewNotificationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "AccountClosedPublisher")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NotificationClient.ListByService`

```go
ctx := context.TODO()
id := notification.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName")

// alternatively `client.ListByService(ctx, id, notification.DefaultListByServiceOperationOptions())` can be used to do batched pagination
items, err := client.ListByServiceComplete(ctx, id, notification.DefaultListByServiceOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `NotificationClient.WorkspaceNotificationCreateOrUpdate`

```go
ctx := context.TODO()
id := notification.NewNotificationNotificationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceId", "AccountClosedPublisher")

read, err := client.WorkspaceNotificationCreateOrUpdate(ctx, id, notification.DefaultWorkspaceNotificationCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NotificationClient.WorkspaceNotificationGet`

```go
ctx := context.TODO()
id := notification.NewNotificationNotificationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceId", "AccountClosedPublisher")

read, err := client.WorkspaceNotificationGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NotificationClient.WorkspaceNotificationListByService`

```go
ctx := context.TODO()
id := notification.NewWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceId")

// alternatively `client.WorkspaceNotificationListByService(ctx, id, notification.DefaultWorkspaceNotificationListByServiceOperationOptions())` can be used to do batched pagination
items, err := client.WorkspaceNotificationListByServiceComplete(ctx, id, notification.DefaultWorkspaceNotificationListByServiceOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
