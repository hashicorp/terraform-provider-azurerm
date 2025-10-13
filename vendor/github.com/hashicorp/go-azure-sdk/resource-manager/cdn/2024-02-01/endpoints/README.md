
## `github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-02-01/endpoints` Documentation

The `endpoints` SDK allows for interaction with Azure Resource Manager `cdn` (API Version `2024-02-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-02-01/endpoints"
```


### Client Initialization

```go
client := endpoints.NewEndpointsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `EndpointsClient.Create`

```go
ctx := context.TODO()
id := endpoints.NewEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "endpointName")

payload := endpoints.Endpoint{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `EndpointsClient.Delete`

```go
ctx := context.TODO()
id := endpoints.NewEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "endpointName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `EndpointsClient.Get`

```go
ctx := context.TODO()
id := endpoints.NewEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "endpointName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EndpointsClient.ListByProfile`

```go
ctx := context.TODO()
id := endpoints.NewProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName")

// alternatively `client.ListByProfile(ctx, id)` can be used to do batched pagination
items, err := client.ListByProfileComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `EndpointsClient.ListResourceUsage`

```go
ctx := context.TODO()
id := endpoints.NewEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "endpointName")

// alternatively `client.ListResourceUsage(ctx, id)` can be used to do batched pagination
items, err := client.ListResourceUsageComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `EndpointsClient.LoadContent`

```go
ctx := context.TODO()
id := endpoints.NewEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "endpointName")

payload := endpoints.LoadParameters{
	// ...
}


if err := client.LoadContentThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `EndpointsClient.PurgeContent`

```go
ctx := context.TODO()
id := endpoints.NewEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "endpointName")

payload := endpoints.PurgeParameters{
	// ...
}


if err := client.PurgeContentThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `EndpointsClient.Start`

```go
ctx := context.TODO()
id := endpoints.NewEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "endpointName")

if err := client.StartThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `EndpointsClient.Stop`

```go
ctx := context.TODO()
id := endpoints.NewEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "endpointName")

if err := client.StopThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `EndpointsClient.Update`

```go
ctx := context.TODO()
id := endpoints.NewEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "endpointName")

payload := endpoints.EndpointUpdateParameters{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `EndpointsClient.ValidateCustomDomain`

```go
ctx := context.TODO()
id := endpoints.NewEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "endpointName")

payload := endpoints.ValidateCustomDomainInput{
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
