
## `github.com/hashicorp/go-azure-sdk/resource-manager/communication/2023-03-31/communicationservices` Documentation

The `communicationservices` SDK allows for interaction with the Azure Resource Manager Service `communication` (API Version `2023-03-31`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/communication/2023-03-31/communicationservices"
```


### Client Initialization

```go
client := communicationservices.NewCommunicationServicesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `CommunicationServicesClient.CheckNameAvailability`

```go
ctx := context.TODO()
id := communicationservices.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

payload := communicationservices.CheckNameAvailabilityRequest{
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


### Example Usage: `CommunicationServicesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := communicationservices.NewCommunicationServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "communicationServiceValue")

payload := communicationservices.CommunicationServiceResource{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CommunicationServicesClient.Delete`

```go
ctx := context.TODO()
id := communicationservices.NewCommunicationServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "communicationServiceValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CommunicationServicesClient.Get`

```go
ctx := context.TODO()
id := communicationservices.NewCommunicationServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "communicationServiceValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CommunicationServicesClient.LinkNotificationHub`

```go
ctx := context.TODO()
id := communicationservices.NewCommunicationServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "communicationServiceValue")

payload := communicationservices.LinkNotificationHubParameters{
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


### Example Usage: `CommunicationServicesClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := communicationservices.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `CommunicationServicesClient.ListBySubscription`

```go
ctx := context.TODO()
id := communicationservices.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `CommunicationServicesClient.ListKeys`

```go
ctx := context.TODO()
id := communicationservices.NewCommunicationServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "communicationServiceValue")

read, err := client.ListKeys(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CommunicationServicesClient.RegenerateKey`

```go
ctx := context.TODO()
id := communicationservices.NewCommunicationServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "communicationServiceValue")

payload := communicationservices.RegenerateKeyParameters{
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


### Example Usage: `CommunicationServicesClient.Update`

```go
ctx := context.TODO()
id := communicationservices.NewCommunicationServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "communicationServiceValue")

payload := communicationservices.CommunicationServiceResourceUpdate{
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
