
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2023-05-01-preview/notificationrecipientuser` Documentation

The `notificationrecipientuser` SDK allows for interaction with the Azure Resource Manager Service `apimanagement` (API Version `2023-05-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2023-05-01-preview/notificationrecipientuser"
```


### Client Initialization

```go
client := notificationrecipientuser.NewNotificationRecipientUserClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `NotificationRecipientUserClient.CheckEntityExists`

```go
ctx := context.TODO()
id := notificationrecipientuser.NewRecipientUserID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "AccountClosedPublisher", "userIdValue")

read, err := client.CheckEntityExists(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NotificationRecipientUserClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := notificationrecipientuser.NewRecipientUserID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "AccountClosedPublisher", "userIdValue")

read, err := client.CreateOrUpdate(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NotificationRecipientUserClient.Delete`

```go
ctx := context.TODO()
id := notificationrecipientuser.NewRecipientUserID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "AccountClosedPublisher", "userIdValue")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NotificationRecipientUserClient.ListByNotification`

```go
ctx := context.TODO()
id := notificationrecipientuser.NewNotificationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "AccountClosedPublisher")

// alternatively `client.ListByNotification(ctx, id)` can be used to do batched pagination
items, err := client.ListByNotificationComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `NotificationRecipientUserClient.WorkspaceNotificationRecipientUserCheckEntityExists`

```go
ctx := context.TODO()
id := notificationrecipientuser.NewNotificationRecipientUserID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "workspaceIdValue", "AccountClosedPublisher", "userIdValue")

read, err := client.WorkspaceNotificationRecipientUserCheckEntityExists(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NotificationRecipientUserClient.WorkspaceNotificationRecipientUserCreateOrUpdate`

```go
ctx := context.TODO()
id := notificationrecipientuser.NewNotificationRecipientUserID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "workspaceIdValue", "AccountClosedPublisher", "userIdValue")

read, err := client.WorkspaceNotificationRecipientUserCreateOrUpdate(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NotificationRecipientUserClient.WorkspaceNotificationRecipientUserDelete`

```go
ctx := context.TODO()
id := notificationrecipientuser.NewNotificationRecipientUserID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "workspaceIdValue", "AccountClosedPublisher", "userIdValue")

read, err := client.WorkspaceNotificationRecipientUserDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NotificationRecipientUserClient.WorkspaceNotificationRecipientUserListByNotification`

```go
ctx := context.TODO()
id := notificationrecipientuser.NewNotificationNotificationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "workspaceIdValue", "AccountClosedPublisher")

// alternatively `client.WorkspaceNotificationRecipientUserListByNotification(ctx, id)` can be used to do batched pagination
items, err := client.WorkspaceNotificationRecipientUserListByNotificationComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
