
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/notificationrecipientuser` Documentation

The `notificationrecipientuser` SDK allows for interaction with Azure Resource Manager `apimanagement` (API Version `2022-08-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/notificationrecipientuser"
```


### Client Initialization

```go
client := notificationrecipientuser.NewNotificationRecipientUserClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `NotificationRecipientUserClient.CheckEntityExists`

```go
ctx := context.TODO()
id := notificationrecipientuser.NewRecipientUserID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "AccountClosedPublisher", "userId")

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
id := notificationrecipientuser.NewRecipientUserID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "AccountClosedPublisher", "userId")

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
id := notificationrecipientuser.NewRecipientUserID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "AccountClosedPublisher", "userId")

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
id := notificationrecipientuser.NewNotificationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "AccountClosedPublisher")

// alternatively `client.ListByNotification(ctx, id)` can be used to do batched pagination
items, err := client.ListByNotificationComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
