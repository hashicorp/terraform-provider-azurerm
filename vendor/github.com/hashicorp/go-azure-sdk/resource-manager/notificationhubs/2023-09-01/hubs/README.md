
## `github.com/hashicorp/go-azure-sdk/resource-manager/notificationhubs/2023-09-01/hubs` Documentation

The `hubs` SDK allows for interaction with Azure Resource Manager `notificationhubs` (API Version `2023-09-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/notificationhubs/2023-09-01/hubs"
```


### Client Initialization

```go
client := hubs.NewHubsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `HubsClient.NotificationHubsCheckNotificationHubAvailability`

```go
ctx := context.TODO()
id := hubs.NewNamespaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName")

payload := hubs.CheckAvailabilityParameters{
	// ...
}


read, err := client.NotificationHubsCheckNotificationHubAvailability(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `HubsClient.NotificationHubsCreateOrUpdate`

```go
ctx := context.TODO()
id := hubs.NewNotificationHubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "notificationHubName")

payload := hubs.NotificationHubResource{
	// ...
}


read, err := client.NotificationHubsCreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `HubsClient.NotificationHubsCreateOrUpdateAuthorizationRule`

```go
ctx := context.TODO()
id := hubs.NewNotificationHubAuthorizationRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "notificationHubName", "authorizationRuleName")

payload := hubs.SharedAccessAuthorizationRuleResource{
	// ...
}


read, err := client.NotificationHubsCreateOrUpdateAuthorizationRule(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `HubsClient.NotificationHubsDebugSend`

```go
ctx := context.TODO()
id := hubs.NewNotificationHubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "notificationHubName")

read, err := client.NotificationHubsDebugSend(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `HubsClient.NotificationHubsDelete`

```go
ctx := context.TODO()
id := hubs.NewNotificationHubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "notificationHubName")

read, err := client.NotificationHubsDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `HubsClient.NotificationHubsDeleteAuthorizationRule`

```go
ctx := context.TODO()
id := hubs.NewNotificationHubAuthorizationRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "notificationHubName", "authorizationRuleName")

read, err := client.NotificationHubsDeleteAuthorizationRule(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `HubsClient.NotificationHubsGet`

```go
ctx := context.TODO()
id := hubs.NewNotificationHubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "notificationHubName")

read, err := client.NotificationHubsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `HubsClient.NotificationHubsGetAuthorizationRule`

```go
ctx := context.TODO()
id := hubs.NewNotificationHubAuthorizationRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "notificationHubName", "authorizationRuleName")

read, err := client.NotificationHubsGetAuthorizationRule(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `HubsClient.NotificationHubsGetPnsCredentials`

```go
ctx := context.TODO()
id := hubs.NewNotificationHubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "notificationHubName")

read, err := client.NotificationHubsGetPnsCredentials(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `HubsClient.NotificationHubsList`

```go
ctx := context.TODO()
id := hubs.NewNamespaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName")

// alternatively `client.NotificationHubsList(ctx, id, hubs.DefaultNotificationHubsListOperationOptions())` can be used to do batched pagination
items, err := client.NotificationHubsListComplete(ctx, id, hubs.DefaultNotificationHubsListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `HubsClient.NotificationHubsListAuthorizationRules`

```go
ctx := context.TODO()
id := hubs.NewNotificationHubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "notificationHubName")

// alternatively `client.NotificationHubsListAuthorizationRules(ctx, id)` can be used to do batched pagination
items, err := client.NotificationHubsListAuthorizationRulesComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `HubsClient.NotificationHubsListKeys`

```go
ctx := context.TODO()
id := hubs.NewNotificationHubAuthorizationRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "notificationHubName", "authorizationRuleName")

read, err := client.NotificationHubsListKeys(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `HubsClient.NotificationHubsRegenerateKeys`

```go
ctx := context.TODO()
id := hubs.NewNotificationHubAuthorizationRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "notificationHubName", "authorizationRuleName")

payload := hubs.PolicyKeyResource{
	// ...
}


read, err := client.NotificationHubsRegenerateKeys(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `HubsClient.NotificationHubsUpdate`

```go
ctx := context.TODO()
id := hubs.NewNotificationHubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "notificationHubName")

payload := hubs.NotificationHubPatchParameters{
	// ...
}


read, err := client.NotificationHubsUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
