
## `github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-06-01/afdendpoints` Documentation

The `afdendpoints` SDK allows for interaction with Azure Resource Manager `cdn` (API Version `2025-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-06-01/afdendpoints"
```


### Client Initialization

```go
client := afdendpoints.NewAFDEndpointsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AFDEndpointsClient.Create`

```go
ctx := context.TODO()
id := afdendpoints.NewAfdEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "afdEndpointName")

payload := afdendpoints.AFDEndpoint{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AFDEndpointsClient.Delete`

```go
ctx := context.TODO()
id := afdendpoints.NewAfdEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "afdEndpointName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AFDEndpointsClient.Get`

```go
ctx := context.TODO()
id := afdendpoints.NewAfdEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "afdEndpointName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AFDEndpointsClient.ListByProfile`

```go
ctx := context.TODO()
id := afdendpoints.NewProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName")

// alternatively `client.ListByProfile(ctx, id)` can be used to do batched pagination
items, err := client.ListByProfileComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AFDEndpointsClient.ListResourceUsage`

```go
ctx := context.TODO()
id := afdendpoints.NewAfdEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "afdEndpointName")

// alternatively `client.ListResourceUsage(ctx, id)` can be used to do batched pagination
items, err := client.ListResourceUsageComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AFDEndpointsClient.PurgeContent`

```go
ctx := context.TODO()
id := afdendpoints.NewAfdEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "afdEndpointName")

payload := afdendpoints.AfdPurgeParameters{
	// ...
}


if err := client.PurgeContentThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AFDEndpointsClient.Update`

```go
ctx := context.TODO()
id := afdendpoints.NewAfdEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "afdEndpointName")

payload := afdendpoints.AFDEndpointUpdateParameters{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AFDEndpointsClient.ValidateCustomDomain`

```go
ctx := context.TODO()
id := afdendpoints.NewAfdEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "afdEndpointName")

payload := afdendpoints.ValidateCustomDomainInput{
	// ...
}


read, err := client.ValidateCustomDomain(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
