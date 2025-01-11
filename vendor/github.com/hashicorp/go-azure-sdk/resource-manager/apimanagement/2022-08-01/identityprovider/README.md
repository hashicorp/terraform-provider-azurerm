
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/identityprovider` Documentation

The `identityprovider` SDK allows for interaction with Azure Resource Manager `apimanagement` (API Version `2022-08-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/identityprovider"
```


### Client Initialization

```go
client := identityprovider.NewIdentityProviderClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `IdentityProviderClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := identityprovider.NewIdentityProviderID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "aad")

payload := identityprovider.IdentityProviderCreateContract{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload, identityprovider.DefaultCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IdentityProviderClient.Delete`

```go
ctx := context.TODO()
id := identityprovider.NewIdentityProviderID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "aad")

read, err := client.Delete(ctx, id, identityprovider.DefaultDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IdentityProviderClient.Get`

```go
ctx := context.TODO()
id := identityprovider.NewIdentityProviderID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "aad")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IdentityProviderClient.GetEntityTag`

```go
ctx := context.TODO()
id := identityprovider.NewIdentityProviderID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "aad")

read, err := client.GetEntityTag(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IdentityProviderClient.ListByService`

```go
ctx := context.TODO()
id := identityprovider.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName")

// alternatively `client.ListByService(ctx, id)` can be used to do batched pagination
items, err := client.ListByServiceComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `IdentityProviderClient.ListSecrets`

```go
ctx := context.TODO()
id := identityprovider.NewIdentityProviderID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "aad")

read, err := client.ListSecrets(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IdentityProviderClient.Update`

```go
ctx := context.TODO()
id := identityprovider.NewIdentityProviderID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "aad")

payload := identityprovider.IdentityProviderUpdateParameters{
	// ...
}


read, err := client.Update(ctx, id, payload, identityprovider.DefaultUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
