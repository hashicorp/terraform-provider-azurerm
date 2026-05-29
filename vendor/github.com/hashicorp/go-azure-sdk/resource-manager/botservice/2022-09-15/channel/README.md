
## `github.com/hashicorp/go-azure-sdk/resource-manager/botservice/2022-09-15/channel` Documentation

The `channel` SDK allows for interaction with Azure Resource Manager `botservice` (API Version `2022-09-15`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/botservice/2022-09-15/channel"
```


### Client Initialization

```go
client := channel.NewChannelClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ChannelClient.Create`

```go
ctx := context.TODO()
id := commonids.NewBotServiceChannelID("12345678-1234-9876-4563-123456789012", "example-resource-group", "botServiceName", "AcsChatChannel")

payload := channel.BotChannel{
	// ...
}


read, err := client.Create(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ChannelClient.Delete`

```go
ctx := context.TODO()
id := commonids.NewBotServiceChannelID("12345678-1234-9876-4563-123456789012", "example-resource-group", "botServiceName", "AcsChatChannel")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ChannelClient.DirectLineRegenerateKeys`

```go
ctx := context.TODO()
id := commonids.NewBotServiceChannelID("12345678-1234-9876-4563-123456789012", "example-resource-group", "botServiceName", "AcsChatChannel")

payload := channel.SiteInfo{
	// ...
}


read, err := client.DirectLineRegenerateKeys(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ChannelClient.EmailCreateSignInURL`

```go
ctx := context.TODO()
id := commonids.NewBotServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "botServiceName")

read, err := client.EmailCreateSignInURL(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ChannelClient.Get`

```go
ctx := context.TODO()
id := commonids.NewBotServiceChannelID("12345678-1234-9876-4563-123456789012", "example-resource-group", "botServiceName", "AcsChatChannel")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ChannelClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewBotServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "botServiceName")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ChannelClient.ListWithKeys`

```go
ctx := context.TODO()
id := commonids.NewBotServiceChannelID("12345678-1234-9876-4563-123456789012", "example-resource-group", "botServiceName", "AcsChatChannel")

read, err := client.ListWithKeys(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ChannelClient.Update`

```go
ctx := context.TODO()
id := commonids.NewBotServiceChannelID("12345678-1234-9876-4563-123456789012", "example-resource-group", "botServiceName", "AcsChatChannel")

payload := channel.BotChannel{
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
