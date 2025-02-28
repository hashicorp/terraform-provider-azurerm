
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/authorization` Documentation

The `authorization` SDK allows for interaction with Azure Resource Manager `apimanagement` (API Version `2024-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/authorization"
```


### Client Initialization

```go
client := authorization.NewAuthorizationClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AuthorizationClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := authorization.NewAuthorizationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "authorizationProviderId", "authorizationId")

payload := authorization.AuthorizationContract{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload, authorization.DefaultCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AuthorizationClient.Delete`

```go
ctx := context.TODO()
id := authorization.NewAuthorizationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "authorizationProviderId", "authorizationId")

read, err := client.Delete(ctx, id, authorization.DefaultDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AuthorizationClient.Get`

```go
ctx := context.TODO()
id := authorization.NewAuthorizationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "authorizationProviderId", "authorizationId")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
