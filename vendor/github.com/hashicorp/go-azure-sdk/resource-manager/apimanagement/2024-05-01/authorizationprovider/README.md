
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/authorizationprovider` Documentation

The `authorizationprovider` SDK allows for interaction with Azure Resource Manager `apimanagement` (API Version `2024-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/authorizationprovider"
```


### Client Initialization

```go
client := authorizationprovider.NewAuthorizationProviderClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AuthorizationProviderClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := authorizationprovider.NewAuthorizationProviderID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "authorizationProviderId")

payload := authorizationprovider.AuthorizationProviderContract{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload, authorizationprovider.DefaultCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AuthorizationProviderClient.Delete`

```go
ctx := context.TODO()
id := authorizationprovider.NewAuthorizationProviderID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "authorizationProviderId")

read, err := client.Delete(ctx, id, authorizationprovider.DefaultDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AuthorizationProviderClient.Get`

```go
ctx := context.TODO()
id := authorizationprovider.NewAuthorizationProviderID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "authorizationProviderId")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AuthorizationProviderClient.ListByService`

```go
ctx := context.TODO()
id := authorizationprovider.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName")

// alternatively `client.ListByService(ctx, id, authorizationprovider.DefaultListByServiceOperationOptions())` can be used to do batched pagination
items, err := client.ListByServiceComplete(ctx, id, authorizationprovider.DefaultListByServiceOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
