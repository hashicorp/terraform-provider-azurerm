
## `github.com/hashicorp/go-azure-sdk/resource-manager/trafficmanager/2018-08-01/profiles` Documentation

The `profiles` SDK allows for interaction with the Azure Resource Manager Service `trafficmanager` (API Version `2018-08-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/trafficmanager/2018-08-01/profiles"
```


### Client Initialization

```go
client := profiles.NewProfilesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ProfilesClient.CheckTrafficManagerRelativeDnsNameAvailability`

```go
ctx := context.TODO()

payload := profiles.CheckTrafficManagerRelativeDnsNameAvailabilityParameters{
	// ...
}


read, err := client.CheckTrafficManagerRelativeDnsNameAvailability(ctx, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProfilesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := profiles.NewTrafficManagerProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "trafficManagerProfileValue")

payload := profiles.Profile{
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


### Example Usage: `ProfilesClient.Delete`

```go
ctx := context.TODO()
id := profiles.NewTrafficManagerProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "trafficManagerProfileValue")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProfilesClient.Get`

```go
ctx := context.TODO()
id := profiles.NewTrafficManagerProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "trafficManagerProfileValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProfilesClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := profiles.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

read, err := client.ListByResourceGroup(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProfilesClient.ListBySubscription`

```go
ctx := context.TODO()
id := profiles.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

read, err := client.ListBySubscription(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProfilesClient.Update`

```go
ctx := context.TODO()
id := profiles.NewTrafficManagerProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "trafficManagerProfileValue")

payload := profiles.Profile{
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
