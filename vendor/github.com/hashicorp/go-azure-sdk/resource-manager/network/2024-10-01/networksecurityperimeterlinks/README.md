
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-10-01/networksecurityperimeterlinks` Documentation

The `networksecurityperimeterlinks` SDK allows for interaction with Azure Resource Manager `network` (API Version `2024-10-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-10-01/networksecurityperimeterlinks"
```


### Client Initialization

```go
client := networksecurityperimeterlinks.NewNetworkSecurityPerimeterLinksClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `NetworkSecurityPerimeterLinksClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := networksecurityperimeterlinks.NewNetworkSecurityPerimeterLinkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkSecurityPerimeterName", "linkName")

payload := networksecurityperimeterlinks.NspLink{
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


### Example Usage: `NetworkSecurityPerimeterLinksClient.Delete`

```go
ctx := context.TODO()
id := networksecurityperimeterlinks.NewNetworkSecurityPerimeterLinkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkSecurityPerimeterName", "linkName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `NetworkSecurityPerimeterLinksClient.Get`

```go
ctx := context.TODO()
id := networksecurityperimeterlinks.NewNetworkSecurityPerimeterLinkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkSecurityPerimeterName", "linkName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NetworkSecurityPerimeterLinksClient.List`

```go
ctx := context.TODO()
id := networksecurityperimeterlinks.NewNetworkSecurityPerimeterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkSecurityPerimeterName")

// alternatively `client.List(ctx, id, networksecurityperimeterlinks.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, networksecurityperimeterlinks.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
