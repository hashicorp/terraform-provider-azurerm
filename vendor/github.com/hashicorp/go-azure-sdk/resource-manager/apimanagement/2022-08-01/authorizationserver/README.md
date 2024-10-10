
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/authorizationserver` Documentation

The `authorizationserver` SDK allows for interaction with Azure Resource Manager `apimanagement` (API Version `2022-08-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/authorizationserver"
```


### Client Initialization

```go
client := authorizationserver.NewAuthorizationServerClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AuthorizationServerClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := authorizationserver.NewAuthorizationServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "authorizationServerName")

payload := authorizationserver.AuthorizationServerContract{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload, authorizationserver.DefaultCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AuthorizationServerClient.Delete`

```go
ctx := context.TODO()
id := authorizationserver.NewAuthorizationServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "authorizationServerName")

read, err := client.Delete(ctx, id, authorizationserver.DefaultDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AuthorizationServerClient.Get`

```go
ctx := context.TODO()
id := authorizationserver.NewAuthorizationServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "authorizationServerName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AuthorizationServerClient.GetEntityTag`

```go
ctx := context.TODO()
id := authorizationserver.NewAuthorizationServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "authorizationServerName")

read, err := client.GetEntityTag(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AuthorizationServerClient.ListByService`

```go
ctx := context.TODO()
id := authorizationserver.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName")

// alternatively `client.ListByService(ctx, id, authorizationserver.DefaultListByServiceOperationOptions())` can be used to do batched pagination
items, err := client.ListByServiceComplete(ctx, id, authorizationserver.DefaultListByServiceOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AuthorizationServerClient.ListSecrets`

```go
ctx := context.TODO()
id := authorizationserver.NewAuthorizationServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "authorizationServerName")

read, err := client.ListSecrets(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AuthorizationServerClient.Update`

```go
ctx := context.TODO()
id := authorizationserver.NewAuthorizationServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "authorizationServerName")

payload := authorizationserver.AuthorizationServerUpdateContract{
	// ...
}


read, err := client.Update(ctx, id, payload, authorizationserver.DefaultUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
