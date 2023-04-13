
## `github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2021-08-01-preview/webhooks` Documentation

The `webhooks` SDK allows for interaction with the Azure Resource Manager Service `containerregistry` (API Version `2021-08-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2021-08-01-preview/webhooks"
```


### Client Initialization

```go
client := webhooks.NewWebHooksClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `WebHooksClient.Create`

```go
ctx := context.TODO()
id := webhooks.NewWebHookID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryValue", "webHookValue")

payload := webhooks.WebhookCreateParameters{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `WebHooksClient.Delete`

```go
ctx := context.TODO()
id := webhooks.NewWebHookID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryValue", "webHookValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `WebHooksClient.Get`

```go
ctx := context.TODO()
id := webhooks.NewWebHookID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryValue", "webHookValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebHooksClient.GetCallbackConfig`

```go
ctx := context.TODO()
id := webhooks.NewWebHookID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryValue", "webHookValue")

read, err := client.GetCallbackConfig(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebHooksClient.List`

```go
ctx := context.TODO()
id := webhooks.NewRegistryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryValue")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebHooksClient.ListEvents`

```go
ctx := context.TODO()
id := webhooks.NewWebHookID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryValue", "webHookValue")

// alternatively `client.ListEvents(ctx, id)` can be used to do batched pagination
items, err := client.ListEventsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebHooksClient.Ping`

```go
ctx := context.TODO()
id := webhooks.NewWebHookID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryValue", "webHookValue")

read, err := client.Ping(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebHooksClient.Update`

```go
ctx := context.TODO()
id := webhooks.NewWebHookID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryValue", "webHookValue")

payload := webhooks.WebhookUpdateParameters{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
