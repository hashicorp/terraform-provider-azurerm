
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/openidconnectprovider` Documentation

The `openidconnectprovider` SDK allows for interaction with the Azure Resource Manager Service `apimanagement` (API Version `2021-08-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/openidconnectprovider"
```


### Client Initialization

```go
client := openidconnectprovider.NewOpenidConnectProviderClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `OpenidConnectProviderClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := openidconnectprovider.NewOpenidConnectProviderID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "openidConnectProviderValue")

payload := openidconnectprovider.OpenidConnectProviderContract{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload, openidconnectprovider.DefaultCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OpenidConnectProviderClient.Delete`

```go
ctx := context.TODO()
id := openidconnectprovider.NewOpenidConnectProviderID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "openidConnectProviderValue")

read, err := client.Delete(ctx, id, openidconnectprovider.DefaultDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OpenidConnectProviderClient.Get`

```go
ctx := context.TODO()
id := openidconnectprovider.NewOpenidConnectProviderID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "openidConnectProviderValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OpenidConnectProviderClient.GetEntityTag`

```go
ctx := context.TODO()
id := openidconnectprovider.NewOpenidConnectProviderID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "openidConnectProviderValue")

read, err := client.GetEntityTag(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OpenidConnectProviderClient.ListByService`

```go
ctx := context.TODO()
id := openidconnectprovider.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue")

// alternatively `client.ListByService(ctx, id, openidconnectprovider.DefaultListByServiceOperationOptions())` can be used to do batched pagination
items, err := client.ListByServiceComplete(ctx, id, openidconnectprovider.DefaultListByServiceOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenidConnectProviderClient.ListSecrets`

```go
ctx := context.TODO()
id := openidconnectprovider.NewOpenidConnectProviderID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "openidConnectProviderValue")

read, err := client.ListSecrets(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OpenidConnectProviderClient.Update`

```go
ctx := context.TODO()
id := openidconnectprovider.NewOpenidConnectProviderID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "openidConnectProviderValue")

payload := openidconnectprovider.OpenidConnectProviderUpdateContract{
	// ...
}


read, err := client.Update(ctx, id, payload, openidconnectprovider.DefaultUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
