
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2026-03-01/expressroutelagauthorizations` Documentation

The `expressroutelagauthorizations` SDK allows for interaction with Azure Resource Manager `network` (API Version `2026-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2026-03-01/expressroutelagauthorizations"
```


### Client Initialization

```go
client := expressroutelagauthorizations.NewExpressRouteLagAuthorizationsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ExpressRouteLagAuthorizationsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := expressroutelagauthorizations.NewExpressRouteLagAuthorizationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "expressRouteLagName", "authorizationName")

payload := expressroutelagauthorizations.ExpressRouteLagAuthorization{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ExpressRouteLagAuthorizationsClient.Delete`

```go
ctx := context.TODO()
id := expressroutelagauthorizations.NewExpressRouteLagAuthorizationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "expressRouteLagName", "authorizationName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ExpressRouteLagAuthorizationsClient.Get`

```go
ctx := context.TODO()
id := expressroutelagauthorizations.NewExpressRouteLagAuthorizationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "expressRouteLagName", "authorizationName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ExpressRouteLagAuthorizationsClient.List`

```go
ctx := context.TODO()
id := expressroutelagauthorizations.NewExpressRouteLagID("12345678-1234-9876-4563-123456789012", "example-resource-group", "expressRouteLagName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ExpressRouteLagAuthorizationsClient.ListKeys`

```go
ctx := context.TODO()
id := expressroutelagauthorizations.NewExpressRouteLagAuthorizationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "expressRouteLagName", "authorizationName")

read, err := client.ListKeys(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
