
## `github.com/hashicorp/go-azure-sdk/resource-manager/media/2021-11-01/liveevents` Documentation

The `liveevents` SDK allows for interaction with the Azure Resource Manager Service `media` (API Version `2021-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/media/2021-11-01/liveevents"
```


### Client Initialization

```go
client := liveevents.NewLiveEventsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `LiveEventsClient.Allocate`

```go
ctx := context.TODO()
id := liveevents.NewLiveEventID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "liveEventValue")

if err := client.AllocateThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `LiveEventsClient.Create`

```go
ctx := context.TODO()
id := liveevents.NewLiveEventID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "liveEventValue")

payload := liveevents.LiveEvent{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload, liveevents.DefaultCreateOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `LiveEventsClient.Delete`

```go
ctx := context.TODO()
id := liveevents.NewLiveEventID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "liveEventValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `LiveEventsClient.Get`

```go
ctx := context.TODO()
id := liveevents.NewLiveEventID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "liveEventValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `LiveEventsClient.List`

```go
ctx := context.TODO()
id := liveevents.NewMediaServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `LiveEventsClient.Reset`

```go
ctx := context.TODO()
id := liveevents.NewLiveEventID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "liveEventValue")

if err := client.ResetThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `LiveEventsClient.Start`

```go
ctx := context.TODO()
id := liveevents.NewLiveEventID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "liveEventValue")

if err := client.StartThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `LiveEventsClient.Stop`

```go
ctx := context.TODO()
id := liveevents.NewLiveEventID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "liveEventValue")

payload := liveevents.LiveEventActionInput{
	// ...
}


if err := client.StopThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `LiveEventsClient.Update`

```go
ctx := context.TODO()
id := liveevents.NewLiveEventID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "liveEventValue")

payload := liveevents.LiveEvent{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
