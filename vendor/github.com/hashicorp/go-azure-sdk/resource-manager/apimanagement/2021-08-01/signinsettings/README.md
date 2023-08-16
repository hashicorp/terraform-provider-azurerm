
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/signinsettings` Documentation

The `signinsettings` SDK allows for interaction with the Azure Resource Manager Service `apimanagement` (API Version `2021-08-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/signinsettings"
```


### Client Initialization

```go
client := signinsettings.NewSignInSettingsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SignInSettingsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := signinsettings.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue")

payload := signinsettings.PortalSigninSettings{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload, signinsettings.DefaultCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SignInSettingsClient.Get`

```go
ctx := context.TODO()
id := signinsettings.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SignInSettingsClient.GetEntityTag`

```go
ctx := context.TODO()
id := signinsettings.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue")

read, err := client.GetEntityTag(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SignInSettingsClient.Update`

```go
ctx := context.TODO()
id := signinsettings.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue")

payload := signinsettings.PortalSigninSettings{
	// ...
}


read, err := client.Update(ctx, id, payload, signinsettings.DefaultUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
