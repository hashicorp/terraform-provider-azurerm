
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-10-01/networksecurityperimeterprofiles` Documentation

The `networksecurityperimeterprofiles` SDK allows for interaction with Azure Resource Manager `network` (API Version `2024-10-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-10-01/networksecurityperimeterprofiles"
```


### Client Initialization

```go
client := networksecurityperimeterprofiles.NewNetworkSecurityPerimeterProfilesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `NetworkSecurityPerimeterProfilesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := networksecurityperimeterprofiles.NewProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkSecurityPerimeterName", "profileName")

payload := networksecurityperimeterprofiles.NspProfile{
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


### Example Usage: `NetworkSecurityPerimeterProfilesClient.Delete`

```go
ctx := context.TODO()
id := networksecurityperimeterprofiles.NewProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkSecurityPerimeterName", "profileName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NetworkSecurityPerimeterProfilesClient.Get`

```go
ctx := context.TODO()
id := networksecurityperimeterprofiles.NewProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkSecurityPerimeterName", "profileName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NetworkSecurityPerimeterProfilesClient.List`

```go
ctx := context.TODO()
id := networksecurityperimeterprofiles.NewNetworkSecurityPerimeterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkSecurityPerimeterName")

// alternatively `client.List(ctx, id, networksecurityperimeterprofiles.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, networksecurityperimeterprofiles.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
