
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/delegationsettings` Documentation

The `delegationsettings` SDK allows for interaction with Azure Resource Manager `apimanagement` (API Version `2022-08-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/delegationsettings"
```


### Client Initialization

```go
client := delegationsettings.NewDelegationSettingsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DelegationSettingsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := delegationsettings.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName")

payload := delegationsettings.PortalDelegationSettings{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload, delegationsettings.DefaultCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DelegationSettingsClient.Get`

```go
ctx := context.TODO()
id := delegationsettings.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DelegationSettingsClient.GetEntityTag`

```go
ctx := context.TODO()
id := delegationsettings.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName")

read, err := client.GetEntityTag(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DelegationSettingsClient.ListSecrets`

```go
ctx := context.TODO()
id := delegationsettings.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName")

read, err := client.ListSecrets(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DelegationSettingsClient.Update`

```go
ctx := context.TODO()
id := delegationsettings.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName")

payload := delegationsettings.PortalDelegationSettings{
	// ...
}


read, err := client.Update(ctx, id, payload, delegationsettings.DefaultUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
