
## `github.com/hashicorp/go-azure-sdk/resource-manager/automation/2015-10-31/webhook` Documentation

The `webhook` SDK allows for interaction with Azure Resource Manager `automation` (API Version `2015-10-31`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/automation/2015-10-31/webhook"
```


### Client Initialization

```go
client := webhook.NewWebhookClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `WebhookClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := webhook.NewWebHookID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "webHookName")

payload := webhook.WebhookCreateOrUpdateParameters{
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


### Example Usage: `WebhookClient.Delete`

```go
ctx := context.TODO()
id := webhook.NewWebHookID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "webHookName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebhookClient.GenerateUri`

```go
ctx := context.TODO()
id := webhook.NewAutomationAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName")

read, err := client.GenerateUri(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebhookClient.Get`

```go
ctx := context.TODO()
id := webhook.NewWebHookID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "webHookName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebhookClient.ListByAutomationAccount`

```go
ctx := context.TODO()
id := webhook.NewAutomationAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName")

// alternatively `client.ListByAutomationAccount(ctx, id, webhook.DefaultListByAutomationAccountOperationOptions())` can be used to do batched pagination
items, err := client.ListByAutomationAccountComplete(ctx, id, webhook.DefaultListByAutomationAccountOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebhookClient.Update`

```go
ctx := context.TODO()
id := webhook.NewWebHookID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "webHookName")

payload := webhook.WebhookUpdateParameters{
	// ...
}


read, err := client.Update(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
