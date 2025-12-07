
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/networksecurityperimeterlinkreferences` Documentation

The `networksecurityperimeterlinkreferences` SDK allows for interaction with Azure Resource Manager `network` (API Version `2025-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/networksecurityperimeterlinkreferences"
```


### Client Initialization

```go
client := networksecurityperimeterlinkreferences.NewNetworkSecurityPerimeterLinkReferencesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `NetworkSecurityPerimeterLinkReferencesClient.Delete`

```go
ctx := context.TODO()
id := networksecurityperimeterlinkreferences.NewLinkReferenceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkSecurityPerimeterName", "linkReferenceName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `NetworkSecurityPerimeterLinkReferencesClient.Get`

```go
ctx := context.TODO()
id := networksecurityperimeterlinkreferences.NewLinkReferenceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkSecurityPerimeterName", "linkReferenceName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NetworkSecurityPerimeterLinkReferencesClient.List`

```go
ctx := context.TODO()
id := networksecurityperimeterlinkreferences.NewNetworkSecurityPerimeterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkSecurityPerimeterName")

// alternatively `client.List(ctx, id, networksecurityperimeterlinkreferences.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, networksecurityperimeterlinkreferences.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
