
## `github.com/hashicorp/go-azure-sdk/resource-manager/healthbot/2025-05-25/healthbots` Documentation

The `healthbots` SDK allows for interaction with Azure Resource Manager `healthbot` (API Version `2025-05-25`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/healthbot/2025-05-25/healthbots"
```


### Client Initialization

```go
client := healthbots.NewHealthBotsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `HealthBotsClient.BotsCreate`

```go
ctx := context.TODO()
id := healthbots.NewHealthBotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "healthBotName")

payload := healthbots.HealthBot{
	// ...
}


if err := client.BotsCreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `HealthBotsClient.BotsDelete`

```go
ctx := context.TODO()
id := healthbots.NewHealthBotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "healthBotName")

if err := client.BotsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `HealthBotsClient.BotsGet`

```go
ctx := context.TODO()
id := healthbots.NewHealthBotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "healthBotName")

read, err := client.BotsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `HealthBotsClient.BotsList`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.BotsList(ctx, id)` can be used to do batched pagination
items, err := client.BotsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `HealthBotsClient.BotsListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.BotsListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.BotsListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `HealthBotsClient.BotsListSecrets`

```go
ctx := context.TODO()
id := healthbots.NewHealthBotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "healthBotName")

read, err := client.BotsListSecrets(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `HealthBotsClient.BotsRegenerateApiJwtSecret`

```go
ctx := context.TODO()
id := healthbots.NewHealthBotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "healthBotName")

read, err := client.BotsRegenerateApiJwtSecret(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `HealthBotsClient.BotsUpdate`

```go
ctx := context.TODO()
id := healthbots.NewHealthBotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "healthBotName")

payload := healthbots.HealthBotUpdateParameters{
	// ...
}


if err := client.BotsUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
