
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/tenantaccess` Documentation

The `tenantaccess` SDK allows for interaction with the Azure Resource Manager Service `apimanagement` (API Version `2021-08-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/tenantaccess"
```


### Client Initialization

```go
client := tenantaccess.NewTenantAccessClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `TenantAccessClient.Create`

```go
ctx := context.TODO()
id := tenantaccess.NewAccessID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "access")

payload := tenantaccess.AccessInformationCreateParameters{
	// ...
}


read, err := client.Create(ctx, id, payload, tenantaccess.DefaultCreateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TenantAccessClient.Get`

```go
ctx := context.TODO()
id := tenantaccess.NewAccessID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "access")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TenantAccessClient.GetEntityTag`

```go
ctx := context.TODO()
id := tenantaccess.NewAccessID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "access")

read, err := client.GetEntityTag(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TenantAccessClient.ListByService`

```go
ctx := context.TODO()
id := tenantaccess.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue")

// alternatively `client.ListByService(ctx, id, tenantaccess.DefaultListByServiceOperationOptions())` can be used to do batched pagination
items, err := client.ListByServiceComplete(ctx, id, tenantaccess.DefaultListByServiceOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `TenantAccessClient.ListSecrets`

```go
ctx := context.TODO()
id := tenantaccess.NewAccessID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "access")

read, err := client.ListSecrets(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TenantAccessClient.RegeneratePrimaryKey`

```go
ctx := context.TODO()
id := tenantaccess.NewAccessID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "access")

read, err := client.RegeneratePrimaryKey(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TenantAccessClient.RegenerateSecondaryKey`

```go
ctx := context.TODO()
id := tenantaccess.NewAccessID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "access")

read, err := client.RegenerateSecondaryKey(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TenantAccessClient.Update`

```go
ctx := context.TODO()
id := tenantaccess.NewAccessID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "access")

payload := tenantaccess.AccessInformationUpdateParameters{
	// ...
}


read, err := client.Update(ctx, id, payload, tenantaccess.DefaultUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
