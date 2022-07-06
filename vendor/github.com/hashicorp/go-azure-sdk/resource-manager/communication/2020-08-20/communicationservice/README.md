
## `github.com/hashicorp/go-azure-sdk/resource-manager/communication/2020-08-20/communicationservice` Documentation

The `communicationservice` SDK allows for interaction with the Azure Resource Manager Service `communication` (API Version `2020-08-20`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/communication/2020-08-20/communicationservice"
```


### Client Initialization

```go
client := communicationservice.NewCommunicationServiceClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `CommunicationServiceClient.CheckNameAvailability`

```go
ctx := context.TODO()
id := communicationservice.NewSubscriptionID()

payload := communicationservice.NameAvailabilityParameters{
	// ...
}


read, err := client.CheckNameAvailability(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CommunicationServiceClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := communicationservice.NewCommunicationServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "communicationServiceValue")

payload := communicationservice.CommunicationServiceResource{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CommunicationServiceClient.Delete`

```go
ctx := context.TODO()
id := communicationservice.NewCommunicationServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "communicationServiceValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CommunicationServiceClient.Get`

```go
ctx := context.TODO()
id := communicationservice.NewCommunicationServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "communicationServiceValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CommunicationServiceClient.LinkNotificationHub`

```go
ctx := context.TODO()
id := communicationservice.NewCommunicationServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "communicationServiceValue")

payload := communicationservice.LinkNotificationHubParameters{
	// ...
}


read, err := client.LinkNotificationHub(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CommunicationServiceClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := communicationservice.NewResourceGroupID()

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `CommunicationServiceClient.ListBySubscription`

```go
ctx := context.TODO()
id := communicationservice.NewSubscriptionID()

// alternatively `client.ListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `CommunicationServiceClient.ListKeys`

```go
ctx := context.TODO()
id := communicationservice.NewCommunicationServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "communicationServiceValue")

read, err := client.ListKeys(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CommunicationServiceClient.RegenerateKey`

```go
ctx := context.TODO()
id := communicationservice.NewCommunicationServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "communicationServiceValue")

payload := communicationservice.RegenerateKeyParameters{
	// ...
}


read, err := client.RegenerateKey(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CommunicationServiceClient.Update`

```go
ctx := context.TODO()
id := communicationservice.NewCommunicationServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "communicationServiceValue")

payload := communicationservice.CommunicationServiceResource{
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
