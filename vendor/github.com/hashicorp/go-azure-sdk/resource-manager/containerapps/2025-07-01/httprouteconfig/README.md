
## `github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/httprouteconfig` Documentation

The `httprouteconfig` SDK allows for interaction with Azure Resource Manager `containerapps` (API Version `2025-07-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/httprouteconfig"
```


### Client Initialization

```go
client := httprouteconfig.NewHTTPRouteConfigClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `HTTPRouteConfigClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := httprouteconfig.NewHTTPRouteConfigID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedEnvironmentName", "httpRouteConfigName")

payload := httprouteconfig.HTTPRouteConfig{
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


### Example Usage: `HTTPRouteConfigClient.Delete`

```go
ctx := context.TODO()
id := httprouteconfig.NewHTTPRouteConfigID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedEnvironmentName", "httpRouteConfigName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `HTTPRouteConfigClient.Get`

```go
ctx := context.TODO()
id := httprouteconfig.NewHTTPRouteConfigID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedEnvironmentName", "httpRouteConfigName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `HTTPRouteConfigClient.List`

```go
ctx := context.TODO()
id := httprouteconfig.NewManagedEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedEnvironmentName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `HTTPRouteConfigClient.Update`

```go
ctx := context.TODO()
id := httprouteconfig.NewHTTPRouteConfigID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedEnvironmentName", "httpRouteConfigName")

payload := httprouteconfig.HTTPRouteConfig{
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
