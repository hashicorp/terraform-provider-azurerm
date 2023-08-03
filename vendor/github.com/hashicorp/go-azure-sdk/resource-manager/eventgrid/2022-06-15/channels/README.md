
## `github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2022-06-15/channels` Documentation

The `channels` SDK allows for interaction with the Azure Resource Manager Service `eventgrid` (API Version `2022-06-15`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2022-06-15/channels"
```


### Client Initialization

```go
client := channels.NewChannelsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ChannelsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := channels.NewChannelID("12345678-1234-9876-4563-123456789012", "example-resource-group", "partnerNamespaceValue", "channelValue")

payload := channels.Channel{
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


### Example Usage: `ChannelsClient.Delete`

```go
ctx := context.TODO()
id := channels.NewChannelID("12345678-1234-9876-4563-123456789012", "example-resource-group", "partnerNamespaceValue", "channelValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ChannelsClient.Get`

```go
ctx := context.TODO()
id := channels.NewChannelID("12345678-1234-9876-4563-123456789012", "example-resource-group", "partnerNamespaceValue", "channelValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ChannelsClient.GetFullUrl`

```go
ctx := context.TODO()
id := channels.NewChannelID("12345678-1234-9876-4563-123456789012", "example-resource-group", "partnerNamespaceValue", "channelValue")

read, err := client.GetFullUrl(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ChannelsClient.ListByPartnerNamespace`

```go
ctx := context.TODO()
id := channels.NewPartnerNamespaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "partnerNamespaceValue")

// alternatively `client.ListByPartnerNamespace(ctx, id, channels.DefaultListByPartnerNamespaceOperationOptions())` can be used to do batched pagination
items, err := client.ListByPartnerNamespaceComplete(ctx, id, channels.DefaultListByPartnerNamespaceOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ChannelsClient.Update`

```go
ctx := context.TODO()
id := channels.NewChannelID("12345678-1234-9876-4563-123456789012", "example-resource-group", "partnerNamespaceValue", "channelValue")

payload := channels.ChannelUpdateParameters{
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
