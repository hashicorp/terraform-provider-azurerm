
## `github.com/hashicorp/go-azure-sdk/resource-manager/digitaltwins/2023-01-31/digitaltwinsinstance` Documentation

The `digitaltwinsinstance` SDK allows for interaction with the Azure Resource Manager Service `digitaltwins` (API Version `2023-01-31`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/digitaltwins/2023-01-31/digitaltwinsinstance"
```


### Client Initialization

```go
client := digitaltwinsinstance.NewDigitalTwinsInstanceClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DigitalTwinsInstanceClient.DigitalTwinsCreateOrUpdate`

```go
ctx := context.TODO()
id := digitaltwinsinstance.NewDigitalTwinsInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "digitalTwinsInstanceValue")

payload := digitaltwinsinstance.DigitalTwinsDescription{
	// ...
}


if err := client.DigitalTwinsCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `DigitalTwinsInstanceClient.DigitalTwinsDelete`

```go
ctx := context.TODO()
id := digitaltwinsinstance.NewDigitalTwinsInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "digitalTwinsInstanceValue")

if err := client.DigitalTwinsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `DigitalTwinsInstanceClient.DigitalTwinsGet`

```go
ctx := context.TODO()
id := digitaltwinsinstance.NewDigitalTwinsInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "digitalTwinsInstanceValue")

read, err := client.DigitalTwinsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DigitalTwinsInstanceClient.DigitalTwinsList`

```go
ctx := context.TODO()
id := digitaltwinsinstance.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.DigitalTwinsList(ctx, id)` can be used to do batched pagination
items, err := client.DigitalTwinsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DigitalTwinsInstanceClient.DigitalTwinsListByResourceGroup`

```go
ctx := context.TODO()
id := digitaltwinsinstance.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.DigitalTwinsListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.DigitalTwinsListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DigitalTwinsInstanceClient.DigitalTwinsUpdate`

```go
ctx := context.TODO()
id := digitaltwinsinstance.NewDigitalTwinsInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "digitalTwinsInstanceValue")

payload := digitaltwinsinstance.DigitalTwinsPatchDescription{
	// ...
}


if err := client.DigitalTwinsUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
