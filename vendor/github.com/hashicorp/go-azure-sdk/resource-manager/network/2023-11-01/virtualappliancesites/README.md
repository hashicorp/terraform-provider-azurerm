
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/virtualappliancesites` Documentation

The `virtualappliancesites` SDK allows for interaction with Azure Resource Manager `network` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/virtualappliancesites"
```


### Client Initialization

```go
client := virtualappliancesites.NewVirtualApplianceSitesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `VirtualApplianceSitesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := virtualappliancesites.NewVirtualApplianceSiteID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkVirtualApplianceName", "virtualApplianceSiteName")

payload := virtualappliancesites.VirtualApplianceSite{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualApplianceSitesClient.Delete`

```go
ctx := context.TODO()
id := virtualappliancesites.NewVirtualApplianceSiteID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkVirtualApplianceName", "virtualApplianceSiteName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualApplianceSitesClient.Get`

```go
ctx := context.TODO()
id := virtualappliancesites.NewVirtualApplianceSiteID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkVirtualApplianceName", "virtualApplianceSiteName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualApplianceSitesClient.List`

```go
ctx := context.TODO()
id := virtualappliancesites.NewNetworkVirtualApplianceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkVirtualApplianceName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
