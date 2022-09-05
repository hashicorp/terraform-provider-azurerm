
## `github.com/hashicorp/go-azure-sdk/resource-manager/notificationhubs/2017-04-01/notificationhubs` Documentation

The `notificationhubs` SDK allows for interaction with the Azure Resource Manager Service `notificationhubs` (API Version `2017-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/notificationhubs/2017-04-01/notificationhubs"
```


### Client Initialization

```go
client := notificationhubs.NewNotificationHubsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `NotificationHubsClient.CheckNotificationHubAvailability`

```go
ctx := context.TODO()
id := notificationhubs.NewNamespaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceValue")

payload := notificationhubs.CheckAvailabilityParameters{
	// ...
}


read, err := client.CheckNotificationHubAvailability(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NotificationHubsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := notificationhubs.NewNotificationHubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceValue", "notificationHubValue")

payload := notificationhubs.NotificationHubCreateOrUpdateParameters{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NotificationHubsClient.CreateOrUpdateAuthorizationRule`

```go
ctx := context.TODO()
id := notificationhubs.NewNotificationHubAuthorizationRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceValue", "notificationHubValue", "authorizationRuleValue")

payload := notificationhubs.SharedAccessAuthorizationRuleCreateOrUpdateParameters{
	// ...
}


read, err := client.CreateOrUpdateAuthorizationRule(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NotificationHubsClient.DebugSend`

```go
ctx := context.TODO()
id := notificationhubs.NewNotificationHubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceValue", "notificationHubValue")
var payload interface{}

read, err := client.DebugSend(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NotificationHubsClient.Delete`

```go
ctx := context.TODO()
id := notificationhubs.NewNotificationHubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceValue", "notificationHubValue")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NotificationHubsClient.DeleteAuthorizationRule`

```go
ctx := context.TODO()
id := notificationhubs.NewNotificationHubAuthorizationRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceValue", "notificationHubValue", "authorizationRuleValue")

read, err := client.DeleteAuthorizationRule(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NotificationHubsClient.Get`

```go
ctx := context.TODO()
id := notificationhubs.NewNotificationHubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceValue", "notificationHubValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NotificationHubsClient.GetAuthorizationRule`

```go
ctx := context.TODO()
id := notificationhubs.NewNotificationHubAuthorizationRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceValue", "notificationHubValue", "authorizationRuleValue")

read, err := client.GetAuthorizationRule(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NotificationHubsClient.GetPnsCredentials`

```go
ctx := context.TODO()
id := notificationhubs.NewNotificationHubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceValue", "notificationHubValue")

read, err := client.GetPnsCredentials(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NotificationHubsClient.List`

```go
ctx := context.TODO()
id := notificationhubs.NewNamespaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceValue")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `NotificationHubsClient.ListAuthorizationRules`

```go
ctx := context.TODO()
id := notificationhubs.NewNotificationHubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceValue", "notificationHubValue")

// alternatively `client.ListAuthorizationRules(ctx, id)` can be used to do batched pagination
items, err := client.ListAuthorizationRulesComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `NotificationHubsClient.ListKeys`

```go
ctx := context.TODO()
id := notificationhubs.NewNotificationHubAuthorizationRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceValue", "notificationHubValue", "authorizationRuleValue")

read, err := client.ListKeys(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NotificationHubsClient.Patch`

```go
ctx := context.TODO()
id := notificationhubs.NewNotificationHubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceValue", "notificationHubValue")

payload := notificationhubs.NotificationHubPatchParameters{
	// ...
}


read, err := client.Patch(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NotificationHubsClient.RegenerateKeys`

```go
ctx := context.TODO()
id := notificationhubs.NewNotificationHubAuthorizationRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceValue", "notificationHubValue", "authorizationRuleValue")

payload := notificationhubs.PolicykeyResource{
	// ...
}


read, err := client.RegenerateKeys(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
