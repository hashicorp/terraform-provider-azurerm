
## `github.com/hashicorp/go-azure-sdk/resource-manager/vmware/2022-05-01/authorizations` Documentation

The `authorizations` SDK allows for interaction with the Azure Resource Manager Service `vmware` (API Version `2022-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/vmware/2022-05-01/authorizations"
```


### Client Initialization

```go
client := authorizations.NewAuthorizationsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AuthorizationsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := authorizations.NewAuthorizationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateCloudValue", "authorizationValue")

payload := authorizations.ExpressRouteAuthorization{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AuthorizationsClient.Delete`

```go
ctx := context.TODO()
id := authorizations.NewAuthorizationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateCloudValue", "authorizationValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AuthorizationsClient.Get`

```go
ctx := context.TODO()
id := authorizations.NewAuthorizationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateCloudValue", "authorizationValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AuthorizationsClient.List`

```go
ctx := context.TODO()
id := authorizations.NewPrivateCloudID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateCloudValue")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
