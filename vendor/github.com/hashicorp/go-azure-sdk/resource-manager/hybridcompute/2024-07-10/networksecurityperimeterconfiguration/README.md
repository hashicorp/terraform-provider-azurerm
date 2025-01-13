
## `github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2024-07-10/networksecurityperimeterconfiguration` Documentation

The `networksecurityperimeterconfiguration` SDK allows for interaction with Azure Resource Manager `hybridcompute` (API Version `2024-07-10`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2024-07-10/networksecurityperimeterconfiguration"
```


### Client Initialization

```go
client := networksecurityperimeterconfiguration.NewNetworkSecurityPerimeterConfigurationClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `NetworkSecurityPerimeterConfigurationClient.GetByPrivateLinkScope`

```go
ctx := context.TODO()
id := networksecurityperimeterconfiguration.NewNetworkSecurityPerimeterConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateLinkScopeName", "networkSecurityPerimeterConfigurationName")

read, err := client.GetByPrivateLinkScope(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NetworkSecurityPerimeterConfigurationClient.ListByPrivateLinkScope`

```go
ctx := context.TODO()
id := networksecurityperimeterconfiguration.NewProviderPrivateLinkScopeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateLinkScopeName")

// alternatively `client.ListByPrivateLinkScope(ctx, id)` can be used to do batched pagination
items, err := client.ListByPrivateLinkScopeComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `NetworkSecurityPerimeterConfigurationClient.ReconcileForPrivateLinkScope`

```go
ctx := context.TODO()
id := networksecurityperimeterconfiguration.NewNetworkSecurityPerimeterConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateLinkScopeName", "networkSecurityPerimeterConfigurationName")

if err := client.ReconcileForPrivateLinkScopeThenPoll(ctx, id); err != nil {
	// handle the error
}
```
